package std

import (
	"crypto/tls"
	"fmt"
	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestStdNested(t *testing.T) {
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	var tlsConfig *tls.Config
	if useSSL {
		tlsConfig = &tls.Config{}
	}
	conn, err := GetStdOpenDBConnection(datastore.Native, datastore.Settings{
		"flatten_nested": 0,
	}, tlsConfig, nil)
	require.NoError(t, err)
	conn.Exec("DROP TABLE std_nested_test")
	if !CheckMinServerVersion(conn, 22, 1, 0) {
		t.Skip(fmt.Errorf("unsupported datastore version"))
		return
	}
	const ddl = `
			CREATE TABLE std_nested_test (
				Col1 Nested(
					  Col1_N1 UInt8
					, Col2_N1 UInt8
				)
				, Col2 Nested(
					  Col1_N2 UInt8
					, Col2_N2 Nested(
						  Col1_N2_N1 UInt8
						, Col2_N2_N1 UInt8
					)
				)
			) Engine MergeTree() ORDER BY tuple()`
	defer func() {
		conn.Exec("DROP TABLE std_nested_test")
	}()
	_, err = conn.Exec(ddl)
	require.NoError(t, err)
	require.NoError(t, err)
	scope, err := conn.Begin()
	require.NoError(t, err)
	batch, err := scope.Prepare("INSERT INTO std_nested_test")
	require.NoError(t, err)
	var (
		col1Data = []map[string]any{
			{
				"Col1_N1": uint8(1),
				"Col2_N1": uint8(20),
			},
			{
				"Col1_N1": uint8(2),
				"Col2_N1": uint8(20),
			},
			{
				"Col1_N1": uint8(3),
				"Col2_N1": uint8(20),
			},
		}
		col2Data = []map[string]any{
			{
				"Col1_N2": uint8(101),
				"Col2_N2": []map[string]any{
					{
						"Col1_N2_N1": uint8(1),
						"Col2_N2_N1": uint8(2),
					},
				},
			},
			{
				"Col1_N2": uint8(201),
				"Col2_N2": []map[string]any{
					{
						"Col1_N2_N1": uint8(3),
						"Col2_N2_N1": uint8(4),
					},
				},
			},
		}
	)

	_, err = batch.Exec(col1Data, col2Data)
	require.NoError(t, err)
	require.NoError(t, scope.Commit())
	var (
		col1 []map[string]any
		col2 []map[string]any
	)
	rows := conn.QueryRow("SELECT * FROM std_nested_test")
	require.NoError(t, rows.Scan(&col1, &col2))
	assert.JSONEq(t, ToJson(col1Data), ToJson(col1))
	assert.JSONEq(t, ToJson(col2Data), ToJson(col2))
}

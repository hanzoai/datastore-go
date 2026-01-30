package std

import (
	"fmt"
	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdBool(t *testing.T) {
	dsns := map[string]datastore.Protocol{"Native": datastore.Native, "Http": datastore.HTTP}
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	for name, protocol := range dsns {
		t.Run(fmt.Sprintf("%s Protocol", name), func(t *testing.T) {
			if conn, err := GetStdDSNConnection(protocol, useSSL, nil); assert.NoError(t, err) {
				if !CheckMinServerVersion(conn, 21, 12, 0) {
					t.Skip(fmt.Errorf("unsupported datastore version"))
					return
				}
				const ddl = `
			CREATE TABLE std_test_bool (
				    Col1 Bool
				  , Col2 Bool
				  , Col3 Array(Bool)
				  , Col4 Nullable(Bool)
				  , Col5 Array(Nullable(Bool))
			) Engine MergeTree() ORDER BY tuple()
		`
				defer func() {
					conn.Exec("DROP TABLE std_test_bool")
				}()
				_, err := conn.Exec(ddl)
				require.NoError(t, err)
				scope, err := conn.Begin()
				require.NoError(t, err)
				batch, err := scope.Prepare("INSERT INTO std_test_bool")
				require.NoError(t, err)
				var val bool
				_, err = batch.Exec(true, false, []bool{true, false, true}, nil, []*bool{&val, nil, &val})
				require.NoError(t, err)
				require.NoError(t, scope.Commit())
				var (
					col1 bool
					col2 bool
					col3 []bool
					col4 *bool
					col5 []*bool
				)
				require.NoError(t, conn.QueryRow("SELECT * FROM std_test_bool").Scan(&col1, &col2, &col3, &col4, &col5))
				assert.Equal(t, true, col1)
				assert.Equal(t, false, col2)
				assert.Equal(t, []bool{true, false, true}, col3)
				require.Nil(t, col4)
				assert.Equal(t, []*bool{&val, nil, &val}, col5)
			}
		})
	}
}

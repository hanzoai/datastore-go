package issues

import (
	"context"
	"strconv"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	datastore_std_tests "github.com/hanzoai/datastore-go/v2/tests/std"
	"github.com/stretchr/testify/require"
)

func Test783(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", datastore.Settings{
			"flatten_nested": 1,
		}, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	ctx := context.Background()
	require.NoError(t, err)
	row := conn.QueryRow(ctx, "SELECT groupArray(('a', ['time1', 'time2'])) as val")
	var x [][]any
	require.NoError(t, row.Scan(&x))
	require.Equal(t, [][]any{{"a", []string{"time1", "time2"}}}, x)
}

func TestStd783(t *testing.T) {
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	conn, err := datastore_std_tests.GetDSNConnection("issues", datastore.Native, useSSL, nil)
	require.NoError(t, err)
	row := conn.QueryRow("SELECT groupArray(('a', ['time1', 'time2'])) as val")
	var x [][]any
	require.NoError(t, row.Scan(&x))
	require.Equal(t, [][]any{{"a", []string{"time1", "time2"}}}, x)
}

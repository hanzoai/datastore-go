package issues

import (
	"context"
	"strconv"
	"testing"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	clickhouse_std_tests "github.com/hanzoai/datastore-go/tests/std"
	"github.com/stretchr/testify/require"
)

func Test762(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		})
	)
	rows, err := conn.Query(context.Background(), "SELECT (NULL, NULL)")
	require.NoError(t, err)
	for rows.Next() {
		var (
			n []any
		)
		require.NoError(t, rows.Scan(&n))
		require.Equal(t, []any{(*any)(nil), (*any)(nil)}, n)
	}

	require.NoError(t, rows.Close())
	require.NoError(t, rows.Err())
}

func Test762Std(t *testing.T) {
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("CLICKHOUSE_USE_SSL", "false"))
	require.NoError(t, err)
	conn, err := clickhouse_std_tests.GetDSNConnection("issues", clickhouse.Native, useSSL, nil)
	rows, err := conn.Query("SELECT tuple(NULL)")
	require.NoError(t, err)
	for rows.Next() {
		var (
			n any
		)
		require.NoError(t, rows.Scan(&n))
		expected := []any{(*any)(nil)}
		require.Equal(t, expected, n)
	}

	require.NoError(t, rows.Close())
	require.NoError(t, rows.Err())
}

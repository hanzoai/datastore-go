package issues

import (
	"context"
	"testing"
	"time"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1066(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", datastore.Settings{
			"max_execution_time": 60,
		}, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	ctx := context.Background()
	require.NoError(t, err)
	const ddl = `
		CREATE TABLE test_1066 (
			Col1 Date32
		) Engine MergeTree() ORDER BY tuple()
		`
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS test_1066")
	}()
	require.NoError(t, conn.Exec(ctx, ddl))

	expectedDate := time.Date(2010, 10, 10, 0, 0, 0, 0, time.UTC)

	require.NoError(t, conn.Exec(ctx, `INSERT INTO test_1066 (Col1) VALUES(?)`, expectedDate))

	row := conn.QueryRow(ctx, "SELECT Col1 FROM test_1066")
	require.NoError(t, err)
	var actualDate time.Time
	require.NoError(t, row.Scan(&actualDate))

	assert.Equal(t, expectedDate, actualDate)
}

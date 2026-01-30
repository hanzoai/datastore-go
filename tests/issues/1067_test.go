package issues

import (
	"context"
	"testing"
	"time"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func Test1067(t *testing.T) {
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
	require.NoError(t, conn.Exec(ctx, ddl))
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS test_1066")
	}()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO test_1066")
	require.NoError(t, err)
	require.NoError(t, batch.Append("1970-01-02"))
	require.NoError(t, batch.Append(time.Now()))
	require.NoError(t, batch.Send())
}

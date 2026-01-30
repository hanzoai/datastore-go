package issues

import (
	"context"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func Test1049(t *testing.T) {
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
		CREATE TABLE test_1049 (
			test String
		) Engine MergeTree() ORDER BY tuple()
		`
	require.NoError(t, conn.Exec(ctx, ddl))
	defer func() {
		conn.Exec(ctx, "DROP TABLE IF EXISTS test_1049")
	}()

	_, err = conn.PrepareBatch(context.Background(), "INSERT INTO test_1049 (\"test\")")
	require.NoError(t, err)
}

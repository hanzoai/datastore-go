package tests

import (
	"context"
	"fmt"
	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuotedDDL(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		ctx := context.Background()
		require.NoError(t, err)
		require.NoError(t, conn.Ping(ctx))
		if !CheckMinServerServerVersion(conn, 21, 9, 0) {
			t.Skip(fmt.Errorf("unsupported datastore version"))
			return
		}
		const ddl = "CREATE TABLE `test_string` (`1` String) Engine MergeTree() ORDER BY tuple()"

		defer func() {
			conn.Exec(ctx, "DROP TABLE `test_string`")
		}()
		require.NoError(t, conn.Exec(ctx, ddl))
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO `test_string`")
		require.NoError(t, err)
		require.NoError(t, batch.Append("A"))
		require.NoError(t, batch.Send())
	})
}

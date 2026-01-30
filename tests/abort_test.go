package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestAbort(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		require.NoError(t, err)
		ctx := context.Background()
		const ddl = `
		CREATE TABLE test_abort (
			Col1 UInt8
		) Engine MergeTree() ORDER BY tuple()
		`
		defer func() {
			conn.Exec(ctx, "DROP TABLE IF EXISTS test_abort")
		}()
		require.NoError(t, conn.Exec(ctx, ddl))
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO test_abort")
		require.NoError(t, err)
		require.NoError(t, batch.Abort())
		if err := batch.Abort(); assert.Error(t, err) {
			assert.Equal(t, datastore.ErrBatchAlreadySent, err)
		}
		batch, err = conn.PrepareBatch(ctx, "INSERT INTO test_abort")
		require.NoError(t, err)
		if assert.NoError(t, batch.Append(uint8(1))) && assert.NoError(t, batch.Send()) {
			var col1 uint8
			if err := conn.QueryRow(ctx, "SELECT * FROM test_abort").Scan(&col1); assert.NoError(t, err) {
				assert.Equal(t, uint8(1), col1)
			}
		}
	})
}

func TestBatchClose(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		require.NoError(t, err)
		ctx := context.Background()

		if protocol == datastore.HTTP {
			// For HTTP, provide specific column names since we can't parse out the null table function
			ctx = datastore.Context(ctx,
				datastore.WithColumnNamesAndTypes([]datastore.ColumnNameAndType{
					{Name: "x", Type: "UInt64"},
				}))
		}

		batch, err := conn.PrepareBatch(ctx, "INSERT INTO function null('x UInt64') VALUES (1)")
		require.NoError(t, err)
		require.NoError(t, batch.Close())
		require.NoError(t, batch.Close()) // No error on multiple calls

		batch, err = conn.PrepareBatch(ctx, "INSERT INTO function null('x UInt64') VALUES (1)")
		require.NoError(t, err)
		if assert.NoError(t, batch.Append(uint8(1))) && assert.NoError(t, batch.Send()) {
			var col1 uint8
			if err := conn.QueryRow(ctx, "SELECT 1").Scan(&col1); assert.NoError(t, err) {
				assert.Equal(t, uint8(1), col1)
			}
		}
	})
}

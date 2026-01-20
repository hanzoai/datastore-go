package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/assert"
)

func TestNothing(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		require.NoError(t, err)
		rows, err := conn.Query(context.Background(), "SELECT NULL FROM system.numbers_mt LIMIT 10")
		require.NoError(t, err)
		var count int
		for rows.Next() {
			var nothing []struct{}
			if !assert.NoError(t, rows.Scan(&nothing)) {
				return
			}
			count++
		}
		require.NoError(t, rows.Close())
		require.NoError(t, rows.Err())
		assert.Equal(t, 10, count)
	})
}

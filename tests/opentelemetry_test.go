package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestOpenTelemetry(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		require.NoError(t, err)
		var count uint64
		rows := conn.QueryRow(datastore.Context(context.Background(), datastore.WithSpan(
			trace.NewSpanContext(trace.SpanContextConfig{
				SpanID:  trace.SpanID{1, 2, 3, 4, 5},
				TraceID: trace.TraceID{5, 4, 3, 2, 1},
			}),
		)), "SELECT COUNT() FROM (SELECT number FROM system.numbers LIMIT 5)")
		require.NoError(t, rows.Scan(&count))
		assert.Equal(t, uint64(5), count)
	})
}

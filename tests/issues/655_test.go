package issues

import (
	"context"
	"testing"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test655 confirms an agreed semantic on failing batch append results with entire batch cancellation.
func Test655(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", datastore.Settings{
			"max_execution_time": 60,
		}, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		ctx = context.Background()
	)

	require.NoError(t, err)
	conn.Exec(ctx, "DROP TABLE test_enum")
	const ddl = `CREATE TABLE test_enum (
				Col1 Enum8 ('Click'=5, 'House'=25)
			) Engine Memory`
	require.NoError(t, conn.Exec(ctx, ddl))

	defer func() {
		conn.Exec(ctx, "DROP TABLE test_enum")
	}()
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO test_enum")
	require.NoError(t, err)
	type request struct {
		Col1 string
	}
	require.Error(t, batch.AppendStruct(&request{Col1: "house"}), "datastore [AppendRow]: (Col1 Enum8('Click' = 5, 'House' = 25)) unknown element \"house\"")
	assert.ErrorContains(t, batch.Send(), "datastore: batch is invalid. check appended data is correct")
}

package issues

import (
	"context"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/assert"
)

func TestIssue476(t *testing.T) {
	var (
		ctx       = context.Background()
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	require.NoError(t, err)

	const ddl = `
			CREATE TABLE issue_476 (
				  Col1 Array(LowCardinality(String))
				, Col2 Array(LowCardinality(String))
			) Engine MergeTree() ORDER BY tuple()
		`
	defer func() {
		conn.Exec(ctx, "DROP TABLE issue_476")
	}()
	require.NoError(t, conn.Exec(ctx, ddl))
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO issue_476")
	require.NoError(t, err)
	require.NoError(t, batch.Append(
		[]string{"A", "B", "C"},
		[]string{},
	))
	require.NoError(t, batch.Send())
	var (
		col1 []string
		col2 []string
	)
	require.NoError(t, conn.QueryRow(ctx, `SELECT * FROM issue_476`).Scan(&col1, &col2))
	assert.Equal(t, []string{"A", "B", "C"}, col1)
	assert.Equal(t, []string{}, col2)
}

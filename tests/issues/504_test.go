package issues

import (
	"context"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/assert"
)

func TestIssue504(t *testing.T) {
	var (
		ctx       = context.Background()
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	require.NoError(t, err)
	require.NoError(t, err)
	var result []struct {
		Col1 string
		Col2 uint64
	}
	const query = `
		SELECT *
		FROM
		(
			SELECT
				'A'    AS Col1,
				number AS Col2
			FROM
			(
				SELECT number
				FROM system.numbers
				LIMIT 5
			)
		)
		WHERE (Col1, Col2) IN (@GS)
		`
	err = conn.Select(ctx, &result, query, datastore.Named("GS", []datastore.GroupSet{
		{Value: []any{"A", 2}},
		{Value: []any{"A", 4}},
	}))
	require.NoError(t, err)
	assert.Equal(t, []struct {
		Col1 string
		Col2 uint64
	}{
		{
			Col1: "A",
			Col2: 2,
		},
		{
			Col1: "A",
			Col2: 4,
		},
	}, result)
}

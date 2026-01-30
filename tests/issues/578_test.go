package issues

import (
	"context"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssue578(t *testing.T) {
	var (
		ctx       = context.Background()
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, nil)
	)
	require.NoError(t, err)
	assert.NoError(t, err)

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO non_existent_table")
	assert.Error(t, err)

	if batch != nil {
		batch.Abort()
	}
}

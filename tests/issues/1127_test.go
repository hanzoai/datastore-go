package issues

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1127(t *testing.T) {
	t.Skip("This test is flaky and needs to be fixed")

	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, nil)
	)
	require.NoError(t, err)

	progressHasTriggered := false
	ctx := datastore.Context(context.Background(), datastore.WithProgress(func(p *datastore.Progress) {
		fmt.Println("progress: ", p)
		progressHasTriggered = true
	}), datastore.WithLogs(func(log *datastore.Log) {
		fmt.Println("log info: ", log)
	}))

	rows, err := conn.Query(ctx, "select number, throwIf(number = 1e6) from system.numbers settings max_block_size = 100")
	require.NoError(t, err)
	defer rows.Close()

	var number uint64
	var throwIf uint8
	for rows.Next() {
		require.NoError(t, rows.Scan(&number, &throwIf))
	}

	assert.Error(t, rows.Err())
	assert.True(t, progressHasTriggered)
}

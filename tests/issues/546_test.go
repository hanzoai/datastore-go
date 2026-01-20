package issues

import (
	"context"
	"fmt"
	"testing"
	"time"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/assert"
)

func TestIssue546(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	require.NoError(t, err)
	ctx := datastore.Context(context.Background(), datastore.WithSettings(datastore.Settings{
		"max_block_size": 2000000,
	}),
		datastore.WithProgress(func(p *datastore.Progress) {
			fmt.Println("progress: ", p)
		}), datastore.WithProfileInfo(func(p *datastore.ProfileInfo) {
			fmt.Println("profile info: ", p)
		}))
	require.NoError(t, conn.Ping(context.Background()))
	if exception, ok := err.(*datastore.Exception); ok {
		fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
	}
	assert.NoError(t, err)

	rows, err := conn.Query(ctx, "SELECT * FROM system.numbers LIMIT 2000000", time.Now())
	assert.NoError(t, err)
	i := 0
	for rows.Next() {
		var (
			col1 uint64
		)
		if err := rows.Scan(&col1); err != nil {
			assert.NoError(t, err)
		}
		i += 1
	}
	require.NoError(t, rows.Close())
	require.NoError(t, rows.Err())
	assert.Equal(t, 2000000, i)
}

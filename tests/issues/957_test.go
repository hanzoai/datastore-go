package issues

import (
	"context"
	"testing"
	"time"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"
)

func Test957(t *testing.T) {
	// given
	ctx := context.Background()
	testEnv, err := datastore_tests.GetTestEnvironment(testSet)
	require.NoError(t, err)

	// when the client is configured to use the test environment
	opts := datastore_tests.ClientOptionsFromEnv(testEnv, datastore.Settings{}, false)
	// and the client is configured to have only 1 connection
	opts.MaxIdleConns = 2
	opts.MaxOpenConns = 1
	// and the client is configured to have a connection lifetime of 1/10 of a second
	opts.ConnMaxLifetime = time.Second / 10
	conn, err := datastore.Open(&opts)
	require.NoError(t, err)

	// then the client should be able to execute queries for 1 second
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		rows, err := conn.Query(ctx, "SELECT 1")
		require.NoError(t, err)
		rows.Close()
	}
}

package issues

import (
	"context"
	"testing"
	"time"

	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func Test1016(t *testing.T) {
	testEnv, err := datastore_tests.GetTestEnvironment("issues")
	require.NoError(t, err)
	conn, err := datastore_tests.TestClientWithDefaultSettings(testEnv)
	require.NoError(t, err)

	rows, err := conn.Query(context.Background(), "SELECT ?;", time.Unix(0, 0).UTC())
	require.NoError(t, err)
	defer rows.Close()

	require.True(t, rows.Next())
	var v time.Time
	err = rows.Scan(&v)
	require.NoError(t, err)
	require.Equal(t, time.Unix(0, 0).UTC(), v)
}

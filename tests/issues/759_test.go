package issues

import (
	"context"
	"testing"
	"time"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/hanzoai/datastore-go/v2/lib/driver"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func Test759(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", datastore.Settings{
			"max_execution_time": 60,
		}, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	require.NoError(t, err)
	timeWant, err := time.Parse(time.RFC3339Nano, "2022-09-15T17:06:31.81718722+04:00")
	require.NoError(t, err)
	testWith(t, conn, timeWant.Local())
	testWith(t, conn, timeWant)

}

func testWith(t *testing.T, conn driver.Conn, timeWant time.Time) {
	date := datastore.DateNamed("Time", timeWant, datastore.NanoSeconds)
	r := conn.QueryRow(context.TODO(), "SELECT @Time", date)

	var timeGot time.Time
	require.NoError(t, r.Scan(&timeGot))
	require.Equal(t, timeGot.Unix(), timeWant.Unix())
}

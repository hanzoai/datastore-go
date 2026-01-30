package issues

import (
	"context"
	"testing"
	"time"

	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestIssue548(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)

	assert.NoError(t, err)
	// give it plenty of time before we conclusively assume deadlock
	timeout := time.After(5 * time.Second)
	done := make(chan bool)
	go func() {
		// should take 1s
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		rows, _ := conn.Query(ctx, "SELECT sleepEachRow(0.001) as Col1 FROM system.numbers LIMIT 1000 SETTINGS max_block_size=10;")
		rows.Close()
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Close() deadlocked")
	case <-done:
	}
}

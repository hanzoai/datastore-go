package issues

import (
	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	datastore_std_tests "github.com/hanzoai/datastore-go/tests/std"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test470PR(t *testing.T) {
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	conn, err := datastore_std_tests.GetDSNConnection("issues", datastore.Native, useSSL, nil)
	require.NoError(t, err)
	const ddl = `
		CREATE TABLE issue_470_pr (
			Col1 Array(String)
		) Engine MergeTree() ORDER BY tuple()
		`
	defer func() {
		conn.Exec("DROP TABLE issue_470_pr")
	}()
	_, err = conn.Exec(ddl)
	require.NoError(t, err)
	scope, err := conn.Begin()
	require.NoError(t, err)
	batch, err := scope.Prepare("INSERT INTO issue_470_pr")
	require.NoError(t, err)
	_, err = batch.Exec(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "converting <nil> to Array(String) is unsupported")
}

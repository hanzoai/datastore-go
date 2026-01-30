package issues

import (
	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	datastore_std_tests "github.com/hanzoai/datastore-go/v2/tests/std"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func Test813(t *testing.T) {
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	conn, err := datastore_std_tests.GetDSNConnection("issues", datastore.Native, useSSL, nil)
	const ddl = `
		CREATE TABLE test_813 (
		  	IntValue Int64,
			Exemplars Nested (
				Attributes Map(LowCardinality(String), String)
			) CODEC(ZSTD(1)) 
		) Engine MergeTree() ORDER BY tuple()
		`
	conn.Exec("DROP TABLE test_813")
	defer func() {
		conn.Exec("DROP TABLE test_813")
	}()
	_, err = conn.Exec(ddl)
	require.NoError(t, err)

	valueArgs := []any{
		int64(14),
		datastore.ArraySet{map[string]string{"array1_key1": "array1_value2", "array1_key2": "array1_value2"}},
	}
	_, err = conn.Exec("INSERT INTO test_813 (IntValue, Exemplars.Attributes) VALUES (?,?)", valueArgs...)
	require.NoError(t, err)
}

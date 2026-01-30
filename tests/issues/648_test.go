package issues

import (
	"context"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	"github.com/stretchr/testify/require"
)

func TestIssue648(t *testing.T) {
	var (
		conn, err = datastore_tests.GetConnectionTCP("issues", datastore.Settings{
			"max_execution_time": 60,
		}, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
	)
	conn.Exec(context.Background(), "DROP TABLE IF EXISTS issue_648")
	require.NoError(t, err)
	require.NoError(t, conn.Exec(
		context.Background(),
		`CREATE TABLE issue_648(
				id Int64,
				arr Array(UInt8), 
				map Map(String, UInt8), 
				mul_arr Array(Array(UInt8)), 
				map_arr Map(UInt8, Array(UInt8)), 
				map_map_arr Map(String, Map(String, Array(UInt8))))
			ENGINE = MergeTree
			ORDER BY (id)`,
	))
	defer func() {
		require.NoError(t, conn.Exec(context.Background(), "DROP TABLE issue_648"))
	}()
	ctx := context.Background()
	require.NoError(t, err)
	for i := uint8(0); i < 10; i++ {
		require.NoError(t, conn.Exec(ctx, "INSERT INTO issue_648 VALUES(?, ?, ?, ?, ?, ?)",
			int64(i),
			datastore.ArraySet{i, i + 1, i + 2},
			map[string]uint8{string([]byte{'A' + i}): i},
			datastore.ArraySet{datastore.ArraySet{i, i + 1}, datastore.ArraySet{i + 1, i + 2}},
			map[uint8][]uint8{i: {i + 1, i + 2}, i + 1: {i + 2, i + 3}},
			map[string]map[string][]uint8{string([]byte{'A' + i}): {string([]byte{'A' + i}): {i + 1, i + 2}}},
		))
	}
	// update array
	require.NoError(t, conn.Exec(ctx, "ALTER TABLE issue_648 UPDATE arr = ? where id = ?", datastore.ArraySet{1, 1}, 0))
	// update map
	require.NoError(t, conn.Exec(ctx, "ALTER TABLE issue_648 UPDATE map = ? where id = ?", map[string]uint8{"a": 1}, 0))
	// update mul_array
	require.NoError(t, conn.Exec(ctx, "ALTER TABLE issue_648 UPDATE mul_arr = ? where id = ?", datastore.ArraySet{datastore.ArraySet{1, 2}, datastore.ArraySet{2, 3}}, 0))
	// update map_arr
	require.NoError(t, conn.Exec(ctx, "ALTER TABLE issue_648 UPDATE map_arr = ? where id = ?", map[uint8][]uint8{1: {1, 2, 3}}, 0))
	// update map_map_arr
	require.NoError(t, conn.Exec(ctx, "ALTER TABLE issue_648 UPDATE map_map_arr = ? where id = ?", map[string]map[string][]uint8{"A": {"B": {1, 2}}}, 0))
}

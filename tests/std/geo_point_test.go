package std

import (
	"context"
	"fmt"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"

	"github.com/hanzoai/datastore-go"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
)

func TestStdGeoPoint(t *testing.T) {
	ctx := datastore.Context(context.Background(), datastore.WithSettings(datastore.Settings{
		"allow_experimental_geo_types": 1,
	}))
	dsns := map[string]datastore.Protocol{"Native": datastore.Native, "Http": datastore.HTTP}
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	for name, protocol := range dsns {
		t.Run(fmt.Sprintf("%s Protocol", name), func(t *testing.T) {
			conn, err := GetStdDSNConnection(protocol, useSSL, nil)
			require.NoError(t, err)
			if !CheckMinServerVersion(conn, 21, 12, 0) {
				t.Skip(fmt.Errorf("unsupported datastore version"))
				return
			}
			const ddl = `
				CREATE TABLE std_test_geo_point (
					Col1 Point
					, Col2 Array(Point)
				) Engine MergeTree() ORDER BY tuple()
				`
			defer func() {
				conn.Exec("DROP TABLE std_test_geo_point")
			}()
			_, err = conn.ExecContext(ctx, ddl)
			require.NoError(t, err)
			scope, err := conn.Begin()
			require.NoError(t, err)
			batch, err := scope.Prepare("INSERT INTO std_test_geo_point")
			require.NoError(t, err)
			_, err = batch.Exec(
				orb.Point{11, 22},
				[]orb.Point{
					{1, 2},
					{3, 4},
				},
			)
			require.NoError(t, err)
			require.NoError(t, scope.Commit())
			var (
				col1 orb.Point
				col2 []orb.Point
			)
			require.NoError(t, conn.QueryRow("SELECT * FROM std_test_geo_point").Scan(&col1, &col2))
			assert.Equal(t, orb.Point{11, 22}, col1)
			assert.Equal(t, []orb.Point{
				{1, 2},
				{3, 4},
			}, col2)
		})
	}
}

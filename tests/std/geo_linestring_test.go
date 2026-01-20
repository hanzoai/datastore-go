package std

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
)

func TestStdGeoLineString(t *testing.T) {
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
				CREATE TABLE std_test_geo_linestring (
					Col1 LineString
					, Col2 Array(LineString)
				) Engine MergeTree() ORDER BY tuple()
				`
			defer func() {
				conn.Exec("DROP TABLE std_test_geo_linestring")
			}()
			_, err = conn.ExecContext(ctx, ddl)
			require.NoError(t, err)
			scope, err := conn.Begin()
			require.NoError(t, err)
			batch, err := scope.Prepare("INSERT INTO std_test_geo_linestring")
			require.NoError(t, err)
			var (
				col1Data = orb.LineString{
					orb.Point{1, 2},
					orb.Point{3, 4},
					orb.Point{5, 6},
				}
				col2Data = []orb.LineString{
					{
						orb.Point{1, 2},
						orb.Point{3, 4},
					},
					{
						orb.Point{5, 6},
						orb.Point{7, 8},
					},
				}
			)
			_, err = batch.Exec(col1Data, col2Data)
			require.NoError(t, err)
			require.NoError(t, scope.Commit())
			var (
				col1 orb.LineString
				col2 []orb.LineString
			)
			require.NoError(t, conn.QueryRow("SELECT * FROM std_test_geo_linestring").Scan(&col1, &col2))
			assert.Equal(t, col1Data, col1)
			assert.Equal(t, col2Data, col2)
		})
	}
}

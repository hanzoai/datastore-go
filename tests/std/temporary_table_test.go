package std

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"

	"github.com/hanzoai/datastore-go"
	"github.com/stretchr/testify/assert"
)

func TestStdTemporaryTable(t *testing.T) {
	dsns := map[string]datastore.Protocol{"Native": datastore.Native, "Http": datastore.HTTP}
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	for name, protocol := range dsns {
		t.Run(fmt.Sprintf("%s Protocol", name), func(t *testing.T) {
			ctx := context.Background()
			if name == "Http" {
				t.Skip("flaky test with HTTP")
				ctx = datastore.Context(ctx, datastore.WithSettings(datastore.Settings{
					"session_id":        "temp_table_test_session",
					"wait_end_of_query": "1",
				}))
			}

			conn, err := GetStdDSNConnection(protocol, useSSL, nil)
			require.NoError(t, err)

			_, err = conn.Exec("DROP TABLE IF EXISTS std_test_temporary_table")
			require.NoError(t, err)
			const ddl = `CREATE TEMPORARY TABLE std_test_temporary_table (
							ID UInt64
						);`
			tx, err := conn.Begin()
			require.NoError(t, err)
			conn.ExecContext(ctx, "DROP TABLE IF EXISTS std_test_temporary_table")
			_, err = tx.ExecContext(ctx, ddl)
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, "INSERT INTO std_test_temporary_table (ID) SELECT number AS ID FROM system.numbers LIMIT 10")
			require.NoError(t, err)
			rows, err := tx.QueryContext(ctx, "SELECT ID AS ID FROM std_test_temporary_table")
			require.NoError(t, err)
			var count int
			for rows.Next() {
				var num int
				if err := rows.Scan(&num); !assert.NoError(t, err) {
					return
				}
				count++
			}
			_, err = tx.QueryContext(ctx, "SELECT ID AS ID1 FROM std_test_temporary_table")
			require.NoError(t, err)
			_, err = conn.Query("SELECT ID AS ID2 FROM std_test_temporary_table")
			require.Error(t, err)
			if name == "Http" {
				assert.Contains(t, err.Error(), "Code: 60")
			} else {
				exception, ok := err.(*datastore.Exception)
				require.True(t, ok)
				assert.Equal(t, int32(60), exception.Code)
			}
			require.Equal(t, 10, count)
			require.NoError(t, tx.Commit())

			_, err = conn.Exec("DROP TABLE IF EXISTS std_test_temporary_table")
			require.NoError(t, err)
			require.NoError(t, conn.Close())
		})
	}
}

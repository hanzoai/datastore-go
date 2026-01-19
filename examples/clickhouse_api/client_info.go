package clickhouse_api

import (
	"context"
	"github.com/hanzoai/datastore-go"
)

func ClientInfo() error {
	conn, err := datastore.Open(&datastore.Options{
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return err
	}

	return conn.Exec(context.TODO(), "SELECT 1")
}

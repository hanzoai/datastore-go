package datastore_api

import (
	"context"
	"github.com/hanzoai/datastore-go/v2"
)

func ClientInfo() error {
	conn, err := datastore.Open(&datastore.Options{
		ClientInfo: datastore.ClientInfo{
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

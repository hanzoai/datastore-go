package datastore_api

import (
	"context"
	"fmt"

	"github.com/hanzoai/datastore-go"
)

func ProgressProfileLogs() error {
	conn, err := GetNativeConnection(datastore.Settings{
		"send_logs_level": "trace",
	}, nil, nil)
	if err != nil {
		return err
	}
	totalRows := uint64(0)
	// use context to pass a call back for progress and profile info
	ctx := datastore.Context(context.Background(), datastore.WithProgress(func(p *datastore.Progress) {
		fmt.Println("progress: ", p)
		totalRows += p.Rows
	}), datastore.WithProfileInfo(func(p *datastore.ProfileInfo) {
		fmt.Println("profile info: ", p)
	}), datastore.WithLogs(func(log *datastore.Log) {
		fmt.Println("log info: ", log)
	}))

	rows, err := conn.Query(ctx, "SELECT number from numbers(1000000) LIMIT 1000000")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
	}

	fmt.Printf("Total Rows: %d\n", totalRows)
	return rows.Err()
}

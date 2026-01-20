package main

import (
	"fmt"
	"time"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests/std"
)

func main() {
	conn := datastore_tests.GetConnectionWithOptions(&datastore.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: datastore.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Settings: datastore.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &datastore.Compression{
			Method: datastore.CompressionLZ4,
		},
		//Debug: true,
	})
	if err := conn.Ping(); err != nil {
		fmt.Printf("1: %v\n", err)
	}
	row := conn.QueryRow("SELECT 1")
	var one int
	if err := row.Scan(&one); err != nil {
		fmt.Printf("2: %v\n", err)
	}
	fmt.Printf("3: %v\n", one)
	if err := conn.Close(); err != nil {
		fmt.Printf("4: %v\n", err)
	}
	fmt.Printf("5\n")
}

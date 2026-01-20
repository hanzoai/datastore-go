package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hanzoai/datastore-go"
	"github.com/hanzoai/datastore-go/lib/driver"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
)

func getClickhouseClient() driver.Conn {
	conn, _ := datastore_tests.GetConnectionWithOptions(&datastore.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: datastore.Auth{
			Database: "",
			Username: "",
			Password: "",
		},
		Settings: datastore.Settings{
			"max_execution_time": 60,
		},
		DialTimeout:     5 * time.Second,
		ConnMaxLifetime: 15 * time.Second,
		Compression: &datastore.Compression{
			Method: datastore.CompressionLZ4,
		},
		// Debug: true,
	})

	return conn
}

func main() {
	conn := getClickhouseClient()
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		var result []struct {
			Test string `ch:"test"`
		}
		sql := `SELECT 'test' AS test FROM system.numbers LIMIT 10`
		if response := conn.Select(context.Background(), &result, sql); response != nil {
			fmt.Println(response.Error())
		}
		fmt.Println(result, conn.Stats())
	})
	http.ListenAndServe(":8080", nil)
}

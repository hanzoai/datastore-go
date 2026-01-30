package datastore_api

import (
	"context"
	"fmt"
	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
)

func QueryWithParameters() error {
	conn, err := GetNativeConnection(nil, nil, nil)
	if err != nil {
		return err
	}

	if !datastore_tests.CheckMinServerServerVersion(conn, 22, 8, 0) {
		return nil
	}

	chCtx := datastore.Context(context.Background(), datastore.WithParameters(datastore.Parameters{
		"str":      "hello",
		"array":    "['a', 'b', 'c']",
		"column":   "number",
		"database": "system",
		"table":    "numbers",
	}))

	row := conn.QueryRow(chCtx, "SELECT {column:Identifier} v, {str:String} s, {array:Array(String)} a FROM {database:Identifier}.{table:Identifier} LIMIT 1 OFFSET 100")
	var (
		col1 uint64
		col2 string
		col3 []string
	)
	if err := row.Scan(&col1, &col2, &col3); err != nil {
		return err
	}
	fmt.Printf("row: col1=%d, col2=%s, col3=%s\n", col1, col2, col3)
	return nil
}

package std

import (
	"fmt"
	"github.com/hanzoai/datastore-go/v2"
	"github.com/hanzoai/datastore-go/v2/tests/std"
)

func QueryWithParameters() error {
	conn, err := GetStdOpenDBConnection(datastore.Native, nil, nil, nil)
	if err != nil {
		return err
	}

	if !std.CheckMinServerVersion(conn, 22, 8, 0) {
		return nil
	}

	row := conn.QueryRow(
		"SELECT {column:Identifier} v, {str:String} s, {array:Array(String)} a FROM {database:Identifier}.{table:Identifier} LIMIT 1 OFFSET 100",
		datastore.Named("num", "42"),
		datastore.Named("str", "hello"),
		datastore.Named("array", "['a', 'b', 'c']"),
		datastore.Named("column", "number"),
		datastore.Named("database", "system"),
		datastore.Named("table", "numbers"),
	)
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

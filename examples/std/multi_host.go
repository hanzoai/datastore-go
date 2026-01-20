package std

import (
	"database/sql"
	"fmt"
	"github.com/hanzoai/datastore-go"
)

func MultiStdHost() error {
	env, err := GetStdTestEnvironment()
	if err != nil {
		return err
	}
	conn, err := datastore.Open(&datastore.Options{
		Addr: []string{"127.0.0.1:9001", "127.0.0.1:9002", fmt.Sprintf("%s:%d", env.Host, env.Port)},
		Auth: datastore.Auth{
			Database: env.Database,
			Username: env.Username,
			Password: env.Password,
		},
		ConnOpenStrategy: datastore.ConnOpenRoundRobin,
	})
	if err != nil {
		return err
	}
	v, err := conn.ServerVersion()
	if err != nil {
		return err
	}
	fmt.Println(v.String())
	return nil
}

func MultiStdHostDSN() error {
	env, err := GetStdTestEnvironment()
	if err != nil {
		return err
	}
	conn, err := sql.Open("datastore", fmt.Sprintf("datastore://127.0.0.1:9001,127.0.0.1:9002,%s:%d?username=%s&password=%s&connection_open_strategy=round_robin", env.Host, env.Port, env.Username, env.Password))
	if err != nil {
		return err
	}
	return conn.Ping()
}

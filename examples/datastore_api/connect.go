package datastore_api

import (
	"fmt"
	"github.com/hanzoai/datastore-go"
)

func Connect() error {
	env, err := GetNativeTestEnvironment()
	if err != nil {
		return err
	}
	conn, err := datastore.Open(&datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, env.Port)},
		Auth: datastore.Auth{
			Database: env.Database,
			Username: env.Username,
			Password: env.Password,
		},
	})
	if err != nil {
		return err
	}
	v, err := conn.ServerVersion()
	fmt.Println(v)
	if err != nil {
		return err
	}
	return nil
}

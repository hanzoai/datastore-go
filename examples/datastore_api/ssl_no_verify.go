package datastore_api

import (
	"crypto/tls"
	"fmt"
	"github.com/hanzoai/datastore-go"
)

func SSLNoVerifyVersion() error {
	env, err := GetNativeTestEnvironment()
	if err != nil {
		return err
	}
	conn, err := datastore.Open(&datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, env.SslPort)},
		Auth: datastore.Auth{
			Database: env.Database,
			Username: env.Username,
			Password: env.Password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
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

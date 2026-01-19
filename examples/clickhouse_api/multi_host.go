package clickhouse_api

import (
	"fmt"
	"math/rand"

	"github.com/hanzoai/datastore-go"
)

func MultiHostVersion() error {
	return multiHostVersion(nil)
}

func MultiHostRoundRobinVersion() error {
	connOpenStrategy := datastore.ConnOpenRoundRobin
	return multiHostVersion(&connOpenStrategy)
}

func MultiHostRandomVersion() error {
	rand.Seed(85206178671753424)
	defer ResetRandSeed()
	connOpenStrategy := datastore.ConnOpenRandom
	return multiHostVersion(&connOpenStrategy)
}

func multiHostVersion(connOpenStrategy *datastore.ConnOpenStrategy) error {
	env, err := GetNativeTestEnvironment()
	if err != nil {
		return err
	}
	options := datastore.Options{
		Addr: []string{"127.0.0.1:9001", "127.0.0.1:9002", fmt.Sprintf("%s:%d", env.Host, env.Port)},
		Auth: clickhouse.Auth{
			Database: env.Database,
			Username: env.Username,
			Password: env.Password,
		},
	}
	if connOpenStrategy != nil {
		options.ConnOpenStrategy = *connOpenStrategy
	}
	conn, err := datastore.Open(&options)
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

package issues

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
	datastore_std_tests "github.com/hanzoai/datastore-go/v2/tests/std"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestIssue570(t *testing.T) {
	env, err := GetIssuesTestEnvironment()
	require.NoError(t, err)
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	var tlsConfig *tls.Config
	dsn := fmt.Sprintf("datastore://%s:%s@%s:%d/default", env.Username, env.Password,
		env.Host, env.Port)
	port := env.Port
	if useSSL {
		tlsConfig = &tls.Config{}
		port = env.SslPort
		dsn = fmt.Sprintf("datastore://%s:%s@%s:%d/default?secure=true", env.Username, env.Password,
			env.Host, env.SslPort)
	}
	require.NoError(t, err)
	// using ParseDNS - defaults shouldn't be set for maxOpenConnections etc
	options, err := datastore.ParseDSN(dsn)
	assert.NoError(t, err)
	conn := datastore_std_tests.GetConnectionWithOptions(options)
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(10)
	assert.NoError(t, conn.Ping())
	conn.Close()

	// check we can pass Options
	options = &datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Auth: datastore.Auth{
			Database: "default",
			Username: env.Username,
			Password: env.Password,
		},
		Compression: &datastore.Compression{
			Method: datastore.CompressionLZ4,
		},
		DialTimeout: time.Second,
		TLS:         tlsConfig,
	}
	conn = datastore_std_tests.GetConnectionWithOptions(options)
	assert.NoError(t, conn.Ping())

	// check we can open with a DSN
	conn, err = sql.Open("datastore", dsn)
	require.NoError(t, err)
	assert.NoError(t, conn.Ping())
}

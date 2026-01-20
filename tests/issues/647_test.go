package issues

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"
)

func TestIssue647(t *testing.T) {
	env, err := GetIssuesTestEnvironment()
	require.NoError(t, err)
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	var tlsConfig *tls.Config
	port := env.Port
	if useSSL {
		tlsConfig = &tls.Config{}
		port = env.SslPort
	}
	options := &datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Auth: datastore.Auth{
			Database: "default",
			Username: env.Username,
			Password: env.Password,
		},
		TLS: tlsConfig,
	}
	conn, err := datastore_tests.GetConnectionWithOptions(options)
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, conn.Ping(ctx))
	//reuse options
	conn2, err := datastore_tests.GetConnectionWithOptions(options)
	require.NoError(t, err)
	require.NoError(t, conn2.Ping(ctx))
}

func TestIssue647_OpenDB(t *testing.T) {
	env, err := GetIssuesTestEnvironment()
	require.NoError(t, err)
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	var tlsConfig *tls.Config
	port := env.Port
	if useSSL {
		tlsConfig = &tls.Config{}
		port = env.SslPort
	}
	options := &datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Auth: datastore.Auth{
			Database: "default",
			Username: env.Username,
			Password: env.Password,
		},
		TLS: tlsConfig,
	}
	conn := datastore.OpenDB(options)
	require.NoError(t, conn.Ping())
	//reuse options
	conn2 := datastore.OpenDB(options)
	require.NoError(t, conn2.Ping())
	// allow nil to be parsed - should work if ClickHouse was available on 9000
	//conn3 := datastore.OpenDB(nil)
	//require.NoError(t, conn3.Ping())
}

func Test647_Connector(t *testing.T) {
	env, err := GetIssuesTestEnvironment()
	require.NoError(t, err)
	useSSL, err := strconv.ParseBool(datastore_tests.GetEnv("DATASTORE_USE_SSL", "false"))
	require.NoError(t, err)
	var tlsConfig *tls.Config
	port := env.Port
	if useSSL {
		tlsConfig = &tls.Config{}
		port = env.SslPort
	}
	options := &datastore.Options{
		Addr: []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Auth: datastore.Auth{
			Database: "default",
			Username: env.Username,
			Password: env.Password,
		},
		TLS: tlsConfig,
	}
	conn := datastore.Connector(options)
	require.NoError(t, sql.OpenDB(conn).Ping())
	// reuse options
	conn2 := datastore.Connector(options)
	require.NoError(t, sql.OpenDB(conn2).Ping())
	// allow nil to be parsed - should work if ClickHouse was available on 9000
	//conn3 := datastore.Connector(nil)
	//require.NoError(t, sql.OpenDB(conn3).Ping())
}

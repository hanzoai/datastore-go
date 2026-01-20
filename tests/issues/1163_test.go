package issues

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	"github.com/stretchr/testify/require"
)

func TestIssue1163(t *testing.T) {
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
	var debugfCalled bool
	options := &datastore.Options{
		Addr:  []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Debug: true,
		Debugf: func(format string, v ...any) {
			debugfCalled = true
		},
		Auth: datastore.Auth{
			Database: "default",
			Username: env.Username,
			Password: env.Password,
		},
		TLS: tlsConfig,
	}
	conn := datastore.Connector(options)
	c, err := conn.Connect(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.True(t, debugfCalled)
}

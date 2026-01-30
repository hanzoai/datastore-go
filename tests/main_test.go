package tests

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/hanzoai/datastore-go/v2/lib/driver"
)

const testSet string = "native"

func TestMain(m *testing.M) {
	os.Exit(Runtime(m, testSet))
}

func GetNativeTestEnvironment() (ClickHouseTestEnvironment, error) {
	return GetTestEnvironment(testSet)
}

func GetNativeConnection(t *testing.T, protocol datastore.Protocol, settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (driver.Conn, error) {
	conn, err := GetConnection(testSet, t, protocol, settings, tlsConfig, compression)
	CleanupNativeConn(t, conn)
	return conn, err
}

func GetNativeConnectionTCP(settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (driver.Conn, error) {
	return GetConnection(testSet, nil, datastore.Native, settings, tlsConfig, compression)
}

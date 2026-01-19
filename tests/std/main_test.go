package std

import (
	"crypto/tls"
	"database/sql"
	"net/url"
	"os"
	"testing"

	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
)

const testSet string = "std"

func TestMain(m *testing.M) {
	os.Exit(datastore_tests.Runtime(m, testSet))
}

func GetStdDSNConnection(protocol clickhouse.Protocol, secure bool, opts url.Values) (*sql.DB, error) {
	return GetDSNConnection(testSet, protocol, secure, opts)
}

func GetStdOpenDBConnection(protocol clickhouse.Protocol, settings clickhouse.Settings, tlsConfig *tls.Config, compression *clickhouse.Compression) (*sql.DB, error) {
	return GetOpenDBConnection(testSet, protocol, settings, tlsConfig, compression)
}

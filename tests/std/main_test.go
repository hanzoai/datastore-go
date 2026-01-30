package std

import (
	"crypto/tls"
	"database/sql"
	"net/url"
	"os"
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
)

const testSet string = "std"

func TestMain(m *testing.M) {
	os.Exit(datastore_tests.Runtime(m, testSet))
}

func GetStdDSNConnection(protocol datastore.Protocol, secure bool, opts url.Values) (*sql.DB, error) {
	return GetDSNConnection(testSet, protocol, secure, opts)
}

func GetStdOpenDBConnection(protocol datastore.Protocol, settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (*sql.DB, error) {
	return GetOpenDBConnection(testSet, protocol, settings, tlsConfig, compression)
}

package datastore_api

import (
	"crypto/tls"
	"math/rand"
	"time"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/hanzoai/datastore-go/v2/lib/driver"
	datastore_tests "github.com/hanzoai/datastore-go/v2/tests"
)

const TestSet string = "examples_datastore_api"

func GetNativeConnection(settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionTCP(TestSet, settings, tlsConfig, compression)
}

func GetHTTPConnection(sessionName string, settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionHTTP(TestSet, sessionName, settings, tlsConfig, compression)
}

func GetNativeTestEnvironment() (datastore_tests.ClickHouseTestEnvironment, error) {
	return datastore_tests.GetTestEnvironment(TestSet)
}

func GetNativeConnectionWithOptions(settings datastore.Settings, tlsConfig *tls.Config, compression *datastore.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionTCP(TestSet, settings, tlsConfig, compression)
}

func CheckMinServerVersion(conn driver.Conn, major, minor, patch uint64) bool {
	return datastore_tests.CheckMinServerServerVersion(conn, major, minor, patch)
}

var randSeed = time.Now().UnixNano()

func ResetRandSeed() {
	rand.Seed(randSeed)
}

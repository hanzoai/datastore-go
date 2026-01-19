package clickhouse_api

import (
	"crypto/tls"
	"math/rand"
	"time"

	"github.com/hanzoai/datastore-go"
	"github.com/hanzoai/datastore-go/lib/driver"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
)

const TestSet string = "examples_clickhouse_api"

func GetNativeConnection(settings clickhouse.Settings, tlsConfig *tls.Config, compression *clickhouse.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionTCP(TestSet, settings, tlsConfig, compression)
}

func GetHTTPConnection(sessionName string, settings clickhouse.Settings, tlsConfig *tls.Config, compression *clickhouse.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionHTTP(TestSet, sessionName, settings, tlsConfig, compression)
}

func GetNativeTestEnvironment() (datastore_tests.ClickHouseTestEnvironment, error) {
	return datastore_tests.GetTestEnvironment(TestSet)
}

func GetNativeConnectionWithOptions(settings clickhouse.Settings, tlsConfig *tls.Config, compression *clickhouse.Compression) (driver.Conn, error) {
	return datastore_tests.GetConnectionTCP(TestSet, settings, tlsConfig, compression)
}

func CheckMinServerVersion(conn driver.Conn, major, minor, patch uint64) bool {
	return datastore_tests.CheckMinServerServerVersion(conn, major, minor, patch)
}

var randSeed = time.Now().UnixNano()

func ResetRandSeed() {
	rand.Seed(randSeed)
}

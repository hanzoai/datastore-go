package std

import (
	"crypto/tls"
	"database/sql"
	"github.com/hanzoai/datastore-go"
	datastore_tests "github.com/hanzoai/datastore-go/tests"
	datastore_tests_std "github.com/hanzoai/datastore-go/tests/std"
)

const TestSet string = "examples_std_api"

func GetStdDSNConnection(protocol clickhouse.Protocol, secure bool, compress string) (*sql.DB, error) {
	return datastore_tests_std.GetDSNConnection(TestSet, protocol, secure, nil)
}

func GetStdOpenDBConnection(protocol clickhouse.Protocol, settings clickhouse.Settings, tlsConfig *tls.Config, compression *clickhouse.Compression) (*sql.DB, error) {
	return datastore_tests_std.GetOpenDBConnection(TestSet, protocol, settings, tlsConfig, compression)
}

func GetStdTestEnvironment() (datastore_tests.ClickHouseTestEnvironment, error) {
	return datastore_tests.GetTestEnvironment(TestSet)
}

func CheckMinServerVersion(conn *sql.DB, major, minor, patch uint64) bool {
	return datastore_tests_std.CheckMinServerVersion(conn, major, minor, patch)
}

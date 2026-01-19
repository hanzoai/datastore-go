package issues

import (
	"os"
	"testing"

	datastore_tests "github.com/hanzoai/datastore-go/tests"
)

const testSet string = "issues"

func TestMain(m *testing.M) {
	os.Exit(datastore_tests.Runtime(m, testSet))
}

func GetIssuesTestEnvironment() (datastore_tests.ClickHouseTestEnvironment, error) {
	return datastore_tests.GetTestEnvironment(testSet)
}

package tests

import (
	"testing"

	"github.com/hanzoai/datastore-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestContributors(t *testing.T) {
	TestProtocols(t, func(t *testing.T, protocol datastore.Protocol) {
		conn, err := GetNativeConnection(t, protocol, nil, nil, &datastore.Compression{
			Method: datastore.CompressionLZ4,
		})
		if assert.NoError(t, err) {
			for _, contributor := range conn.Contributors() {
				t.Log(contributor)
			}
		}
	})
}

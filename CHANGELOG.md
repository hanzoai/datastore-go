# Changelog

## 1.0.0 (2026-01-19)

Initial release of datastore-go, a rebranded and streamlined ClickHouse Go client.

### Features

* Native ClickHouse TCP protocol support
* HTTP protocol support for proxied connections
* `database/sql` interface compatibility
* Connection pooling with failover and load balancing
* LZ4/ZSTD compression support
* Batch inserts with columnar support
* Struct marshaling/unmarshaling
* Query parameters and named placeholders
* OpenTelemetry integration
* Async insert support
* TLS/SSL support
* External tables support
* JSON type support
* Variant and Dynamic type support
* Time and Time64 datatype support

### Breaking Changes from clickhouse-go v2

This is a clean break from the v2 line with no backwards compatibility:

* Module path changed to `github.com/hanzoai/datastore-go`
* Package name changed from `clickhouse` to `datastore`
* DSN scheme changed from `clickhouse://` to `datastore://`
* SQL driver registration name changed from `clickhouse` to `datastore`

### Migration from clickhouse-go v2

Update your imports:

```go
// Before
import "github.com/ClickHouse/clickhouse-go/v2"
clickhouse.Open(&clickhouse.Options{...})

// After
import "github.com/hanzoai/datastore-go"
datastore.Open(&datastore.Options{...})
```

Update your DSN strings:

```go
// Before
sql.Open("clickhouse", "clickhouse://localhost:9000/default")

// After
sql.Open("datastore", "datastore://localhost:9000/default")
```

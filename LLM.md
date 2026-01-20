# LLM.md - Datastore-Go Project Documentation

## Project Overview

Datastore-go is a Go client library for ClickHouse databases, rebranded from the original ClickHouse-go project. It provides both native TCP protocol and HTTP protocol support for connecting to ClickHouse servers.

**Repository:** `github.com/hanzoai/datastore-go`
**Version:** 1.0.0
**Package:** `datastore`

## Architecture

### Core Components

```
datastore-go/
├── datastore.go          # Main connection pooling and driver
├── datastore_std.go      # database/sql interface implementation
├── datastore_options.go  # Connection options and DSN parsing
├── datastore_rows.go     # Row handling and scanning
├── conn.go               # Native TCP connection implementation
├── conn_http.go          # HTTP connection implementation
├── conn_batch.go         # Batch insert operations
├── context.go            # Context helpers (WithSettings, WithProgress, etc.)
├── bind.go               # Query parameter binding
├── lib/
│   ├── column/           # Column type implementations
│   ├── proto/            # Wire protocol implementation
│   ├── binary/           # Binary encoding/decoding
│   └── driver/           # Driver interfaces
└── tests/                # Integration tests
```

### Connection Flow

1. **Native Protocol (TCP):** Direct binary protocol connection to port 9000
2. **HTTP Protocol:** HTTP/HTTPS connections to port 8123, useful for proxied environments

### Key Interfaces

- `datastore.Conn` - Native connection interface
- `*sql.DB` - Standard database/sql interface
- `driver.Batch` - Batch insert interface
- `driver.Rows` - Query result interface

## Usage Patterns

### Native Interface
```go
conn, err := datastore.Open(&datastore.Options{
    Addr: []string{"localhost:9000"},
    Auth: datastore.Auth{
        Database: "default",
        Username: "default",
    },
})
defer conn.Close()

// Query
rows, err := conn.Query(ctx, "SELECT * FROM table WHERE id = ?", id)

// Batch insert
batch, err := conn.PrepareBatch(ctx, "INSERT INTO table")
batch.Append(values...)
batch.Send()
```

### database/sql Interface
```go
db := datastore.OpenDB(&datastore.Options{...})
// or
db, err := sql.Open("datastore", "datastore://localhost:9000/default")

// Standard sql operations
rows, err := db.QueryContext(ctx, "SELECT * FROM table")
```

### DSN Format
```
datastore://username:password@host1:9000,host2:9000/database?dial_timeout=30s&compress=lz4
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `datastore.go` | Main connection pool, `Open()` function |
| `datastore_std.go` | `database/sql` driver registration, `OpenDB()` |
| `datastore_options.go` | `Options` struct, `ParseDSN()` |
| `context.go` | Context helpers like `WithSettings()`, `WithProgress()` |
| `conn_batch.go` | Batch insert implementation |
| `bind.go` | Parameter binding logic |
| `lib/column/` | All ClickHouse column type handlers |
| `lib/proto/` | Native protocol encoding |

## Testing

```bash
# Run all tests (requires ClickHouse server or Docker)
go test ./...

# Run specific test
go test -run TestBatch ./tests/

# With verbose output
go test -v ./...
```

Tests use testcontainers-go to automatically spin up ClickHouse instances.

## Build Commands

```bash
# Build
go build ./...

# Lint
golangci-lint run

# Format
go fmt ./...
```

## Important Design Decisions

1. **Connection Pooling:** Native connections are pooled with configurable limits
2. **Compression:** LZ4/ZSTD compression at column level for inserts
3. **Batch Operations:** Columnar batch inserts for high performance
4. **Context Propagation:** Query settings, progress, logs passed via context
5. **Error Handling:** Structured errors with column information

## Migration Notes

This is version 1.0.0, a clean break from clickhouse-go v2:

- Module: `github.com/hanzoai/datastore-go` (was `github.com/ClickHouse/clickhouse-go/v2`)
- Package: `datastore` (was `clickhouse`)
- DSN Scheme: `datastore://` (was `clickhouse://`)
- Driver Name: `datastore` (was `clickhouse`)

## Dependencies

Key external dependencies:
- `github.com/ClickHouse/ch-go` - Low-level protocol encoding
- `github.com/paulmach/orb` - Geospatial types
- `go.opentelemetry.io/otel` - OpenTelemetry integration

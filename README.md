# Datastore

Golang SQL database client for [ClickHouse](https://clickhouse.com/).

## Key features

* Uses ClickHouse native format for optimal performance
* Supports native ClickHouse TCP client-server protocol
* Compatibility with [`database/sql`](#std-databasesql-interface)
* HTTP protocol support for transport
* Marshal rows into structs ([ScanStruct](examples/datastore_api/scan_struct.go), [Select](examples/datastore_api/select_struct.go))
* Unmarshal struct to row ([AppendStruct](benchmark/v2/write-native-struct/main.go))
* Connection pool (for both TCP-Native and HTTP)
* Failover and load balancing
* [Bulk write support](examples/datastore_api/batch.go)
* [PrepareBatch options](#preparebatch-options)
* [AsyncInsert](benchmark/v2/write-async/main.go)
* Named and numeric placeholders support
* LZ4/ZSTD compression support
* External data
* [Query parameters](examples/std/query_parameters.go)

Support for the ClickHouse protocol advanced features using `Context`:

* Query ID
* Quota Key
* Settings
* [Query parameters](examples/datastore_api/query_parameters.go)
* OpenTelemetry
* Execution events:
	* Logs
	* Progress
	* Profile info
	* Profile events

## Install

```sh
go get -u github.com/hanzoai/datastore-go
```

## `datastore` interface (native interface)

```go
conn, err := datastore.Open(&datastore.Options{
	Addr: []string{"127.0.0.1:9000"},
	Auth: datastore.Auth{
		Database: "default",
		Username: "default",
		Password: "",
	},
	DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
		dialCount++
		var d net.Dialer
		return d.DialContext(ctx, "tcp", addr)
	},
	Debug: true,
	Debugf: func(format string, v ...any) {
		fmt.Printf(format+"\n", v...)
	},
	Settings: datastore.Settings{
		"max_execution_time": 60,
	},
	Compression: &datastore.Compression{
		Method: datastore.CompressionLZ4,
	},
	DialTimeout:      time.Second * 30,
	MaxOpenConns:     5,
	MaxIdleConns:     5,
	ConnMaxLifetime:  time.Duration(10) * time.Minute,
	ConnOpenStrategy: datastore.ConnOpenInOrder,
	BlockBufferSize: 10,
	MaxCompressionBuffer: 10240,
	ClientInfo: datastore.ClientInfo{
		Products: []struct {
			Name    string
			Version string
		}{
			{Name: "my-app", Version: "0.1"},
		},
	},
})
if err != nil {
	return err
}
return conn.Ping(context.Background())
```

## `database/sql` interface

### OpenDB

```go
conn := datastore.OpenDB(&datastore.Options{
	Addr: []string{"127.0.0.1:9999"},
	Auth: datastore.Auth{
		Database: "default",
		Username: "default",
		Password: "",
	},
	TLS: &tls.Config{
		InsecureSkipVerify: true,
	},
	Settings: datastore.Settings{
		"max_execution_time": 60,
	},
	DialTimeout: time.Second * 30,
	Compression: &datastore.Compression{
		Method: datastore.CompressionLZ4,
	},
	Debug: true,
	BlockBufferSize: 10,
	MaxCompressionBuffer: 10240,
	ClientInfo: datastore.ClientInfo{
		Products: []struct {
			Name    string
			Version string
		}{
			{Name: "my-app", Version: "0.1"},
		},
	},
})
conn.SetMaxIdleConns(5)
conn.SetMaxOpenConns(10)
conn.SetConnMaxLifetime(time.Hour)
```

### DSN

* hosts - comma-separated list of single address hosts for load-balancing and failover
* username/password - auth credentials
* database - select the current default database
* dial_timeout - a duration string (default 30s)
* connection_open_strategy - random/round_robin/in_order (default in_order)
* debug - enable debug output (boolean value)
* compress - specify the compression algorithm: `none` (default), `zstd`, `lz4`, `lz4hc`, `gzip`, `deflate`, `br`
* compress_level - Level of compression (algorithm-specific)
* block_buffer_size - size of block buffer (default 2)
* read_timeout - a duration string (default 5m)
* max_compression_buffer - max size (bytes) of compression buffer (default 10MiB)
* client_info_product - optional list of product name and version pairs

Example:

```sh
datastore://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&read_timeout=30s&max_execution_time=60
```

### HTTP Support

The native format can be used over the HTTP protocol:

```sh
http://host1:8123,host2:8123/database?dial_timeout=200ms&max_execution_time=60
```

Or using `OpenDB`:

```go
conn := datastore.OpenDB(&datastore.Options{
	Addr: []string{"127.0.0.1:8123"},
	Auth: datastore.Auth{
		Database: "default",
		Username: "default",
		Password: "",
	},
	Settings: datastore.Settings{
		"max_execution_time": 60,
	},
	DialTimeout: 30 * time.Second,
	Compression: &datastore.Compression{
		Method: datastore.CompressionLZ4,
	},
	Protocol: datastore.HTTP,
})
```

## Compression

ZSTD/LZ4 compression is supported over native and http protocols. This is performed column by column at a block level and is only used for inserts.

## TLS/SSL

Set a non-nil `tls.Config` pointer in the Options struct to establish a secure connection:

```go
conn := datastore.OpenDB(&datastore.Options{
	...
	TLS: &tls.Config{
		InsecureSkipVerify: false,
	},
	...
})
```

### HTTPS

Use `https` in your DSN string:

```sh
https://host1:8443,host2:8443/database?dial_timeout=200ms&max_execution_time=60
```

## Async insert

[Async insert](https://clickhouse.com/docs/optimize/asynchronous-inserts) is supported via `WithAsync()` helper:

```go
ctx := datastore.Context(ctx, datastore.WithAsync(true))
```

## PrepareBatch options

Available options:
- [WithReleaseConnection](examples/datastore_api/batch_release_connection.go) - return connection to pool after PrepareBatch

## Examples

### native interface

* [batch](examples/datastore_api/batch.go)
* [batch with release connection](examples/datastore_api/batch_release_connection.go)
* [native async insert](examples/datastore_api/async_native.go)
* [http async insert](examples/datastore_api/async_http.go)
* [batch struct](examples/datastore_api/append_struct.go)
* [columnar](examples/datastore_api/columnar_insert.go)
* [scan struct](examples/datastore_api/scan_struct.go)
* [query parameters](examples/datastore_api/query_parameters.go)
* [bind params](examples/datastore_api/bind.go)
* [client info](examples/datastore_api/client_info.go)

### std `database/sql` interface

* [batch](examples/std/batch.go)
* [native async insert](examples/std/async_native.go)
* [http async insert](examples/std/async_http.go)
* [open db](examples/std/connect.go)
* [query parameters](examples/std/query_parameters.go)
* [bind params](examples/std/bind.go)
* [client info](examples/std/client_info.go)

## License

Apache License 2.0

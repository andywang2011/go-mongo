# go-mongo

Thin, testable MongoDB client for Go, built on the [official mongo-driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo).

```bash
go get github.com/andywang2011/go-mongo@v0.2.0
```

## Design

- No viper, no global singletons, no `log.Fatal` inside the library.
- One `Client` per process; share it across goroutines (driver connection pool is thread-safe).
- Credentials and topology come from your app via `Config`.
- Call `Close(ctx)` on shutdown (use a timeout context in production).

## Quick start

```go
cfg := mongo.Config{
    Address:  "127.0.0.1:27017",
    Database: "myapp",
}
if err := cfg.Validate(); err != nil { /* handle */ }

client, err := mongo.New(cfg)
if err != nil { /* handle */ }
defer client.Close(context.Background())

coll := client.Database().Collection("users")
// ...
```

`ConfigDefault()` sets pool and timeout defaults only; you must still set `Database` (and hosts).

## Cluster topologies

| Mode | Config |
|------|--------|
| Single node | `Address: "127.0.0.1:27017"` or `DirectConnection: true` for dev |
| Replica set | `Addresses: []string{"h1:27017","h2:27017"}`, `ReplicaSet: "rs0"` |
| Atlas / SRV | `Address: "mongodb+srv://cluster.example.net"` (full URI; `Addresses` ignored) |
| Full URI | `Address: "mongodb://user:pass@host1,host2/db?replicaSet=rs0"` |

Auth is embedded in the URI when `Username` / `Password` are set (`authSource` defaults to `admin`).

## Connection pool (high concurrency)

| Field | Default | Notes |
|-------|---------|--------|
| `MaxPoolSize` | 100 | Set ≥ expected concurrent DB operations |
| `MinPoolSize` | 0 | Warm connections; production often `MaxPoolSize/4` |
| `MaxConnIdleTime` | 0 (driver default) | Close idle sockets |
| `MaxConnecting` | 0 (driver default) | Cap concurrent handshakes |
| `ConnectTimeout` | 10s | Connect + initial ping |
| `ServerSelectionTO` | 5s | Replica set / SRV selection |
| `SocketTimeout` | 0 | Per-operation socket timeout when set |

Do not call `New` per HTTP request; reuse one client for the process lifetime.

## Read preference & write concern

- `ReadPreference`: `primary`, `primaryPreferred`, `secondary`, `secondaryPreferred`, `nearest`
- `WriteConcern`: optional `WriteConcernConfig` (`W`, `J`, `WTimeout`)

## TLS

Set `TLSEnable`, `TLSCA`, and optionally `TLSCert` / `TLSKey` for mutual TLS (PEM strings, same style as go-kafka).

## API

| Method | Description |
|--------|-------------|
| `New(cfg)` | Connect, ping, return `Client` |
| `Database()` | Default `*mongo.Database` |
| `Ping(ctx)` | Health check |
| `RawClient()` | Underlying `*mongo.Client` |
| `Close(ctx)` | Disconnect |

`Conn.DB` is kept for backward compatibility; prefer `Database()`.

## Examples

- [examples/basic](examples/basic/main.go) — insert / find / close
- [examples/replica](examples/replica/main.go) — replica set hosts + pool tuning

## Limitations

- No built-in config file loading (use viper/Nacos in your app).
- Integration tests against a real MongoDB are optional (`-tags=integration`).

---

## 中文摘要

面向应用框架的 MongoDB 封装：配置由调用方注入，进程内单例共享连接池，支持副本集/SRV/TLS/读写偏好/WriteConcern，英文 godoc，可 mock 的 `Client` 接口。

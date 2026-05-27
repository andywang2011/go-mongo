# Changelog

## Unreleased

### Added

- `doc.go`, `errors.go` (`ErrInvalidConfig`)
- `BuildURI` in `uri.go`; `applyClientOptions` in `options.go`
- Cluster: `Addresses`, `ReplicaSet`, `ReadPreference`, `DirectConnection`
- `WriteConcernConfig`, TLS PEM fields, `MaxConnIdleTime`, `SocketTimeout`, `MaxConnecting`
- `ConfigDefault()`, injectable `connectMongo` / `pingAfterConnect` for tests
- `Client.Ping`, `Conn.RawClient()`; `Client()` deprecated alias
- English README, examples (`basic`, `replica`), `goleak` in `TestMain`
- CI coverage gate ≥ 85%

### Changed

- **Breaking:** default `MaxPoolSize` is now **100** (was 10)
- Auth via connection URI (no duplicate `SetAuth`)
- Godoc and package docs in English

## v0.1.0

- Initial thin wrapper around mongo-driver with `Config`, `New`, `Close`, mock-friendly `Client`

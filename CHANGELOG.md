# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Breaking

- Migrated to `go.mongodb.org/mongo-driver/v2`. `Client.Database()`, `Conn.RawClient()` and `Conn.DB` now expose v2 `*mongo.Database` / `*mongo.Client`; update your imports to `go.mongodb.org/mongo-driver/v2/...` (including `bson`).
- `Config.SocketTimeout` now maps to the driver's `ClientOptions.Timeout` (v2 removed `SocketTimeout`).
- `WriteConcernConfig.WTimeout` is still validated but has no effect (v2 removed `wtimeout`).

### Unchanged

- `New(ctx, cfg)`, `Close`, `Ping`, `BuildURI`, `RedactedURI`, `Config`/`ConfigDefault` keep their signatures.

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



### Changed

- **Breaking:** removed deprecated `Conn.Client()`; use `RawClient()`.
- `Config.Validate` now validates `WriteConcern` (W, WTimeout, W=0 with J).
- Failed ping cleanup disconnects with a bounded timeout context.
- **Breaking:** `New` takes a `context.Context` first (combined with `ConnectTimeout`); `Conn` is safe for concurrent `Close`.

### Added

- `RedactedURI(cfg)` / `Config.RedactedURI()` for log-safe connection URIs.

## [0.1.0] - 2026-11-07

### Added

- Config-driven Mongo client with replica set, SRV, and TLS options.

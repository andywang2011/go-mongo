// Package mongo provides a thin, testable wrapper around the official MongoDB Go driver.
//
// Design constraints (for use in application frameworks):
//   - No viper, no global singletons, no log.Fatal inside the library.
//   - Credentials and topology settings are supplied by the caller via Config.
//   - Use New once per process and share the returned Client across goroutines (driver pool is thread-safe).
//   - Always call Close with a context on shutdown.
//
// Cluster topologies are expressed via Config (hosts, replica set, mongodb+srv, read preference).
// Connection pool sizing (MaxPoolSize, MinPoolSize, idle timeouts) is configured on Config for high concurrency.
package mongo

package mongo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

// Read preference mode names for Config.ReadPreference.
const (
	ReadPrimary            = "primary"
	ReadPrimaryPreferred   = "primaryPreferred"
	ReadSecondary          = "secondary"
	ReadSecondaryPreferred = "secondaryPreferred"
	ReadNearest            = "nearest"
)

// WriteConcernConfig optional write concern (nil uses driver default).
type WriteConcernConfig struct {
	// W is the write acknowledgment (e.g. 1, "majority"). Empty defaults to 1 when J or WTimeout is set.
	W string
	// J requests journal acknowledgment.
	J bool
	// WTimeout is the max time to wait for W.
	WTimeout time.Duration
}

// Config holds MongoDB client settings (filled by the application, not the library).
type Config struct {
	// Address is host:port, comma-separated hosts, or a full mongodb:// / mongodb+srv:// URI.
	Address string
	// Addresses is an explicit host list; used when Address is not a full URI (overrides host part).
	Addresses []string

	// Username is embedded in the connection URI when set.
	Username string
	// Password is embedded in the connection URI when set.
	Password string
	// Database is the default database name (required).
	Database string
	// AuthSource is the authentication database (default "admin" when Username is set).
	AuthSource string

	// ReplicaSet sets the replicaSet query parameter.
	ReplicaSet string
	// ReadPreference is primary, primaryPreferred, secondary, secondaryPreferred, or nearest.
	ReadPreference string
	// DirectConnection sets directConnection=true (useful for single-node development).
	DirectConnection bool

	// WriteConcern optional write concern; nil keeps driver default.
	WriteConcern *WriteConcernConfig

	// MaxPoolSize is the maximum connections in the pool (default 100).
	MaxPoolSize uint64
	// MinPoolSize is the minimum warm connections in the pool.
	MinPoolSize uint64
	// MaxConnIdleTime closes idle connections after this duration.
	MaxConnIdleTime time.Duration
	// MaxConnecting limits connections being established concurrently.
	MaxConnecting uint64

	// ConnectTimeout is used for initial connect and ping.
	ConnectTimeout time.Duration
	// ServerSelectionTO is server selection timeout.
	ServerSelectionTO time.Duration
	// SocketTimeout is the socket timeout for operations.
	SocketTimeout time.Duration

	// TLSEnable enables TLS when TLSCA is set (and optional client cert/key for mutual TLS).
	TLSEnable bool
	// TLSInsecureSkipVerify skips server certificate verification (dev/debug only).
	TLSInsecureSkipVerify bool
	// TLSCA is the CA PEM to verify the server certificate.
	TLSCA string
	// TLSCert is the client certificate PEM (mutual TLS).
	TLSCert string
	// TLSKey is the client private key PEM.
	TLSKey string
}

const (
	defaultHost              = "127.0.0.1:27017"
	defaultMaxPoolSize       = 100
	defaultConnectTimeout    = 10 * time.Second
	defaultServerSelectionTO = 5 * time.Second
	defaultAuthSource        = "admin"
)

// ConfigDefault returns pool and timeout defaults only; Database and hosts must still be set before Validate.
func ConfigDefault() Config {
	return Config{
		MaxPoolSize:       defaultMaxPoolSize,
		ConnectTimeout:    defaultConnectTimeout,
		ServerSelectionTO: defaultServerSelectionTO,
	}
}

// Validate checks required fields and applies defaults.
func (c *Config) Validate() error {
	if c.Database == "" {
		return fmt.Errorf("%w: database is required", ErrInvalidConfig)
	}
	if c.MaxPoolSize == 0 {
		c.MaxPoolSize = defaultMaxPoolSize
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = defaultConnectTimeout
	}
	if c.ServerSelectionTO == 0 {
		c.ServerSelectionTO = defaultServerSelectionTO
	}
	if c.ReadPreference != "" {
		if _, err := parseReadPreference(c.ReadPreference); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
		}
	}
	if c.WriteConcern != nil {
		if err := c.WriteConcern.validate(); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
		}
	}
	if c.TLSEnable {
		if c.TLSCA == "" {
			return fmt.Errorf("%w: TLSCA is required when TLSEnable is true", ErrInvalidConfig)
		}
	}
	return nil
}

func (w *WriteConcernConfig) validate() error {
	if w == nil {
		return nil
	}
	return nil
}

// BuildTLS returns *tls.Config when TLS is enabled and CA material is present.
func (c Config) BuildTLS() (*tls.Config, error) {
	if !c.TLSEnable || c.TLSCA == "" {
		return nil, nil
	}
	cfg := &tls.Config{InsecureSkipVerify: c.TLSInsecureSkipVerify}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(c.TLSCA)) {
		return nil, fmt.Errorf("mongo tls ca parse failed")
	}
	cfg.RootCAs = pool
	if c.TLSCert != "" && c.TLSKey != "" {
		cert, err := tls.X509KeyPair([]byte(c.TLSCert), []byte(c.TLSKey))
		if err != nil {
			return nil, fmt.Errorf("mongo tls client cert: %w", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	return cfg, nil
}

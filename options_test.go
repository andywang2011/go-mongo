package mongo

import (
	"testing"
	"time"
)

func Test_applyClientOptions_poolAndTimeouts(t *testing.T) {
	cfg := Config{
		Database:          "app",
		Address:           "127.0.0.1:27017",
		MaxPoolSize:       50,
		MinPoolSize:       5,
		MaxConnIdleTime:   time.Minute,
		MaxConnecting:     3,
		ConnectTimeout:    2 * time.Second,
		ServerSelectionTO: 3 * time.Second,
		SocketTimeout:     4 * time.Second,
		DirectConnection:  true,
		ReplicaSet:        "rs0",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	opts, err := applyClientOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opts.MaxPoolSize == nil || *opts.MaxPoolSize != 50 {
		t.Fatalf("MaxPoolSize %+v", opts.MaxPoolSize)
	}
}

func Test_applyClientOptions_readPreference(t *testing.T) {
	cfg := Config{
		Database:          "app",
		Address:           "127.0.0.1:27017",
		ReadPreference:    ReadSecondary,
		MaxPoolSize:       10,
		ConnectTimeout:    time.Second,
		ServerSelectionTO: time.Second,
	}
	opts, err := applyClientOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opts.ReadPreference == nil {
		t.Fatal("expected read preference")
	}
}

func Test_applyClientOptions_writeConcern(t *testing.T) {
	cfg := Config{
		Database:          "app",
		Address:           "127.0.0.1:27017",
		MaxPoolSize:       10,
		ConnectTimeout:    time.Second,
		ServerSelectionTO: time.Second,
		WriteConcern: &WriteConcernConfig{
			W: "majority",
			J: true,
		},
	}
	opts, err := applyClientOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opts.WriteConcern == nil {
		t.Fatal("expected write concern")
	}
}

func Test_buildWriteConcern_invalidW(t *testing.T) {
	_, err := buildWriteConcern(&WriteConcernConfig{W: "invalid"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_parseReadPreference_unknown(t *testing.T) {
	_, err := parseReadPreference("bogus")
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_applyClientOptions_tls(t *testing.T) {
	caPEM, certPEM, keyPEM := testCertPEM(t)
	cfg := Config{
		Database:          "app",
		Address:           "127.0.0.1:27017",
		MaxPoolSize:       10,
		ConnectTimeout:    time.Second,
		ServerSelectionTO: time.Second,
		TLSEnable:         true,
		TLSCA:             caPEM,
		TLSCert:           certPEM,
		TLSKey:            keyPEM,
	}
	opts, err := applyClientOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opts.TLSConfig == nil {
		t.Fatal("expected TLS config")
	}
}

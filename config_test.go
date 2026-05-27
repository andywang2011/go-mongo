package mongo

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	cfg := Config{Database: "app"}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.MaxPoolSize != defaultMaxPoolSize {
		t.Fatalf("MaxPoolSize = %d", cfg.MaxPoolSize)
	}
}

func TestConfig_Validate_missingDatabase(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_customTimeouts(t *testing.T) {
	cfg := Config{
		Database:          "db",
		ConnectTimeout:    3 * time.Second,
		ServerSelectionTO: 2 * time.Second,
		MaxPoolSize:       5,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.MaxPoolSize != 5 {
		t.Fatal("MaxPoolSize overwritten")
	}
}

func TestConfig_Validate_readPreference(t *testing.T) {
	cfg := Config{Database: "db", ReadPreference: "nope"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_tlsRequiresCA(t *testing.T) {
	cfg := Config{Database: "db", TLSEnable: true}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigDefault(t *testing.T) {
	c := ConfigDefault()
	if c.MaxPoolSize != defaultMaxPoolSize {
		t.Fatalf("MaxPoolSize %d", c.MaxPoolSize)
	}
}

func TestConfig_BuildTLS_disabled(t *testing.T) {
	tlsCfg, err := Config{Database: "db"}.BuildTLS()
	if err != nil || tlsCfg != nil {
		t.Fatalf("got %v err %v", tlsCfg, err)
	}
}

func TestConfig_BuildTLS_badCA(t *testing.T) {
	_, err := Config{TLSEnable: true, TLSCA: "not-pem"}.BuildTLS()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_BuildTLS_withCA(t *testing.T) {
	caPEM, _, _ := testCertPEM(t)
	cfg := Config{TLSEnable: true, TLSCA: caPEM}
	tlsCfg, err := cfg.BuildTLS()
	if err != nil || tlsCfg == nil || tlsCfg.RootCAs == nil {
		t.Fatalf("BuildTLS: %v cfg %v", err, tlsCfg)
	}
}

func TestConfig_BuildTLS_mutual(t *testing.T) {
	caPEM, certPEM, keyPEM := testCertPEM(t)
	cfg := Config{
		TLSEnable: true,
		TLSCA:     caPEM,
		TLSCert:   certPEM,
		TLSKey:    keyPEM,
	}
	tlsCfg, err := cfg.BuildTLS()
	if err != nil || len(tlsCfg.Certificates) != 1 {
		t.Fatalf("BuildTLS: %v certs %d", err, len(tlsCfg.Certificates))
	}
}

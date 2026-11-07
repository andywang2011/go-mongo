package mongo

import (
	"context"
	"testing"
	"time"
)

func Test_defaultConnectMongo(t *testing.T) {
	t.Run("ok lazy connect", func(t *testing.T) {
		c, err := defaultConnectMongo(context.Background(), validTestConfig())
		if err != nil {
			t.Fatal(err)
		}
		_ = c.Disconnect(context.Background())
	})
	t.Run("cancelled ctx", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := defaultConnectMongo(ctx, validTestConfig()); err == nil {
			t.Fatal("expected ctx error")
		}
	})
	t.Run("invalid read preference", func(t *testing.T) {
		cfg := validTestConfig()
		cfg.ReadPreference = "bogus"
		if _, err := defaultConnectMongo(context.Background(), cfg); err == nil {
			t.Fatal("expected error")
		}
	})
}

func Test_parseReadPreference_allModes(t *testing.T) {
	for _, m := range []string{"", ReadPrimary, ReadPrimaryPreferred, ReadSecondary, ReadSecondaryPreferred, ReadNearest} {
		if rp, err := parseReadPreference(m); err != nil || rp == nil {
			t.Errorf("mode %q: %v", m, err)
		}
	}
}

func Test_buildWriteConcern_v2(t *testing.T) {
	if wc, err := buildWriteConcern(nil); wc != nil || err != nil {
		t.Fatal("nil config must yield nil")
	}
	if wc, err := buildWriteConcern(&WriteConcernConfig{}); wc != nil || err != nil {
		t.Fatalf("empty config must yield nil, got %v %v", wc, err)
	}
	wc, err := buildWriteConcern(&WriteConcernConfig{W: "3"})
	if err != nil || wc.W != 3 {
		t.Fatalf("W=3: %v %v", wc, err)
	}
	wc, err = buildWriteConcern(&WriteConcernConfig{WTimeout: time.Second})
	if err != nil || wc.W != 1 {
		t.Fatalf("WTimeout-only defaults W=1: %v %v", wc, err)
	}
	wc, err = buildWriteConcern(&WriteConcernConfig{J: true})
	if err != nil || wc.W != 1 || wc.Journal == nil || !*wc.Journal {
		t.Fatalf("J: %v %v", wc, err)
	}
}

func Test_applyClientOptions_socketTimeoutMapsToTimeout(t *testing.T) {
	cfg := validTestConfig()
	cfg.SocketTimeout = 7 * time.Second
	opts, err := applyClientOptions(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if opts.Timeout == nil || *opts.Timeout != 7*time.Second {
		t.Fatalf("Timeout = %v", opts.Timeout)
	}
	cfg.SocketTimeout = 0
	if opts, _ = applyClientOptions(cfg); opts.Timeout != nil {
		t.Fatal("Timeout must be unset when SocketTimeout is 0")
	}
}

func Test_applyClientOptions_errors(t *testing.T) {
	cfg := validTestConfig()
	cfg.ReadPreference = "bogus"
	if _, err := applyClientOptions(cfg); err == nil {
		t.Error("expected read preference error")
	}
	cfg = validTestConfig()
	cfg.WriteConcern = &WriteConcernConfig{W: "x"}
	if _, err := applyClientOptions(cfg); err == nil {
		t.Error("expected write concern error")
	}
	cfg = validTestConfig()
	cfg.TLSEnable = true
	cfg.TLSCA = "not pem"
	if _, err := applyClientOptions(cfg); err == nil {
		t.Error("expected TLS error")
	}
}

func TestConfig_BuildTLS_badClientCert(t *testing.T) {
	caPEM, certPEM, _ := testCertPEM(t)
	cfg := Config{TLSEnable: true, TLSCA: caPEM, TLSCert: certPEM, TLSKey: "bad"}
	if _, err := cfg.BuildTLS(); err == nil {
		t.Fatal("expected client cert error")
	}
	cfg = Config{TLSEnable: true, TLSCA: "bad"}
	if _, err := cfg.BuildTLS(); err == nil {
		t.Fatal("expected CA parse error")
	}
}

func Test_hostList(t *testing.T) {
	if got := hostList(Config{Address: "mongodb://h:1"}); got != "mongodb://h:1" {
		t.Errorf("got %q", got)
	}
	if got := hostList(Config{Addresses: []string{"a:1", "b:2"}}); got != "a:1,b:2" {
		t.Errorf("got %q", got)
	}
	if got := hostList(Config{}); got != defaultHost {
		t.Errorf("got %q", got)
	}
	if got := hostList(Config{Address: "x:1"}); got != "x:1" {
		t.Errorf("got %q", got)
	}
}

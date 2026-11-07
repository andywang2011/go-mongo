package mongo

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

func parseReadPreference(name string) (*readpref.ReadPref, error) {
	switch name {
	case "", ReadPrimary:
		return readpref.Primary(), nil
	case ReadPrimaryPreferred:
		return readpref.PrimaryPreferred(), nil
	case ReadSecondary:
		return readpref.Secondary(), nil
	case ReadSecondaryPreferred:
		return readpref.SecondaryPreferred(), nil
	case ReadNearest:
		return readpref.Nearest(), nil
	default:
		return nil, fmt.Errorf("unknown read preference %q", name)
	}
}

func buildWriteConcern(wc *WriteConcernConfig) (*writeconcern.WriteConcern, error) {
	if wc == nil {
		return nil, nil
	}

	result := &writeconcern.WriteConcern{}
	// WTimeout is not supported by mongo-driver v2's WriteConcern; use
	// Config.SocketTimeout (client-side operation timeout) to bound writes.
	if wc.W != "" {
		if wc.W == "majority" {
			result.W = "majority"
		} else if n, err := strconv.Atoi(wc.W); err == nil {
			result.W = n
		} else {
			return nil, fmt.Errorf("invalid write concern W %q", wc.W)
		}
	} else if wc.J || wc.WTimeout > 0 {
		result.W = 1
	}
	if wc.J {
		journal := true
		result.Journal = &journal
	}
	if result.W == nil && result.Journal == nil {
		return nil, nil
	}
	return result, nil
}

// applyClientOptions maps Config to mongo-driver ClientOptions (URI carries auth and topology query params).
func applyClientOptions(cfg Config) (*options.ClientOptions, error) {
	uri := BuildURI(cfg)
	opts := options.Client().ApplyURI(uri)

	opts.SetMaxPoolSize(cfg.MaxPoolSize)
	if cfg.MinPoolSize > 0 {
		opts.SetMinPoolSize(cfg.MinPoolSize)
	}
	if cfg.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	}
	if cfg.MaxConnecting > 0 {
		opts.SetMaxConnecting(cfg.MaxConnecting)
	}
	opts.SetConnectTimeout(cfg.ConnectTimeout)
	opts.SetServerSelectionTimeout(cfg.ServerSelectionTO)
	if cfg.SocketTimeout > 0 {
		// v2 removed SocketTimeout; Timeout is the client-side per-operation limit.
		opts.SetTimeout(cfg.SocketTimeout)
	}
	if cfg.DirectConnection {
		opts.SetDirect(true)
	}
	if cfg.ReplicaSet != "" {
		opts.SetReplicaSet(cfg.ReplicaSet)
	}

	if cfg.ReadPreference != "" {
		rp, err := parseReadPreference(cfg.ReadPreference)
		if err != nil {
			return nil, err
		}
		opts.SetReadPreference(rp)
	}

	wc, err := buildWriteConcern(cfg.WriteConcern)
	if err != nil {
		return nil, err
	}
	if wc != nil {
		opts.SetWriteConcern(wc)
	}

	if cfg.TLSEnable {
		tlsCfg, err := cfg.BuildTLS()
		if err != nil {
			return nil, err
		}
		if tlsCfg != nil {
			opts.SetTLSConfig(tlsCfg)
		}
	}

	return opts, nil
}

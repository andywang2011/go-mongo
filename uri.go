package mongo

import (
	"net/url"
	"strings"
)

func isFullURI(s string) bool {
	return strings.HasPrefix(s, "mongodb://") || strings.HasPrefix(s, "mongodb+srv://")
}

func hostList(cfg Config) string {
	if isFullURI(cfg.Address) {
		return cfg.Address
	}
	if len(cfg.Addresses) > 0 {
		return strings.Join(cfg.Addresses, ",")
	}
	addr := cfg.Address
	if addr == "" {
		addr = defaultHost
	}
	return addr
}

// BuildURI builds a MongoDB connection URI from cfg (call Validate first for defaults).
func BuildURI(cfg Config) string {
	if isFullURI(cfg.Address) {
		return cfg.Address
	}

	hosts := hostList(cfg)
	u := &url.URL{
		Scheme: "mongodb",
		Host:   hosts,
		Path:   "/" + cfg.Database,
	}

	if cfg.Username != "" {
		u.User = url.UserPassword(cfg.Username, cfg.Password)
	}

	q := url.Values{}
	if cfg.AuthSource != "" {
		q.Set("authSource", cfg.AuthSource)
	} else if cfg.Username != "" {
		q.Set("authSource", defaultAuthSource)
	}
	if cfg.ReplicaSet != "" {
		q.Set("replicaSet", cfg.ReplicaSet)
	}
	if cfg.ReadPreference != "" {
		q.Set("readPreference", cfg.ReadPreference)
	}
	if cfg.DirectConnection {
		q.Set("directConnection", "true")
	}
	if len(q) > 0 {
		u.RawQuery = q.Encode()
	}
	return u.String()
}

// redactedPassword replaces a password in log-safe URIs.
const redactedPassword = "xxxxx"

// RedactedURI returns BuildURI(cfg) with the password replaced by a placeholder, safe for logs.
func RedactedURI(cfg Config) string {
	return redactURIPassword(BuildURI(cfg))
}

// RedactedURI returns the connection URI with the password redacted.
func (c Config) RedactedURI() string {
	return RedactedURI(c)
}

// redactURIPassword masks the password in the userinfo section of a mongodb URI.
func redactURIPassword(uri string) string {
	schemeEnd := strings.Index(uri, "://")
	if schemeEnd < 0 {
		return uri
	}
	start := schemeEnd + len("://")
	rest := uri[start:]
	authEnd := strings.IndexAny(rest, "/?#")
	if authEnd < 0 {
		authEnd = len(rest)
	}
	authority := rest[:authEnd]
	at := strings.LastIndex(authority, "@")
	if at < 0 {
		return uri
	}
	userinfo := authority[:at]
	colon := strings.Index(userinfo, ":")
	if colon < 0 {
		return uri
	}
	return uri[:start] + userinfo[:colon] + ":" + redactedPassword + rest[at:]
}

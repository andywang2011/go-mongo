package mongo

import (
	"net/url"
	"strings"
	"testing"
)

func TestBuildURI_simple(t *testing.T) {
	got := BuildURI(Config{Address: "127.0.0.1:27017", Database: "app"})
	want := "mongodb://127.0.0.1:27017/app"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBuildURI_fullURIPassthrough(t *testing.T) {
	full := "mongodb://host:27017/mydb"
	if got := BuildURI(Config{Address: full, Database: "ignored"}); got != full {
		t.Fatalf("got %q", got)
	}
}

func TestBuildURI_srvPassthrough(t *testing.T) {
	full := "mongodb+srv://cluster.example.com"
	if got := BuildURI(Config{Address: full, Database: "db"}); got != full {
		t.Fatalf("got %q", got)
	}
}

func TestBuildURI_defaultHost(t *testing.T) {
	got := BuildURI(Config{Database: "app"})
	if !strings.HasPrefix(got, "mongodb://127.0.0.1:27017/app") {
		t.Fatalf("got %q", got)
	}
}

func TestBuildURI_addresses(t *testing.T) {
	got := BuildURI(Config{
		Addresses:  []string{"h1:27017", "h2:27017"},
		Database:   "app",
		ReplicaSet: "rs0",
	})
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "h1:27017,h2:27017" {
		t.Fatalf("host %q", u.Host)
	}
	if u.Query().Get("replicaSet") != "rs0" {
		t.Fatalf("query %v", u.Query())
	}
}

func TestBuildURI_authAndReadPref(t *testing.T) {
	got := BuildURI(Config{
		Address:          "127.0.0.1:27017",
		Database:         "app",
		Username:         "u",
		Password:         "p",
		AuthSource:       "admin",
		ReadPreference:   ReadPrimaryPreferred,
		DirectConnection: true,
	})
	if !strings.Contains(got, "u:p@") {
		t.Fatalf("missing auth: %q", got)
	}
	if !strings.Contains(got, "readPreference=primaryPreferred") {
		t.Fatalf("missing readPreference: %q", got)
	}
	if !strings.Contains(got, "directConnection=true") {
		t.Fatalf("missing directConnection: %q", got)
	}
}

func TestRedactedURI_builtURI(t *testing.T) {
	cfg := Config{Address: "h:27017", Database: "app", Username: "u", Password: "s3cr3t/@p"}
	got := RedactedURI(cfg)
	if strings.Contains(got, "s3cr3t") {
		t.Fatalf("password leaked: %s", got)
	}
	if !strings.HasPrefix(got, "mongodb://u:xxxxx@h:27017/app") {
		t.Fatalf("got %q", got)
	}
	if cfg.RedactedURI() != got {
		t.Fatal("method must match function")
	}
}

func TestRedactURIPassword(t *testing.T) {
	tests := []struct{ in, want string }{
		{"mongodb://u:p@h1:1,h2:2/db?x=a@b", "mongodb://u:xxxxx@h1:1,h2:2/db?x=a@b"},
		{"mongodb+srv://u:p%40ss@cluster.example.com", "mongodb+srv://u:xxxxx@cluster.example.com"},
		{"mongodb://u:p@w@h/db", "mongodb://u:xxxxx@h/db"},
		{"mongodb://user@h/db", "mongodb://user@h/db"},
		{"mongodb://h/db", "mongodb://h/db"},
		{"mongodb://h/db?authSource=a@b:c", "mongodb://h/db?authSource=a@b:c"},
		{"not-a-uri", "not-a-uri"},
		{"mongodb://u:p@h", "mongodb://u:xxxxx@h"},
	}
	for _, tt := range tests {
		if got := redactURIPassword(tt.in); got != tt.want {
			t.Errorf("redact(%q)=%q want %q", tt.in, got, tt.want)
		}
	}
}

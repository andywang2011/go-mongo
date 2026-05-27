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
		Addresses: []string{"h1:27017", "h2:27017"},
		Database:  "app",
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
		Address:        "127.0.0.1:27017",
		Database:       "app",
		Username:       "u",
		Password:       "p",
		AuthSource:     "admin",
		ReadPreference: ReadPrimaryPreferred,
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

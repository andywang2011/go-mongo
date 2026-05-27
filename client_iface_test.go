package mongo

import (
	"context"
	"testing"
)

func TestConn_implementsClient(t *testing.T) {
	var _ Client = (*Conn)(nil)
}

func TestConn_Database_prefersMethod(t *testing.T) {
	c := &Conn{DB: nil}
	if c.Database() != nil {
		t.Fatal("nil DB")
	}
}

func TestFakeClient_ping(t *testing.T) {
	f := &fakeClient{}
	if err := f.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}

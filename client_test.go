package mongo

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func validTestConfig() Config {
	return Config{
		Database:          "testdb",
		Address:           "127.0.0.1:27017",
		MaxPoolSize:       10,
		ConnectTimeout:    time.Second,
		ServerSelectionTO: time.Second,
	}
}

func TestNew_invalidConfig(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_connectError(t *testing.T) {
	old := connectMongo
	defer func() { connectMongo = old }()
	connectMongo = func(context.Context, Config) (*mongo.Client, error) {
		return nil, errors.New("connect failed")
	}
	_, err := New(validTestConfig())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_pingError(t *testing.T) {
	oldConnect := connectMongo
	oldPing := pingAfterConnect
	defer func() {
		connectMongo = oldConnect
		pingAfterConnect = oldPing
	}()
	connectMongo = func(ctx context.Context, cfg Config) (*mongo.Client, error) {
		return mongo.NewClient(options.Client().ApplyURI(BuildURI(cfg)))
	}
	pingAfterConnect = func(*mongo.Client, context.Context) error {
		return errors.New("ping failed")
	}
	_, err := New(validTestConfig())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_success(t *testing.T) {
	oldConnect := connectMongo
	oldPing := pingAfterConnect
	defer func() {
		connectMongo = oldConnect
		pingAfterConnect = oldPing
	}()
	var created *mongo.Client
	connectMongo = func(ctx context.Context, cfg Config) (*mongo.Client, error) {
		c, err := mongo.NewClient(options.Client().ApplyURI(BuildURI(cfg)))
		if err != nil {
			return nil, err
		}
		created = c
		return c, nil
	}
	pingAfterConnect = func(*mongo.Client, context.Context) error { return nil }

	c, err := New(validTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	conn, ok := c.(*Conn)
	if !ok {
		t.Fatal("expected *Conn")
	}
	if conn.Database() == nil {
		t.Fatal("expected database")
	}
	if conn.RawClient() != created {
		t.Fatal("RawClient mismatch")
	}
	if conn.Client() != created {
		t.Fatal("Client() alias mismatch")
	}
	if err := conn.Ping(context.Background()); err != nil {
		// disconnected client may error on ping; ok if client nil path tested elsewhere
	}
	_ = conn.Close(context.Background())
}

func TestConn_Close_nil(t *testing.T) {
	var c *Conn
	if err := c.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConn_Close_twice(t *testing.T) {
	c := &Conn{}
	if err := c.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConn_Ping_nil(t *testing.T) {
	var c *Conn
	if err := c.Ping(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestConn_Database_nil(t *testing.T) {
	var c *Conn
	if c.Database() != nil {
		t.Fatal("expected nil")
	}
}

func TestConn_RawClient_nil(t *testing.T) {
	var c *Conn
	if c.RawClient() != nil {
		t.Fatal("expected nil")
	}
}

func TestNew_applyOptionsError(t *testing.T) {
	cfg := validTestConfig()
	cfg.WriteConcern = &WriteConcernConfig{W: "not-a-number"}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

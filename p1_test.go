package mongo

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func swapSeams(t *testing.T) {
	t.Helper()
	oc, op, od := connectMongo, pingAfterConnect, disconnectTimeout
	t.Cleanup(func() { connectMongo, pingAfterConnect, disconnectTimeout = oc, op, od })
}

func TestNew_connectErrorWithClientDisconnects(t *testing.T) {
	swapSeams(t)
	connectMongo = func(_ context.Context, cfg Config) (*mongo.Client, error) {
		c, err := mongo.Connect(options.Client().ApplyURI(BuildURI(cfg)))
		if err != nil {
			return nil, err
		}
		return c, errors.New("connect failed")
	}
	if _, err := New(context.Background(), validTestConfig()); err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_cleanupUsesFreshContextWhenConnectCtxExpired(t *testing.T) {
	swapSeams(t)
	var created *mongo.Client
	connectMongo = func(_ context.Context, cfg Config) (*mongo.Client, error) {
		c, err := mongo.Connect(options.Client().ApplyURI(BuildURI(cfg)))
		created = c
		return c, err
	}
	pingAfterConnect = func(_ *mongo.Client, ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}
	cfg := validTestConfig()
	cfg.ConnectTimeout = 20 * time.Millisecond
	_, err := New(context.Background(), cfg)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	// Second Disconnect on the already-disconnected client reports the driver's
	// "already disconnected" error, proving the cleanup Disconnect ran.
	if created == nil || created.Disconnect(context.Background()) == nil {
		t.Fatal("expected cleanup to have disconnected the client")
	}
}

func TestNew_parentCtxCancelled(t *testing.T) {
	swapSeams(t)
	connectMongo = func(ctx context.Context, _ Config) (*mongo.Client, error) { return nil, ctx.Err() }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := New(ctx, validTestConfig()); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestNew_nilCtxTolerated(t *testing.T) {
	swapSeams(t)
	connectMongo = func(context.Context, Config) (*mongo.Client, error) { return nil, errors.New("x") }
	//lint:ignore SA1012 intentional nil context in test
	if _, err := New(nil, validTestConfig()); err == nil {
		t.Fatal("expected error")
	}
}

func TestDisconnectQuietly_nil(t *testing.T) { disconnectQuietly(nil) }

func TestConn_concurrentCloseAndAccess(t *testing.T) {
	c, err := mongo.Connect(options.Client().ApplyURI(BuildURI(validTestConfig())))
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{client: c, DB: c.Database("x")}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = conn.Database()
			_ = conn.RawClient()
			_ = conn.Close(context.Background())
		}()
	}
	wg.Wait()
	if conn.RawClient() != nil || conn.Database() != nil {
		t.Fatal("expected closed conn")
	}
	if err := conn.Ping(context.Background()); !errors.Is(err, ErrInvalidConfig) {
		t.Fatal(err)
	}
}

func TestConfig_Validate_writeConcernCombos(t *testing.T) {
	cfg := validTestConfig()
	cfg.WriteConcern = &WriteConcernConfig{W: "0", J: true}
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidConfig) {
		t.Fatal(err)
	}
	cfg.WriteConcern = &WriteConcernConfig{W: "majority", J: true, WTimeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	cfg.WriteConcern = &WriteConcernConfig{W: "x"}
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidConfig) {
		t.Fatal(err)
	}
}

func TestRedactedURI_hidesPassword(t *testing.T) {
	cfg := validTestConfig()
	cfg.Username, cfg.Password = "u", "s3cr3t"
	if got := cfg.RedactedURI(); got == "" || contains(got, "s3cr3t") || !contains(got, "xxxxx") {
		t.Fatal(got)
	}
	if got := RedactedURI(Config{Database: "d"}); contains(got, "xxxxx") {
		t.Fatal(got)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

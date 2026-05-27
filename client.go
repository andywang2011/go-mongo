package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// connectMongo is replaced in tests; production uses defaultConnectMongo.
var connectMongo = defaultConnectMongo

// pingAfterConnect is replaced in tests to avoid a live server.
var pingAfterConnect = func(c *mongo.Client, ctx context.Context) error {
	return c.Ping(ctx, nil)
}

func defaultConnectMongo(ctx context.Context, cfg Config) (*mongo.Client, error) {
	opts, err := applyClientOptions(cfg)
	if err != nil {
		return nil, err
	}
	return mongo.Connect(ctx, opts)
}

// New connects, pings, and returns a Client. Call Close on shutdown.
func New(cfg Config) (Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	c, err := connectMongo(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := pingAfterConnect(c, ctx); err != nil {
		_ = c.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	return &Conn{
		client: c,
		cfg:    cfg,
		DB:     c.Database(cfg.Database),
	}, nil
}

// Ping checks server reachability.
func (m *Conn) Ping(ctx context.Context) error {
	if m == nil || m.client == nil {
		return ErrInvalidConfig
	}
	return m.client.Ping(ctx, nil)
}

// RawClient returns the underlying driver client for advanced use.
func (m *Conn) RawClient() *mongo.Client {
	if m == nil {
		return nil
	}
	return m.client
}

// Client returns the underlying driver client. Deprecated: use RawClient.
func (m *Conn) Client() *mongo.Client {
	return m.RawClient()
}

// Close disconnects from MongoDB. Safe on nil receiver and idempotent after first close.
func (m *Conn) Close(ctx context.Context) error {
	if m == nil || m.client == nil {
		return nil
	}
	err := m.client.Disconnect(ctx)
	m.client = nil
	m.DB = nil
	return err
}

package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// disconnectTimeout bounds cleanup Disconnect calls after a failed connect or ping.
var disconnectTimeout = 5 * time.Second

// connectMongo is replaced in tests; production uses defaultConnectMongo.
var connectMongo = defaultConnectMongo

// pingAfterConnect is replaced in tests to avoid a live server.
var pingAfterConnect = func(c *mongo.Client, ctx context.Context) error {
	return c.Ping(ctx, nil)
}

// defaultConnectMongo builds driver options and connects. mongo-driver v2's
// Connect takes no context (it does not block on the network); ctx bounds the
// subsequent ping and is honoured here only for early cancellation.
func defaultConnectMongo(ctx context.Context, cfg Config) (*mongo.Client, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	opts, err := applyClientOptions(cfg)
	if err != nil {
		return nil, err
	}
	return mongo.Connect(opts)
}

// New connects, pings, and returns a Client. Call Close on shutdown.
// ctx bounds the connect and ping (together with Config.ConnectTimeout).
func New(ctx context.Context, cfg Config) (Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cctx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	c, err := connectMongo(cctx, cfg)
	if err != nil {
		disconnectQuietly(c)
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := pingAfterConnect(c, cctx); err != nil {
		disconnectQuietly(c)
		return nil, fmt.Errorf("mongo ping: %w", err)
	}
	return &Conn{
		client: c,
		cfg:    cfg,
		DB:     c.Database(cfg.Database),
	}, nil
}

// disconnectQuietly releases a partially initialised client using a fresh
// bounded context, so cleanup still runs when the connect ctx has expired.
func disconnectQuietly(c *mongo.Client) {
	if c == nil {
		return
	}
	dctx, dcancel := context.WithTimeout(context.Background(), disconnectTimeout)
	defer dcancel()
	_ = c.Disconnect(dctx)
}

// Ping checks server reachability.
func (m *Conn) Ping(ctx context.Context) error {
	if m == nil {
		return ErrInvalidConfig
	}
	m.mu.RLock()
	c := m.client
	m.mu.RUnlock()
	if c == nil {
		return ErrInvalidConfig
	}
	return c.Ping(ctx, nil)
}

// RawClient returns the underlying driver client for advanced use.
func (m *Conn) RawClient() *mongo.Client {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.client
}

// Close disconnects from MongoDB. Safe on nil receiver and idempotent after first close.
func (m *Conn) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	c := m.client
	m.client = nil
	m.DB = nil
	m.mu.Unlock()
	if c == nil {
		return nil
	}
	return c.Disconnect(ctx)
}

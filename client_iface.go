package mongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Client is the MongoDB connection abstraction (mock-friendly).
type Client interface {
	Database() *mongo.Database
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

var _ Client = (*Conn)(nil)

// Conn holds a connected mongo.Client and the default database.
type Conn struct {
	mu     sync.RWMutex
	client *mongo.Client
	cfg    Config
	// DB is the default database; prefer Database() for new code.
	DB *mongo.Database
}

// Database returns the default database handle.
func (m *Conn) Database() *mongo.Database {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.DB
}

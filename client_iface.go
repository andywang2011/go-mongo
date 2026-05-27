package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
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
	client *mongo.Client
	cfg    Config
	DB     *mongo.Database // prefer Database(); kept for backward compatibility
}

// Database returns the default database handle.
func (m *Conn) Database() *mongo.Database {
	if m == nil {
		return nil
	}
	return m.DB
}

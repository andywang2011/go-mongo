package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

type fakeClient struct {
	db *mongo.Database
}

func (f *fakeClient) Database() *mongo.Database { return f.db }
func (f *fakeClient) Ping(context.Context) error { return nil }
func (f *fakeClient) Close(context.Context) error { return nil }

var _ Client = (*fakeClient)(nil)

func TestFakeClient_implementsClient(t *testing.T) {
	var _ Client = (*fakeClient)(nil)
}

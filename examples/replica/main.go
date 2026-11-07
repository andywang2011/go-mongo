// Replica set hosts with tuned connection pool (requires a real cluster).
package main

import (
	"context"
	"log"
	"time"

	mongo "github.com/andywang2011/go-mongo"
)

func main() {
	cfg := mongo.Config{
		Addresses:       []string{"mongo1:27017", "mongo2:27017", "mongo3:27017"},
		ReplicaSet:      "rs0",
		Database:        "app",
		ReadPreference:  mongo.ReadPrimaryPreferred,
		MaxPoolSize:     100,
		MinPoolSize:     25,
		MaxConnIdleTime: 60 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Close(ctx)
	}()

	if err := client.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("connected to replica set", cfg.ReplicaSet)
}

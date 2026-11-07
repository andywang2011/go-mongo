// Basic insert/find with graceful Close.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mongo "github.com/andywang2011/go-mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	cfg := mongo.Config{
		Address:  "127.0.0.1:27017",
		Database: "example",
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	client, err := mongo.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Close(ctx)
	}()

	coll := client.Database().Collection("demo")
	ctx := context.Background()

	_, err = coll.InsertOne(ctx, bson.M{"hello": "world"})
	if err != nil {
		log.Fatal(err)
	}

	var doc bson.M
	if err := coll.FindOne(ctx, bson.M{"hello": "world"}).Decode(&doc); err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc)
}

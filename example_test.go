package mongo_test

import (
	"context"
	"fmt"

	"github.com/andywang2011/go-mongo"
)

func ExampleNew() {
	cfg := mongo.ConfigDefault()
	cfg.Database = "app"
	cfg.Address = "127.0.0.1:65530"

	client, err := mongo.New(cfg)
	if err != nil {
		fmt.Println("unavailable")
		return
	}
	defer client.Close(context.Background())
	fmt.Println("ok")
	// Output: unavailable
}

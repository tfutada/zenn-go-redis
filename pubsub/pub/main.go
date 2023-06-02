package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"time"
)

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		PoolSize: 1000,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Publish messages
	for i := 0; i < 10; i++ {
		msg := genKey()
		err := rdb.Publish(ctx, "mychannel1", msg).Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("publish:", msg)
		time.Sleep(1 * time.Second)
	}
}

func genKey() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("failed to generate UUID: %v", err)
		panic(err)
	}
	return fmt.Sprintf("msg:%s", u.String())
}

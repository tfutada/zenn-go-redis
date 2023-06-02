package main

import (
	"context"
	"fmt"
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

	pubsub := rdb.Subscribe(ctx, "mychannel1")
	defer pubsub.Close()

	ch := pubsub.Channel(redis.WithChannelSize(0))

	for msg := range ch {
		time.Sleep(30 * time.Second)
		fmt.Printf("receive: %v %v \n", msg.Channel, msg.Payload)
	}
	// unreachable here
}

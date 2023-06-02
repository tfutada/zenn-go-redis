package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

const THREADS = 2

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		PoolSize: 1000,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	raceCondition(ctx, rdb)
}

func raceCondition(ctx context.Context, rdb *redis.Client) {
	var wg sync.WaitGroup
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//
			subscriber(ctx, rdb) // never return
		}()
	}
	wg.Wait()
}

// consume message from channel
func subscriber(ctx context.Context, rdb *redis.Client) {

	pubsub := rdb.Subscribe(ctx, "mychannel1")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch { // infinite loop
		fmt.Println(msg.Channel, msg.Payload)
	}
	// unreachable here
}

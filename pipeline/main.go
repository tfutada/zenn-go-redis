package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func main() {
	println("Started...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pipeLine(rdb)
}

func pipeLine(rdb *redis.Client) {
	var incr *redis.IntCmd

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "pipelined_counter")
		incr = pipe.Incr(ctx, "pipelined_counter")
		pipe.Expire(ctx, "pipelined_counter", time.Hour)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// The value is available only after the pipeline is executed.
	fmt.Println(incr.Val())
}

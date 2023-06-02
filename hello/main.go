package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		PoolSize: 1000,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Set(ctx, "mykey1", "hoge", 0)
	ret, err := rdb.Get(ctx, "mykey1").Result()
	if err != nil {
		println("Error: ", err)
		return
	}

	println("Result: ", ret)
}

package main

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
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

	//
	limiter := redis_rate.NewLimiter(rdb)

	for i := 0; i < 100; i++ {
		res, err := limiter.Allow(ctx, "user-ip:192.168.5.789", redis_rate.PerSecond(10))
		if err != nil {
			panic(err)
		}

		if res.Allowed == 1 {
			println("return 200 ok")
		} else {
			println("return 429 too many requests")
		}

		time.Sleep(50 * time.Millisecond)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

const THREADS = 2
const NumOfGet = 10

var ctx = context.Background()

func main() {
	println("Started...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	raceCondition(rdb)
}

func raceCondition(rdb *redis.Client) {
	var result []int64

	var wg sync.WaitGroup
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//
			ret := pipeLine(rdb)
			result = append(result, ret)
		}()
	}
	wg.Wait()

	fmt.Printf("done. result counter> %v \n", result)
}

func pipeLine(rdb *redis.Client) int64 {
	key := "pipelined_counter"
	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		for i := 0; i < NumOfGet; i++ {
			pipe.Get(ctx, key)
		}
		return nil
	})

	value, err := cmds[len(cmds)-1].(*redis.StringCmd).Result()
	if err != nil {
		fmt.Println("Error converting result to integer: ", err)
		return -1
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Println("Error parsing string to int64: ", err)
		return -1
	}

	return intValue
}

// incr, get(1000), incr, get(1000), incr, get(1000)
// incr, get(500), incr,

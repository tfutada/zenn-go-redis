package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

const THREADS = 2
const NumOfGet = 1000

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
	var result []string

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

func pipeLine(rdb *redis.Client) string {
	key := "pipelined_counter"
	cmds, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		for i := 0; i < NumOfGet; i++ {
			pipe.Get(ctx, key)
		}
		return nil
	})

	value, err := cmds[len(cmds)-1].(*redis.StringCmd).Result()
	if err != nil {
		fmt.Println("error: ", err)
		return "-1"
	}

	return value
}

// incr, get(1000), incr, get(1000), incr, get(1000)
// incr, get(500), incr,

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

const THREADS = 2
const maxRetries = 10

var ctx = context.Background()

func main() {
	println("Started...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	raceCondition(rdb)

	count, err := rdb.Get(ctx, "optimistic_counter").Result()
	if err != nil {
		fmt.Println("Error getting: ", err)
		return
	}
	fmt.Println("Final value: ", count)
}

func raceCondition(rdb *redis.Client) {
	var wg sync.WaitGroup
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := increment(rdb, "optimistic_counter")
			if err != nil {
				fmt.Println("Error incrementing: ", err)
				return
			}
		}()
	}
	wg.Wait()
}

// Increment transactional increments the key using GET and SET commands.
func increment(rdb *redis.Client, key string) error {
	// Transactional function.
	txf := func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int() // Get the current value or zero.
		if err != nil && err != redis.Nil {
			return err
		}

		n++ // Actual operation (local in optimistic lock).

		// Operation is committed only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}

	// Retry if the key has been changed.
	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			return nil // Success.
		}
		if err == redis.TxFailedErr {
			fmt.Println("Optimistic lock lost. Retry.")
			continue
		}
		return err // Return any other error.
	}
	return errors.New("increment reached maximum number of retries")
}

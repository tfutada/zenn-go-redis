package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

const THREADS = 3

type Action func(msg *redis.Message) error

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

// simulate that multiple k8s instances subscribe to the same channel
// each instance can receive copy of the same message. i.e. broadcast
func raceCondition(ctx context.Context, rdb *redis.Client) {
	var wg sync.WaitGroup
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			subscriber(ctx, rdb, i) // never return
		}()
	}
	wg.Wait()
}

// consume message from channel
func subscriber(ctx context.Context, rdb *redis.Client, i int) {

	pubsub := rdb.Subscribe(ctx, "mychannel1")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch { // infinite loop
		//
		err := LockMessage(ctx, rdb, msg, func(msg *redis.Message) error {
			// only one thread(k8s instance) can be here
			fmt.Printf("process msg: thread #%d, %v \n", i, msg.Payload)
			time.Sleep(10 * time.Second)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	// unreachable here
}

func LockMessage(ctx context.Context, rdb *redis.Client, msg *redis.Message, action Action) error {
	rLock := msg.Payload // name of lock key should be unique id of message

	nx, err := rdb.SetNX(ctx, rLock, true, 1*time.Hour).Result()
	if err != nil {
		return err // something wrong with Redis or network
	}
	if !nx {
		return nil // other thread have got the lock.
		//return errors.New("cannot to get a lock as someone else has it")
	} else {
		fmt.Printf("lock on %s \n", rLock) // get the lock successfully
	}

	// now the thread are able to enter CRITICAL SECTION
	err = action(msg) // run a task with the message
	if err != nil {
		return err
	}

	return nil
}

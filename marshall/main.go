package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	Name    string
	Version int64
	Code    int64
}

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		PoolSize: 1000,
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	s := &Message{Name: "Android", Version: 13, Code: 1}

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := string(b)

	rdb.Set(ctx, "mykey1", m, 0)
	ret, err := rdb.Get(ctx, "mykey1").Result()
	if err != nil {
		println("Error: ", err)
		return
	}
	println("Result: ", ret)
	var mm Message
	u := json.Unmarshal([]byte(ret), &mm)
	if u != nil {
		fmt.Println("Error unmarshalling:", u)
		return
	}
	fmt.Printf("%#v", mm)
}

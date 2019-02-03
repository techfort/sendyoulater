package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := InitAPI(r)

	if err != nil {
		panic(err)
	}

	fmt.Println("this is the api...")
}

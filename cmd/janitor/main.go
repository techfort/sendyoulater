package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/techfort/sendyoulater"
)

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	j := sendyoulater.NewJanitor(r)
	count, err := j.CleanUp(10, "test*")
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
	}
	fmt.Printf("Count: %v", count)
}

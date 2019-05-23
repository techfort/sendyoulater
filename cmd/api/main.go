package main

import (
	"fmt"

	"github.com/techfort/sendyoulater"
)

func main() {
	v := sendyoulater.Env()
	_, err := InitAPI(v)
	if err != nil {
		panic(err)
	}

	fmt.Println("this is the api...")
}

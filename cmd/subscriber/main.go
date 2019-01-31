package main

import (
	"fmt"
	"strings"

	"github.com/techfort/forward"
)

func main() {
	ps := forward.NewPubSub(forward.PubSubConfig{
		Addr:       "localhost:6379",
		KeyspaceID: 0,
		Pattern:    "*",
	})
	events, errs := ps.Channel()
	for e := range events {
		fmt.Println(fmt.Sprintf("%+v", e))
		if e.Type == "expired" {
			actionKey := strings.Replace(e.Key, "shdw:", "", 1)
			fmt.Println(actionKey)
		}
	}
	<-errs
}

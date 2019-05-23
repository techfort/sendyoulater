package main

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/techfort/sendyoulater"

	"github.com/techfort/forward"
)

func main() {
	v := sendyoulater.Env()
	ps := forward.NewPubSub(forward.PubSubConfig{
		Addr:       v.GetString("redis_url"),
		KeyspaceID: 0,
		Pattern:    "*",
	})

	r := redis.NewClient(&redis.Options{
		Addr: v.GetString("redis_url"),
	})

	fmt.Println("CLIENT_ID", v.GetString("oauth_clientid"))

	smsHandler, emailHandler := sendyoulater.NewSMSHandler(r, v), sendyoulater.NewEmailHandler(r, v)

	events, errs := ps.Channel()
	for e := range events {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("failed to process event")
				}
				fmt.Println(err.Error())
			}
		}()

		fmt.Println(fmt.Sprintf("%+v", e))
		if e.Type == "expired" {

			action, key := sendyoulater.ParseShadowKey(e.Key)
			fmt.Println(fmt.Sprintf("Processing action: %v", action))
			if action == "email" {
				err := emailHandler.Handle(key)
				if err != nil {
					panic(errors.Wrap(err, "failed to process email"))
				}
			}
			if action == "sms" {
				err := smsHandler.Handle(key)
				if err != nil {
					panic(errors.Wrap(err, "failed to process sms"))
				}
			}
		}
	}
	<-errs
}

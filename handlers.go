package sendyoulater

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Handler interface
type Handler interface {
	Handle(key string) error
}

// EventHandler is the base type
type EventHandler struct {
	store
}

// SMSHandler handlers sms
type SMSHandler struct {
	EventHandler
}

// EmailHandler handles Emails
type EmailHandler struct {
	EventHandler
}

// NewSMSHandler returns a SMSHandler
func NewSMSHandler(r *redis.Client) Handler {
	return SMSHandler{EventHandler{store{r}}}
}

// NewEmailHandler returns an Email Handler
func NewEmailHandler(r *redis.Client) Handler {
	return EmailHandler{EventHandler{store{r}}}
}

// Handle handles the actual sms event
func (sms SMSHandler) Handle(key string) error {
	fmt.Println(fmt.Sprintf("Handling SMS %+v", key))
	_, actionKey := ParseShadowKey(key)
	smsMap, err := sms.HGetAll(actionKey).Result()
	fmt.Println(fmt.Sprintf("SMS: %+v", smsMap))
	return err
}

// Handle handles the actual email event
func (email EmailHandler) Handle(key string) error {
	fmt.Println(fmt.Sprintf("Handling Email: %+v", key))
	_, actionKey := ParseShadowKey(key)
	emailMap, err := email.HGetAll(actionKey).Result()
	fmt.Println(fmt.Sprintf("Email: %+v", emailMap))
	return err
}

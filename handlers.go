package sendyoulater

import (
	"fmt"

	"github.com/techfort/forward"
)

// Handler interface
type Handler interface {
	Handle() error
}

// EventHandler is the base type
type EventHandler struct {
	Event forward.RedisKV
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
func NewSMSHandler(e forward.RedisKV) Handler {
	return SMSHandler{EventHandler{e}}
}

// NewEmailHandler returns an Email Handler
func NewEmailHandler(e forward.RedisKV) Handler {
	return EmailHandler{EventHandler{e}}
}

// Handle handles the actual sms event
func (sms SMSHandler) Handle() error {
	fmt.Println(fmt.Sprintf("Handling SMS %+v", sms.Event))
	return nil
}

// Handle handles the actual email event
func (email EmailHandler) Handle() error {
	fmt.Println(fmt.Sprintf("Handling Email: %+v", email.Event))
	return nil
}

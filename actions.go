package sendyoulater

import "time"

// Action base type
type Action struct {
	UserID    string
	Timestamp time.Time
	Delay     time.Duration
}

// EmailAction is an Action for emails
type EmailAction struct {
	Action
	To      string
	Subject string
	Body    string
}

// SMSAction is an Action for SMS
type SMSAction struct {
	Action
	To   string
	Body string
}

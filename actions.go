package sendyoulater

import (
	"time"

	"github.com/pkg/errors"
)

// Action base type
type Action struct {
	ID        string
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

// FromMap inflates an EmailAction from a redis result
func (a *EmailAction) FromMap(result map[string]string) (EmailAction, error) {
	delay, err := time.ParseDuration(result["Delay"])
	if err != nil {
		return *a, errors.Wrap(err, "cannot convert Delay value to integer")
	}
	timestamp, err := time.Parse(TimeFormat, result["Timestamp"])
	if err != nil {
		return *a, errors.Wrap(err, "cannot convert timestamp")
	}
	a.Action = Action{ID: result["ID"], UserID: result["UserID"], Delay: time.Duration(delay) * time.Second, Timestamp: timestamp}
	a.To = result["To"]
	a.Body = result["Body"]
	a.Subject = result["Subject"]
	return *a, err
}

// SMSAction is an Action for SMS
type SMSAction struct {
	Action
	To   string
	Body string
}

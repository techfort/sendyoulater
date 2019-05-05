package sendyoulater

import (
	"fmt"
	"strings"
	"time"
)

const (
	// SMS literal
	SMS = "sms"
	// Email literal
	Email = "email"
	// TimeFormat for datetime formatting
	TimeFormat = "2006-01-02"
	// Shadow unformatted key
	Shadow = `shd:`
	// UserKEY key
	UserKEY = `user:%v`
	// PlanKEY key
	PlanKEY = `plan:%v`
	// EmailActionKEY unformatted key
	EmailActionKEY = `u:%v:email:%v`
	// ShadowEmailActionKEY unformatted key
	ShadowEmailActionKEY = Shadow + EmailActionKEY
	// SMSActionKEY key
	SMSActionKEY = `u:%v:sms:%v`
	// ShadowSMSActionKEY key
	ShadowSMSActionKEY = Shadow + SMSActionKEY
	// EmailActionsForUserKEY key
	EmailActionsForUserKEY = `eafu:%v`
	// SMSActionsForUserKEY key
	SMSActionsForUserKEY = `safu:%v`
)

// KeyPlan returns the formatted plan key
func KeyPlan(name string) string {
	return fmt.Sprintf(PlanKEY, name)
}

// KeyUser returns the formatted user key
func KeyUser(userID string) string {
	return fmt.Sprintf(UserKEY, userID)
}

// DateID returns an id composed of the date and the counter since a counter is monthly valid
func DateID(counter int64) string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%v-%v-%v-%v", year, month, day, counter)
}

// KeysEmailAction returns key and shadow key for email actions
func KeysEmailAction(userID string, counter int64) (string, string) {
	dateID := DateID(counter)
	return fmt.Sprintf(EmailActionKEY, userID, dateID), fmt.Sprintf(ShadowEmailActionKEY, userID, dateID)
}

// KeyEmailActionsForUser returns the formatted key
func KeyEmailActionsForUser(userID string) string {
	return fmt.Sprintf(EmailActionsForUserKEY, userID)
}

// KeySMSActionsForUser returns the formatted key
func KeySMSActionsForUser(userID string) string {
	return fmt.Sprintf(SMSActionsForUserKEY, userID)
}

// KeySMSAction returns shadow and key for sms action
func KeySMSAction(userID string, counter int64) (string, string) {
	dateID := DateID(counter)
	return fmt.Sprintf(SMSActionKEY, userID, dateID), fmt.Sprintf(ShadowSMSActionKEY, userID, dateID)
}

// ParseShadowKey returns the type of action given a certain shadow key
func ParseShadowKey(shadowKey string) (string, string) {
	chunks := strings.Split(shadowKey, ":")
	return chunks[3], strings.Replace(shadowKey, Shadow, "", 1)
}

package sendyoulater

import (
	"fmt"
	"strings"
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
)

// KeyPlan returns the formatted plan key
func KeyPlan(name string) string {
	return fmt.Sprintf(PlanKEY, name)
}

// KeyUser returns the formatted user key
func KeyUser(userID string) string {
	return fmt.Sprintf(UserKEY, userID)
}

// KeysEmailAction returns key and shadow key for email actions
func KeysEmailAction(userID string, counter int64) (string, string) {
	return fmt.Sprintf(EmailActionKEY, userID, counter), fmt.Sprintf(ShadowEmailActionKEY, userID, counter)
}

// KeySMSAction returns shadow and key for sms action
func KeySMSAction(userID string, counter int64) (string, string) {
	return fmt.Sprintf(SMSActionKEY, userID, counter), fmt.Sprintf(ShadowSMSActionKEY, userID, counter)
}

// ParseShadowKey returns the type of action given a certain shadow key
func ParseShadowKey(shadowKey string) (string, string) {
	chunks := strings.Split(shadowKey, ":")
	return chunks[3], strings.Replace(shadowKey, Shadow, "", 1)
}

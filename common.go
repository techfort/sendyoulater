package sendyoulater

import "fmt"

const (
	// TimeFormat for datetime formatting
	TimeFormat = "2006-01-02 15:04:05 +0000 UTC"
	// Shadow unformatted key
	Shadow = `shd:`
	// UserKEY key
	UserKEY = `u:%v`
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

// KeysEmailAction returns key and shadow key for email actions
func KeysEmailAction(userID string, counter int64) (string, string) {
	return fmt.Sprintf(EmailActionKEY, userID, counter), fmt.Sprintf(ShadowEmailActionKEY, userID, counter)
}

// KeySMSAction returns shadow and key for sms action
func KeySMSAction(userID string, counter int64) (string, string) {
	return fmt.Sprintf(SMSActionKEY, userID, counter), fmt.Sprintf(ShadowSMSActionKEY, userID, counter)
}

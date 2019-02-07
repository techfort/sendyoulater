package sendyoulater

import (
	"time"

	"github.com/pkg/errors"
)

type smsRepo struct {
	store
}

// SMSActionRepo interface to save sms actions
type SMSActionRepo interface {
	SaveSMSAction(user User, to, text string, ex time.Duration) (SMSAction, error)
}

type smsUseCase struct {
	store
	UserRepository UserRepo
	PlanRepository PlanRepo
	SMSRepository  SMSActionRepo
}

// SMSUseCase for saving SMS actions based on user and plans
type SMSUseCase interface {
	SaveSMS(userID, to, text string, ex time.Duration) (SMSAction, error)
}

// NewSMSActionRepo returns a SMSActionRepo
func (s store) NewSMSActionRepo() SMSActionRepo {
	return smsRepo{s}
}

// NewSMSUseCase returns a SMSUseCase
func (s store) NewSMSUseCase(ur UserRepo, pr PlanRepo, sar SMSActionRepo) SMSUseCase {
	return smsUseCase{s, ur, pr, sar}
}

// SaveSMS saves and SMS action
func (suc smsUseCase) SaveSMS(userID, to, text string, ex time.Duration) (SMSAction, error) {
	var sa SMSAction
	user, err := suc.UserRepository.ByID(userID)
	if err != nil {
		return sa, err
	}
	plan, err := suc.PlanRepository.ByName(user.Plan)
	if err != nil {
		return sa, err
	}
	if user.SMSCounter <= plan.MaxSMS {
		sa, err := suc.SMSRepository.SaveSMSAction(user, to, text, ex)
		if err != nil {
			return sa, errors.Wrap(err, "failed to save sms action")
		}
		user.SMSCounter++
		_, err = suc.UserRepository.Update(user)
	} else {
		return sa, errors.New("Maximum number of SMS reached for user")
	}
	return sa, nil
}

// SaveSMSAction is a repo method that saves the action to redis
func (sr smsRepo) SaveSMSAction(user User, to, text string, ex time.Duration) (SMSAction, error) {
	actionKey, shadowKey := KeySMSAction(user.UserID, user.SMSCounter+1)
	pipe := sr.TxPipeline()
	pipe.Set(shadowKey, ex, ex)
	sa := SMSAction{Action: Action{Delay: ex, Timestamp: time.Now(), UserID: user.UserID}, Body: text, To: to}
	pipe.HMSet(actionKey, map[string]interface{}{
		"Action": SMS,
		"To":     to,
		"Body":   text,
	})
	_, err := pipe.Exec()
	return sa, err
}

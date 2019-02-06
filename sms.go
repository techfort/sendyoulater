package sendyoulater

import "time"

type smsRepo struct {
	store
}

type SMSActionRepo interface {
	SaveSMSAction(user User, plan Plan, to, text string, ex time.Duration) (SMSAction, error)
}

type smsUseCase struct {
	store
	UserRepository UserRepo
	PlanRepository PlanRepo
	SMSRepository  SMSActionRepo
}

type SMSUseCase interface {
	SaveSMS(userID, to, text string, ex time.Duration) (SMSAction, error)
}

func (s store) NewSMSActionRepo() SMSActionRepo {
	return smsRepo{s}
}

func (suc smsUseCase) SaveSMS(userID, to, text string, ex time.Duration) (SMSAction, error) {
	return SMSAction{}, nil
}

func (sr smsRepo) SaveSMSAction(user User, plan Plan, to, text string, ex time.Duration) (SMSAction, error) {
	return SMSAction{}, nil
}

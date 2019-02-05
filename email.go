package sendyoulater

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type emailUseCase struct {
	store
	UserRepository UserRepo
	PlanRepository PlanRepo
}

type emailRepo struct {
	store
}

// EmailActionRepo interface for saving/updating email actions
type EmailActionRepo interface {
	SaveEmailAction(user User, plan Plan, subject, body, to string, ex time.Duration) (EmailAction, error)
}

func (s store) NewEmailRepo() EmailActionRepo {
	return emailRepo{s}
}

// EmailUseCase interface to save EmailActions
type EmailUseCase interface {
	SaveEmailActions(userID, subject, body string, to []string, ex time.Duration) ([]EmailAction, error)
}

func (s store) NewEmailUseCase(u UserRepo, p PlanRepo) EmailUseCase {
	return emailUseCase{s, u, p}
}

func (s store) SaveEmailAction(user User, plan Plan, subject, body, to string, ex time.Duration) (EmailAction, error) {
	shadow, actionKey := KeysEmailAction(user.UserID, user.EmailCounter)
	if _, err := s.Set(shadow, "email", ex).Result(); err != nil {
		return EmailAction{}, errors.Wrap(err, "error setting shadow key, email action not saved")
	}
	_, err := s.HMSet(actionKey, map[string]interface{}{
		"UserID":    user.UserID,
		"Timestamp": time.Now().Format(TimeFormat),
		"To":        to,
		"Subject":   subject,
		"Body":      body,
	}).Result()
	return EmailAction{Action: Action{UserID: user.UserID, Timestamp: time.Now(), Delay: ex}, To: to, Subject: subject, Body: body}, err
}

func (s emailUseCase) SaveEmailActions(userID, subject, body string, to []string, ex time.Duration) ([]EmailAction, error) {
	var (
		err     error
		user    User
		plan    Plan
		actions []EmailAction
	)
	user, err = s.UserRepository.ByID(userID)
	if err != nil {
		return []EmailAction{}, errors.Wrap(err, fmt.Sprintf("failed to retrieve user %v", userID))
	}
	if plan, err = s.PlanRepository.ByName(user.Plan); err != nil {
		return []EmailAction{}, errors.Wrap(err, fmt.Sprintf("failed to retrieve plan %v", user.Plan))
	}
	g, _ := errgroup.WithContext(context.Background())
	if (user.EmailCounter + int64(len(to))) < plan.MaxEmails {
		for _, rec := range to {
			g.Go(func() error {
				rec := rec
				a, err := s.SaveEmailAction(user, plan, subject, body, rec, ex)
				actions = append(actions, a)
				if err != nil {
					return err
				}
				user.EmailCounter++
				return err
			})
		}
	}
	if err := g.Wait(); err != nil {
		panic(errors.Wrap(err, "failed to save email action"))
	}

	return actions, err
}

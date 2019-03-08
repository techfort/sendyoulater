package sendyoulater

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type sessionStore struct {
	*redis.Client
}

// SessionStore intetrface for setting sessions
type SessionStore interface {
	Set(sessionID, userID string) error
	Get(sessionID string) (string, error)
	Touch(sessionID string) error
}

// NewSessionStore returns a SessionStore interface
func NewSessionStore(r *redis.Client) SessionStore {
	return sessionStore{r}
}

// Get returns the content of a sessionID key
func (s sessionStore) Get(sessionID string) (string, error) {
	return s.Client.Get(sessionID).Result()
}

// Set associates a sessionID witha user id, timeout set to an hour
func (s sessionStore) Set(sessionID, userID string) error {
	_, err := s.Client.Set(sessionID, userID, time.Duration(60)*time.Minute).Result()
	if err != nil {
		return errors.Wrap(err, "could not set session")
	}
	return err
}

// Touch refreshes a session
func (s sessionStore) Touch(sessionID string) error {
	ok, err := s.Client.Expire(sessionID, time.Duration(60)*time.Minute).Result()
	if !ok {
		return errors.New("Could not touch session")
	}
	if err != nil {
		return errors.Wrap(err, "Could not touch session")
	}
	return err
}

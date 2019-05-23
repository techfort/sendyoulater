package sendyoulater

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Handler interface
type Handler interface {
	Handle(key string) error
}

// EventHandler is the base type
type EventHandler struct {
	store
	Viper *viper.Viper
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
func NewSMSHandler(r *redis.Client, v *viper.Viper) Handler {
	return SMSHandler{EventHandler{store{r}, v}}
}

// NewEmailHandler returns an Email Handler
func NewEmailHandler(r *redis.Client, v *viper.Viper) Handler {
	return EmailHandler{EventHandler{store{r}, v}}
}

// Handle handles the actual sms event
func (sms SMSHandler) Handle(key string) error {
	fmt.Println(fmt.Sprintf("Handling SMS %+v", key))
	_, actionKey := ParseShadowKey(key)
	smsMap, err := sms.HGetAll(actionKey).Result()
	fmt.Println(fmt.Sprintf("SMS: %+v", smsMap))
	return err
}

// Handle handles the actual email event
func (email EmailHandler) Handle(key string) error {
	fmt.Println(fmt.Sprintf("Handling Email: %+v", key))
	_, actionKey := ParseShadowKey(key)
	emailMap, err := email.HGetAll(actionKey).Result()
	userID, err := UserIDFromKey(actionKey)
	if err != nil {
		return errors.Wrap(err, "failed to parse key")
	}
	store := NewStore(email.Client)
	ur := store.NewUserRepo()
	user, err := ur.ByID(userID)
	if err != nil {
		return errors.Wrap(err, "cannot find user")
	}
	tok := &oauth2.Token{
		AccessToken:  user.Token,
		RefreshToken: user.RefreshToken,
		TokenType:    "Bearer",
	}
	conf := Oauth2Config(email.Viper)
	client := conf.Client(context.Background(), tok)

	rtr, err := RefreshToken(client, conf, tok)
	if err != nil {
		return errors.Wrap(err, "failed to refresh token")
	}
	fmt.Println("TOKEN REFRESH SUCCESSFUL", rtr)
	user.Token = rtr.AccessToken
	if _, err := ur.Update(user); err != nil {
		return errors.Wrap(err, "failed to update user with refreshed token info")
	}
	newToken := &oauth2.Token{
		AccessToken:  rtr.AccessToken,
		RefreshToken: user.RefreshToken,
		TokenType:    "Bearer",
	}
	client = conf.Client(context.Background(), newToken)
	srv, err := gmail.New(client)
	if err != nil {
		return errors.Wrap(err, "failed to create service from updated token client")
	}
	message := Message(user.UserID, emailMap["To"], emailMap["Subject"], emailMap["Body"])
	_, err = srv.Users.Messages.Send(userID, &message).Do()
	if err != nil {
		if gapiErr, ok := err.(*googleapi.Error); ok {
			fmt.Println(fmt.Sprintf("Error: %+v", gapiErr))
			if gapiErr.Code == 401 {
				return errors.Wrap(gapiErr, "failed to send message, unauthorised")
			}
			return errors.Wrap(gapiErr, "failed to send message")
		}
	}
	actionsForUserKey := KeyEmailActionsForUser(userID)
	if _, err := email.SRem(actionsForUserKey, actionKey).Result(); err != nil {
		return errors.Wrap(err, "failed to remove action from list of actions for user")
	}
	return err
}

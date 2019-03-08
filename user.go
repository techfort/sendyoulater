package sendyoulater

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// User represents a user
type User struct {
	UserID       string
	FirstName    string
	LastName     string
	Plan         string
	Company      string
	SMSCounter   int64
	EmailCounter int64
	PeriodStart  time.Time
	Token        string
	RefreshToken string
}

type userRepo struct {
	store
}

// UserRepo wraps the User related methods of store
type UserRepo interface {
	ByID(userID string) (User, error)
	Save(userID, firstname, lastname, plan, company, token, refresh string) (User, error)
	Update(user User) (User, error)
}

func (s store) NewUserRepo() UserRepo {
	return userRepo{s}
}

// User retrieves a user
func (u userRepo) ByID(userID string) (User, error) {
	var (
		err    error
		user   User
		smsc   int64
		emailc int64
		start  time.Time
	)
	ret, err := u.HGetAll(KeyUser(userID)).Result()
	fmt.Println(fmt.Sprintf("ret: %v, %+v", ret["UserID"], ret))
	if err != nil {
		return user, err
	}

	if emailc, err = strconv.ParseInt(ret["EmailCounter"], 10, 64); err != nil {
		return user, errors.Wrap(err, fmt.Sprintf("Failed to convert EmailCounter"))
	}
	if smsc, err = strconv.ParseInt(ret["SMSCounter"], 10, 64); err != nil {
		return user, errors.Wrap(err, fmt.Sprintf("Failed to convert SMSCounter"))
	}

	fmt.Println(ret["PeriodStart"])
	if start, err = time.Parse(TimeFormat, ret["PeriodStart"]); err != nil {
		return user, errors.Wrap(err, fmt.Sprintf("Failed to convert PeriodStart"))
	}
	user.UserID = ret["UserID"]
	user.FirstName = ret["FirstName"]
	user.LastName = ret["LastName"]
	user.Company = ret["Company"]
	user.Plan = ret["Plan"]
	user.PeriodStart = start
	user.EmailCounter = emailc
	user.SMSCounter = smsc
	user.Token = ret["Token"]
	user.RefreshToken = ret["RefreshToken"]
	return user, err
}

func (u userRepo) Save(userID, firstname, lastname, plan, company, token, refresh string) (User, error) {
	userMap := map[string]interface{}{
		"UserID":       userID,
		"FirstName":    firstname,
		"LastName":     lastname,
		"Plan":         plan,
		"Company":      company,
		"EmailCounter": 0,
		"SMSCounter":   0,
		"PeriodStart":  time.Now().Format(TimeFormat),
		"Token":        token,
		"RefreshToken": refresh,
	}
	if _, err := u.HMSet(KeyUser(userID), userMap).Result(); err != nil {
		return User{}, errors.Wrap(err, fmt.Sprintf("failed to save user: %+v", userMap))
	}
	return u.ByID(userID)
}

func (u userRepo) Update(user User) (User, error) {
	if _, err := u.HMSet(KeyUser(user.UserID), map[string]interface{}{
		"FirstName":    user.FirstName,
		"LastName":     user.LastName,
		"Company":      user.Company,
		"EmailCounter": user.EmailCounter,
		"SMSCounter":   user.SMSCounter,
		"Token":        user.Token,
		"RefreshToken": user.RefreshToken,
	}).Result(); err != nil {
		return User{}, errors.Wrap(err, fmt.Sprintf("failed to update user: %+v", user))
	}
	return u.ByID(user.UserID)
}

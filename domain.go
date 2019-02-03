package sendyoulater

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Plan is a representation of usage plan
type Plan struct {
	Name      string
	MaxEmails int64
	MaxSMS    int64
}

var (
	Plans = map[string]Plan{
		"basic":      Plan{"basic", 100, 100},
		"enterprise": Plan{"enterprise", 100, 100},
	}
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
}

type userRepo struct {
	store
}

// UserRepo wraps the User related methods of store
type UserRepo interface {
	UserByID(userID string) (User, error)
}

// UserByID retrieves a user from cache
func (ur userRepo) UserByID(userID string) (User, error) {
	return ur.User(userID)
}

// PlanRepo interface that wwraps plan related functions of Store
type PlanRepo interface {
	PlanByName(plan string)
}

type store struct {
	*redis.Client
}

// PlanByName retrieves a plan by its name
func (s store) PlanByName(name string) (Plan, error) {
	var (
		ret       map[string]string
		plan      Plan
		err       error
		maxemails int64
		maxsms    int64
	)
	ret, err = s.HGetAll(name).Result()
	if maxemails, err := strconv.ParseInt(ret["MaxEmails"], 10, 64); err != nil {
		return plan, errors.Wrap(err, "failed to parse maxemails value")
	}
	if maxsms, err := strconv.ParseInt(ret["MaxSMS"], 10, 64); err != nil {
		return plan, errors.Wrap(err, "failed to parse MaxSMS value")
	}
	plan.Name = ret["Name"]
	plan.MaxEmails = maxemails
	plan.MaxSMS = maxsms
	return plan, err
}

// User retrieves a user
func (s store) User(userID string) (User, error) {
	var (
		err    error
		user   User
		ret    map[string]string
		smsc   int64
		emailc int64
		start  time.Time
	)
	ret, err = s.HGetAll(userID).Result()
	if err != nil {
		return user, err
	}
	user.UserID = ret["UserID"]
	user.FirstName = ret["FirstName"]
	user.LastName = ret["LastName"]
	user.Company = ret["Company"]
	user.Plan = ret["Plan"]
	if smsc, err := strconv.ParseInt(ret["SMSCounter"], 10, 64); err != nil {
		return user, err
	}
	if emailc, err := strconv.ParseInt(ret["EmailCounter"], 10, 64); err != nil {
		return user, err
	}
	if start, err := time.Parse(ret["PeriodStart"], TimeFormat); err != nil {
		return user, err
	}
	user.PeriodStart = start
	return user, err
}

func (s store) SetEmailAction(userID, subject, body string, to []string, ex time.Duration) error {
	var (
		err  error
		user User
	)
	user, err = s.User(userID)

}

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/techfort/sendyoulater"
)

// Context is the dis Context for API requests
type Context struct {
	echo.Context
	Redis *redis.Client
}

const (
	// ErrMessage std error
	ErrMessage = `{"message":"unable to process request"}`
)

// BadRequest returns a byte array for a json blob error response
func BadRequest() []byte {
	return []byte(ErrMessage)
}

// Store returns the repo factory
func (c Context) Store() sendyoulater.Store {
	return sendyoulater.NewStore(c.Redis)
}

// Err returns a standard error response
func (c Context) Err() error {
	return c.JSONBlob(http.StatusInternalServerError, BadRequest())
}

// RoutesGET returns get routes
func RoutesGET() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"/user/:id":   UserByID,
		"/plan/:name": PlanByName,
	}
}

// RoutesPOST returns POST routes
func RoutesPOST() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"/action/email/save":  SaveEmailAction,
		"/remove":             Remove,
		"/user/save":          SaveUser,
		"user/:userId/update": UpdateUser,
	}
}

// InitAPI starts the API
func InitAPI(r *redis.Client) (*echo.Echo, error) {
	e := echo.New()
	e = Config(e, r)
	err := e.Start(":1666")
	return e, errors.Wrap(err, "failed to start API")
}

// Config provides middleware and wraps the context of each request
func Config(e *echo.Echo, r *redis.Client) *echo.Echo {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{c, r}
			return h(cc)
		}
	})
	for route, handler := range RoutesGET() {
		e.GET(route, handler)
	}
	for route, handler := range RoutesPOST() {
		e.POST(route, handler)
	}
	return e
}

// SaveEmailAction sets the timer for an action
func SaveEmailAction(c echo.Context) error {
	cc := c.(*Context)
	store := cc.Store()
	ur, pr, er := store.NewUserRepo(), store.NewPlanRepo(), store.NewEmailActionRepo()
	var m echo.Map
	err := cc.Bind(&m)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err.Error()))
		return cc.Err()
	}
	euc := store.NewEmailUseCase(ur, pr, er)
	userID := m["userId"].(string)
	subject := m["subject"].(string)
	body := m["body"].(string)
	toStr := m["to"].(string)
	to := strings.Split(toStr, ",")
	ex, err := time.ParseDuration(m["ex"].(string))
	fmt.Println("Saving email aciton...")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err.Error()))
		return cc.Err()
	}
	ea, err := euc.SaveEmailActions(userID, subject, body, to, ex)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err.Error()))
		return cc.Err()
	}
	return cc.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"status": "ok", "message": "%v actions saved" }`, len(ea))))
}

// UpdateUser updates user information
func UpdateUser(c echo.Context) error {
	cc := c.(*Context)
	var user sendyoulater.User
	if err := cc.Bind(user); err != nil {
		return cc.JSONBlob(http.StatusInternalServerError, []byte(`{"message": "failed to retrieve data"}`))
	}
	ur := cc.Store().NewUserRepo()
	_, err := ur.Update(user)
	if err != nil {
		return cc.Err()
	}
	return cc.JSONBlob(http.StatusOK, []byte(`{"status":"ok", "message": "user updated correctly"`))
}

// SaveUser is the handler for saving user info
func SaveUser(c echo.Context) error {
	var m echo.Map
	cc := c.(*Context)
	if err := cc.Bind(&m); err != nil {
		return err
	}
	ur := cc.Store().NewUserRepo()
	user, err := ur.Save(m["userId"].(string), m["firstname"].(string), m["lastname"].(string), m["plan"].(string), m["company"].(string))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(fmt.Sprintf("User: %+v", user))
	return cc.JSONBlob(http.StatusOK, []byte(`{ "message": "user saved correctly"}`))
}

// Remove removes an action timer
func Remove(c echo.Context) error {
	return nil
}

// UserByID returns a user by id
func UserByID(c echo.Context) error {
	cc := c.(*Context)
	id := cc.Param("id")
	ur := cc.Store().NewUserRepo()
	user, err := ur.ByID(id)
	if err != nil {
		fmt.Println(err.Error())
		return cc.JSONBlob(http.StatusInternalServerError, []byte(`{"message": "could not find user"}`))
	}
	return cc.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"message":"ok", "user":"%+v"}`, user.UserID)))
}

// PlanByName retrieves a plan by name
func PlanByName(c echo.Context) error {
	name := c.Param("name")
	cc := c.(*Context)
	pr := cc.Store().NewPlanRepo()
	p, err := pr.ByName(name)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERR: %v", err.Error()))
		return cc.Err()
	}
	return cc.JSONBlob(http.StatusOK, []byte(fmt.Sprintf(`{"message": "ok", "plan": "%+v"`, p)))
}

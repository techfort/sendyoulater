package main

import (
	"fmt"
	"net/http"

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

// Store returns the repo factory
func (c Context) Store() sendyoulater.Store {
	return sendyoulater.NewStore(c.Redis)
}

// RoutesGET returns get routes
func RoutesGET() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"/user/:id": UserByID,
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
	action := new(sendyoulater.EmailAction)
	if err := cc.Bind(action); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%+v", action))
	return nil
}

// UpdateUser updates user information
func UpdateUser(c echo.Context) error {
	return nil
}

// SaveUser is the handler for saving user info
func SaveUser(c echo.Context) error {
	var m echo.Map
	cc := c.(*Context)
	if err := cc.Bind(&m); err != nil {
		return err
	}
	// TODO: continue here...
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

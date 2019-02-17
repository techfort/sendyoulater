package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/urlshortener/v1"

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
		"login":       Login,
		"auth":        Auth,
		"token":       Token,
	}
}

// RoutesPOST returns POST routes
func RoutesPOST() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"/action/email/save":  SaveEmailAction,
		"/remove":             Remove,
		"/user/save":          SaveUser,
		"user/:userId/update": UpdateUser,
		"init":                initData,
	}
}

// InitAPI starts the API
func InitAPI(r *redis.Client) (*echo.Echo, error) {
	e := echo.New()
	e = Config(e, r)
	err := e.Start(":1323")
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
	e.File("/", "/static/syl-ui/public/index.html")
	e.Use(middleware.Static("/static/syl-ui/public"))

	for route, handler := range RoutesGET() {
		e.GET(route, handler)
	}
	for route, handler := range RoutesPOST() {
		e.POST(route, handler)
	}
	return e
}

func Ui(c echo.Context) error {
	return nil
}

func Login(c echo.Context) error {
	url := fmt.Sprintf("%v", conf.AuthCodeURL("state", oauth2.AccessTypeOffline))
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

var (
	gmailClient *http.Client
	conf        = &oauth2.Config{
		ClientID:     "541640626027-l7s3mcv05cbdhqsq0vf54tcvpprb6s63.apps.googleusercontent.com",
		ClientSecret: "5jbcSzmUBPjFKww6BsoEKpC8",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			urlshortener.UrlshortenerScope,
			gmail.GmailSendScope,
			"https://www.googleapis.com/auth/plus.me",
		},
		RedirectURL: "http://localhost:1323/auth",
	}
)

func Auth(c echo.Context) error {

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	ctx := context.Background()
	fmt.Println(fmt.Sprintf("%+v", c.QueryParams()))
	var p string
	p = c.QueryParam("code")
	var code string
	code = p
	_, err := fmt.Scan(&code)
	if err != nil {
		fmt.Println("code scan failed")
	} else {
		fmt.Printf("%+v", code)
	}

	tok, err := conf.Exchange(ctx, p)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	client := conf.Client(ctx, tok)

	request, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	request.Header.Set("Bearer",
		"ya29.GluzBhHLhGGwcTUfmLvU0VyzfIdgLseObhQLrxeWS0wHmgHRsfJ9uMHALPsslE6JoY4vt13AUVXvH6IJCs7jqxH-7hVBTTFvF-y8mg4TThZ58bV-5BF1qEkrfExj",
	)
	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	b, err := ioutil.ReadAll(res.Body)

	fmt.Printf("RES: %+v", string(b))
	srv, err := gmail.New(client)
	f, err := srv.Users.GetProfile("sendyoulater@gmail.com").Do()
	response := map[string]interface{}{
		"token":   tok,
		"profile": f,
	}
	return c.JSON(http.StatusOK, response)
}

func Token(c echo.Context) error {
	tok := c.QueryParam("token")
	return c.JSON(http.StatusOK, tok)
}

func Service() error {
	return nil
}

func initData(c echo.Context) error {
	cc := c.(*Context)
	pr, ur := cc.Store().NewPlanRepo(), cc.Store().NewUserRepo()
	pr.SavePlan("basic", 100, 100)
	ur.Save("joe", "joe", "minichino", "basic", "sendyoulater")
	return cc.JSONBlob(http.StatusOK, []byte(`{"message":"ok"}`))
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

func SaveSMSAction(c echo.Context) error {
	// cc := c.(*Context)
	// store := cc.Store()
	// ur, pr, sr := store.NewUserRepo(), store.NewPlanRepo(), store.NewSMSRepo()
	// suc := store.
	return nil
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

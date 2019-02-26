package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/techfort/sendyoulater"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/urlshortener/v1"
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
		"token":       Token,
		"auth":        Auth,
		"check":       CheckSession,
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
	session.NewCookieStore()

	store, err := session.NewRedisStore(20, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to redis and create session store"))
	}
	e.Use(session.Sessions("GSESSION", store))
	e.File("/", "../../static/signin/index.html")
	e.Static("/ui", "../../static/syl-ui/dist")
	e.Static("/css", "../../static/syl-ui/dist/css")
	e.Static("/js", "../../static/syl-ui/dist/js")

	for route, handler := range RoutesGET() {
		e.GET(route, handler)
	}
	for route, handler := range RoutesPOST() {
		e.POST(route, handler)
	}
	return e
}

// Login handler
func Login(c echo.Context) error {
	url := fmt.Sprintf("%v", conf.AuthCodeURL("state", oauth2.AccessTypeOffline))
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

var (
	gmailClient *http.Client
	conf        = &oauth2.Config{
		ClientID:     "541640626027-75fep8r1ptdd377l73qhb2f03pckc0po.apps.googleusercontent.com",
		ClientSecret: "7Qbre2sDMxnPHZwa_6DASz4j",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			urlshortener.UrlshortenerScope,
			gmail.GmailSendScope,
			"openid",
			"profile",
			"email",
		},
		RedirectURL: "http://localhost:1323/token",
	}
)

// Token is the callback from google
func Token(c echo.Context) error {
	ctx := context.Background()
	code := c.QueryParam("code")
	fmt.Printf("CODE: %v", code)
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	fmt.Println("Token", fmt.Sprintf("%+v", tok))
	client := conf.Client(oauth2.NoContext, tok)
	url := "https://www.googleapis.com/oauth2/v3/userinfo"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	req.Header.Set("Content-Type", "application/json")
	emailres, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"token": tok,
			"error": err,
		})
	}
	body, err := ioutil.ReadAll(emailres.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	type UserInfo struct {
		Sub       string `json:"sub"`
		Name      string `json:"name"`
		GivenName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Email     string `json:"email"`
	}
	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.SetCookie(&http.Cookie{
		Name:  "SessionID",
		Value: userInfo.Email,
	})
	fmt.Println(tok, userInfo)

	html := `<!DOCTYPE html>
		<html>
		<body>
		<script>
			window.opener.postMessage({ loginSuccessful: true, email: "%v" }, "http://localhost:1323");
			window.close();
		</script>
		</body>
		</html>`
	return c.HTML(http.StatusOK, fmt.Sprintf(html, userInfo.Email))
}

// Auth is unused for now
func Auth(c echo.Context) error {
	var m echo.Map
	err := c.Bind(&m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, m)
}

func initData(c echo.Context) error {
	cc := c.(*Context)
	pr, ur := cc.Store().NewPlanRepo(), cc.Store().NewUserRepo()
	pr.SavePlan("basic", 100, 100)
	ur.Save("joe", "joe", "minichino", "basic", "sendyoulater")
	return cc.JSONBlob(http.StatusOK, []byte(`{"message":"ok"}`))
}

// CheckSession checks if the session is alive
func CheckSession(c echo.Context) error {
	cc := c.(*Context)
	cookie, _ := cc.Cookie("SessionID")
	fmt.Println(cookie.Value)
	res := map[string]interface{}{
		"status":  "ok",
		"message": fmt.Sprintf("Logged in as %v", cookie.Value),
	}

	return c.JSON(http.StatusOK, res)
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

// SaveSMSAction is the handler for saving sms actions
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

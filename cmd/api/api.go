package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/spf13/viper"
	"github.com/techfort/sendyoulater"
	"golang.org/x/oauth2"
)

// Context is the dis Context for API requests
type Context struct {
	echo.Context
	Redis *redis.Client
}

const (
	// ErrMessage std error
	ErrMessage = `{"message":"unable to process request"}`
	// CookieName is the name of the cookie
	CookieName = "SessionID"
)

// BadRequest returns a byte array for a json blob error response
func BadRequest() []byte {
	return []byte(ErrMessage)
}

// Store returns the repo factory
func (c Context) Store() sendyoulater.Store {
	return sendyoulater.NewStore(c.Redis)
}

// SessionID generates an id
func (c Context) SessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return fmt.Sprintf(`syl_sess:%v`, base64.URLEncoding.EncodeToString(b))
}

// Err returns a standard error response
func (c Context) Err() error {
	return c.JSONBlob(http.StatusInternalServerError, BadRequest())
}

// SessionStore instantiates a session store
func (c Context) SessionStore() sendyoulater.SessionStore {
	return sendyoulater.NewSessionStore(c.Redis)
}

// RoutesGET returns get routes
func RoutesGET() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"user":        UserData,
		"/user/:id":   UserByID,
		"/plan/:name": PlanByName,
		"login":       Login,
		"token":       Token,
		"auth":        Auth,
		"check":       CheckSession,
		"loadData":    LoadData,
		"init":        initData,
	}
}

// RoutesPOST returns POST routes
func RoutesPOST() map[string]echo.HandlerFunc {
	return map[string]echo.HandlerFunc{
		"/action/email/save":  SaveEmailAction,
		"/remove":             Remove,
		"/user/save":          SaveUser,
		"user/:userId/update": UpdateUser,
		"loginfromfe":         LoginFromFE,
	}
}

// InitAPI starts the API
func InitAPI(v *viper.Viper) (*echo.Echo, error) {
	e := echo.New()
	fmt.Println("ENV", v.GetString("api_port"), v.GetString("redis_url"))
	r := redis.NewClient(&redis.Options{
		Addr: v.GetString("redis_url"), //"localhost:6379",
	})
	e = Config(e, r)
	err := e.Start(":" + v.GetString("api_port")) //":1323"
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

// LoginFromFE handles the session creation from the frontend
func LoginFromFE(c echo.Context) error {
	cc := c.(*Context)
	var m echo.Map
	if err := c.Bind(&m); err != nil {
		return c.JSONBlob(http.StatusInternalServerError, []byte(fmt.Sprintf(`{ "message": "fail", "error": "%v" }`, err.Error())))
	}

	sessionID := cc.SessionID()
	cookie := new(http.Cookie)
	cookie.Name = CookieName
	cookie.Value = sessionID
	cookie.Expires = time.Now().Add(time.Duration(60) * time.Minute)
	cc.SetCookie(cookie)

	fmt.Println(fmt.Sprintf("SessionID: %v", cookie.Value))
	cc.SessionStore().Set(sessionID, m["Email"].(string))
	return cc.JSON(http.StatusOK, m)
}

var (
	gmailClient *http.Client
	v           = sendyoulater.Env()
	conf        = sendyoulater.Oauth2Config(v)
)

// Token is the callback from google
func Token(c echo.Context) error {
	ctx := context.Background()
	cc := c.(*Context)
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
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	ur := cc.Store().NewUserRepo()
	var (
		user sendyoulater.User
	)
	user, err = ur.ByID(userInfo.Email)
	if user.UserID == "" {
		user, err = ur.Save(userInfo.Email, userInfo.GivenName, userInfo.LastName, "basic", "", tok.AccessToken, tok.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		user.RefreshToken = tok.RefreshToken
		user.Token = tok.AccessToken
		user, err = ur.Update(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	sessionID := cc.SessionID()
	c.SetCookie(&http.Cookie{
		Name:   CookieName,
		Value:  sessionID,
		Secure: true,
	})
	err = cc.SessionStore().Set(sessionID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
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

// UserData returns the user data for the logged in user
func UserData(c echo.Context) error {
	cc := c.(*Context)
	cookie, err := cc.Cookie(CookieName)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error in UserData: %v", err))
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve cookie"))
	}
	sessionID := cookie.Value
	userID, err := cc.SessionStore().Get(sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to retrieve session"))
	}
	ur := cc.Store().NewUserRepo()
	user, err := ur.ByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to find user"))
	}
	return c.JSON(http.StatusOK, user)
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
	pr := cc.Store().NewPlanRepo()
	pr.SavePlan("basic", 100, 100)
	return cc.JSONBlob(http.StatusOK, []byte(`{"message":"ok"}`))
}

// CheckSession checks if the session is alive
func CheckSession(c echo.Context) error {
	cc := c.(*Context)
	cookie, err := cc.Cookie("SessionID")
	if err != nil {
		fmt.Println(err.Error())
		return cc.JSON(http.StatusInternalServerError, err)
	}
	res := map[string]interface{}{
		"status":  "ok",
		"message": fmt.Sprintf("Logged in as %v", cookie.Value),
	}

	return c.JSON(http.StatusOK, res)
}

// LoadData returns all actions for user
func LoadData(c echo.Context) error {
	cc := c.(*Context)
	store := cc.Store()
	ur, er := store.NewUserRepo(), store.NewEmailActionRepo()
	userID := cc.QueryParam("user")
	if userID == "" {
		return cc.JSON(http.StatusInternalServerError, errors.New("Missing userId parameter"))
	}
	user, err := ur.ByID(userID)
	if err != nil {
		return cc.JSON(http.StatusInternalServerError, errors.New("cannot retrieve user"))
	}

	eas, err := er.EmailsOfUser(user)
	if err != nil {
		return cc.Err()
	}
	return cc.JSON(http.StatusOK, eas)
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
		fmt.Println(fmt.Sprintf("Error converting duration: %v", err.Error()))
		return cc.Err()
	}
	ea, err := euc.SaveEmailActions(userID, subject, body, to, ex)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error saving email action: %v", err.Error()))
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
	user, err := ur.Save(m["userId"].(string), m["firstname"].(string), m["lastname"].(string), m["plan"].(string), m["company"].(string), "", "")
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

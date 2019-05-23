package sendyoulater

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	urlshortener "google.golang.org/api/urlshortener/v1"
)

// Env returns a viper object with all the env vars
func Env() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("syl")
	v.BindEnv("oauth_clientID")
	v.BindEnv("oauth_client_secret")
	v.BindEnv("redis_url")
	err := v.BindEnv("api_port")
	if err != nil {
		fmt.Println("ERROR setting env var", err.Error())
	}
	v.BindEnv("oauth_redirect_url")
	return v
}

// Oauth2Config returns the oauth config object
func Oauth2Config(v *viper.Viper) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     v.GetString("oauth_clientID"),      //"541640626027-75fep8r1ptdd377l73qhb2f03pckc0po.apps.googleusercontent.com",
		ClientSecret: v.GetString("oauth_client_secret"), //"7Qbre2sDMxnPHZwa_6DASz4j",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			urlshortener.UrlshortenerScope,
			gmail.GmailSendScope,
			"openid",
			"profile",
			"email",
		},
		RedirectURL: v.GetString("oauth_redirect_url"), //"http://localhost:1323/token",
	}
}

// SetUpStripe sets the env var to the stripe key
func SetUpStripe(v *viper.Viper) error {
	stripeKey := v.GetString("stripe_secret_key")
	if stripeKey == "" {
		return errors.New("No secret key found for stripe")
	}
	stripe.Key = stripeKey
	return nil
}

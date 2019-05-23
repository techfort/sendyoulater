package sendyoulater

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/api/gmail/v1"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// UserIDFromKey returns the userID from an action key
func UserIDFromKey(key string) (string, error) {
	fail := ""
	chunks := strings.Split(key, ":")
	if len(chunks) != 4 {
		return fail, errors.Errorf("Error parsing key: %v", key)
	}
	if chunks[1] == "" {
		return fail, errors.Errorf("Missing user ID in key: %v", key)
	}
	return chunks[1], nil
}

// RefreshTokenResponse is a struct holding the data of a refresh token request
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

const (
	// RefreshTokenURL is the url for the refresh token request
	RefreshTokenURL = "https://www.googleapis.com/oauth2/v4/token"
)

// RefreshToken attempts to refresh the token from the google api
func RefreshToken(client *http.Client, conf *oauth2.Config, oldToken *oauth2.Token) (RefreshTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", oldToken.RefreshToken)
	data.Set("client_id", conf.ClientID)
	data.Set("client_secret", conf.ClientSecret)
	fail := RefreshTokenResponse{}
	req, err := http.NewRequest("POST", RefreshTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fail, errors.Wrap(err, "failed to create request for refresh token")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return fail, errors.Wrap(err, "error executing token refresh request")
	}
	var rtr RefreshTokenResponse
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fail, errors.Wrap(err, "failed to read body of response from refresh token request")
	}
	if err := json.Unmarshal(body, &rtr); err != nil {
		return fail, errors.Wrap(err, "failed to unmrashal json of response")
	}
	return rtr, err
}

// Message takes a few parameters and forms a gmail message
func Message(from, to, subject, body string) gmail.Message {
	var message gmail.Message
	messageStr := []byte(
		"From: " + from + "\r\n" +
			"Reply-to: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" + body)

	message.Raw = base64.StdEncoding.EncodeToString(messageStr)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)
	return message
}

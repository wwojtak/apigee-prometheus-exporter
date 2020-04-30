package token

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// This is a hard-coded value that the API requires in the Authorization header.
// It is NOT client credential.

const clientAuth = "ZWRnZWNsaTplZGdlY2xpc2VjcmV0"

var username string
var password string

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ApigeeToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expires      int64
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}

func ApigeeClient(storedToken *ApigeeToken) (*ApigeeToken, error) {
	postData := url.Values{}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	if storedToken.RefreshToken != "" {
		log.Info("Refreshing token")
		postData.Set("grant_type", "refresh_token")
		postData.Set("refresh_token", storedToken.RefreshToken)
	} else {
		log.Info("Requesting new token")
		username = os.Getenv("APIGEE_USERNAME")
		password = os.Getenv("APIGEE_PASSWORD")
		if username == "" || password == "" {
			log.Fatal("Apigee credentials not provided")
			return nil, errors.New("token request failed")
		}
		postData.Set("username", username)
		postData.Set("password", password)
		postData.Set("grant_type", "password")
	}
	request, err := http.NewRequest("POST", "https://login.apigee.com/oauth/token",
		strings.NewReader(postData.Encode()))
	if err != nil {
		log.WithFields(log.Fields{"details": err}).Fatal("Creating request to token endpoint failed")
		return nil, errors.New("token request failed")
	} else {
		request.Header.Set("Accept", "application/json;charset=utf-8")
		request.Header.Set("Authorization", "Basic "+clientAuth)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		resp, err := client.Do(request)
		if err != nil {
			log.WithFields(log.Fields{"details": err}).Fatal("Call to token endpoint failed")
			return nil, errors.New("token request failed")
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				errResp := new(ErrorResponse)
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					log.WithFields(log.Fields{"details": err}).Fatal("Cannot decode JSON response")
					return nil, errors.New("token request failed")
				} else {
					log.WithFields(log.Fields{"status_code": resp.StatusCode, "details": errResp}).Fatal(
						"Call to token endpoint failed")
					return nil, errors.New("token request failed")
				}
			} else {
				tokenResp := new(ApigeeToken)
				err = json.NewDecoder(resp.Body).Decode(&tokenResp)
				tokenResp.Expires = time.Now().Unix() + int64(tokenResp.ExpiresIn)
				if err != nil {
					log.WithFields(log.Fields{"details": err}).Fatal(
						"Cannot decode JSON response")
					return nil, errors.New("token request failed")
				} else {
					log.Info("New token received")
					return tokenResp, nil
				}
			}
		}
	}
}

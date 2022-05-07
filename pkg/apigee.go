package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	SSOClientCredentials = "ZWRnZWNsaTplZGdlY2xpc2VjcmV0"
)

type ApigeeConfig struct {
	Username      string
	Password      string
	AccessToken   string
	Mfa           string
	OauthTokenUrl string
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func ApigeeLogin(config ApigeeConfig) (string, error) {
	requestForm := url.Values{
		"grant_type": []string{"password"},
		"username":   []string{config.Username},
		"password":   []string{config.Password},
	}
	req, err := http.NewRequest(http.MethodPost, config.OauthTokenUrl, bytes.NewBufferString(requestForm.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic"+" "+SSOClientCredentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := url.Values{}
	q.Add("mfa_token", config.Mfa)
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return "", &RequestError{StatusCode: resp.StatusCode, Err: err}
		}
		return "", &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}

	token := &OauthToken{}
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return "", err
	}
	fmt.Println(token.RefreshToken)

	return token.AccessToken, nil
}

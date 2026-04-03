// Package linkedinoauth2 handles the OAuth 2.0 with Google
package linkedinoauth2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/SyafaHadyan/worku/internal/constants"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

type LinkedInOAuthItf interface {
	LinkedInOAuthConfig() *oauth2.Config
	GetUserInfo(token string) (dto.ResponseLinkedInOAuth, error)
	GenerateRandomState() string
}

type LinkedInOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func New(env *env.Env) *LinkedInOAuth {
	return &LinkedInOAuth{
		ClientID:     env.LinkedInClientID,
		ClientSecret: env.LinkedInClientSecret,
		RedirectURL:  env.LinkedInRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     linkedin.Endpoint,
	}
}

func (l *LinkedInOAuth) LinkedInOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     l.ClientID,
		ClientSecret: l.ClientSecret,
		RedirectURL:  l.RedirectURL,
		Scopes:       l.Scopes,
		Endpoint:     l.Endpoint,
	}
}

func (l *LinkedInOAuth) GetUserInfo(token string) (dto.ResponseLinkedInOAuth, error) {
	var data dto.ResponseLinkedInOAuth

	reqURL, err := url.Parse(string(constants.LinkedInOAuthGetUserInfo))
	if err != nil {
		return dto.ResponseLinkedInOAuth{}, err
	}

	pToken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {pToken},
		},
	}

	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return dto.ResponseLinkedInOAuth{}, err
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return dto.ResponseLinkedInOAuth{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return dto.ResponseLinkedInOAuth{}, err
	}

	log.Println(data)

	return data, nil
}

func (l *LinkedInOAuth) GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	return state
}

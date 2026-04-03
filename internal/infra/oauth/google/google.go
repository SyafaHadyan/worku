// Package google handles the OAuth 2.0 with Google
package googleoauth2

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/SyafaHadyan/worku/internal/constants"
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthItf interface {
	GoogleOAuthConfig() *oauth2.Config
	GetUserInfo(token string) (dto.ResponseGoogleOAuth, error)
	GenerateRandomState() string
}

type GoogleOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func New(env *env.Env) *GoogleOAuth {
	return &GoogleOAuth{
		ClientID:     env.GoogleClientID,
		ClientSecret: env.GoogleClientSecret,
		RedirectURL:  env.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func (g *GoogleOAuth) GoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURL,
		Scopes:       g.Scopes,
		Endpoint:     g.Endpoint,
	}
}

func (g *GoogleOAuth) GetUserInfo(token string) (dto.ResponseGoogleOAuth, error) {
	var data dto.ResponseGoogleOAuth

	reqURL, err := url.Parse(string(constants.GoogleOAuthGetUserInfo))
	if err != nil {
		return dto.ResponseGoogleOAuth{}, err
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
		return dto.ResponseGoogleOAuth{}, err
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return dto.ResponseGoogleOAuth{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return dto.ResponseGoogleOAuth{}, err
	}

	return data, nil
}

func (g *GoogleOAuth) GenerateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	return state
}

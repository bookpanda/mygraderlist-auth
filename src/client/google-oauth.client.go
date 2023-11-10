package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type GoogleOauthClient struct {
	oauthConfig *oauth2.Config
}

func NewGoogleOauthClient(oauthConfig *oauth2.Config) *GoogleOauthClient {
	return &GoogleOauthClient{
		oauthConfig,
	}
}

type GoogleUserEmailResponse struct {
	Email     string `json:"email"`
	Firstname string `json:"given_name"`
	Lastname  string `json:"family_name"`
}

var (
	InvalidCode   = errors.New("Invalid code")
	HttpError     = errors.New("Unable to get user info")
	IOError       = errors.New("Unable to read google response")
	InvalidFormat = errors.New("Google sent unexpected format")
)

func (c *GoogleOauthClient) GetUserEmail(code string) (*GoogleUserEmailResponse, error) {
	token, err := c.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Error().Err(err).Msg("Unable to exchange oauth token")
		return nil, InvalidCode
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Error().Err(err).Msg("Unable to get user info")
		return nil, HttpError
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read user info response")
		return nil, IOError
	}

	var parsedResponse GoogleUserEmailResponse
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		log.Error().Err(err).Msg("Google send unexpected response")
		return nil, InvalidFormat
	}

	return &parsedResponse, nil
}

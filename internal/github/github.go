package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Service holds GitHub OAuth configuration.
type Service struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	OAuthAPIURL  string `json:"oauth_api_url"`
	UserAPIURL   string `json:"user_api_url"`
}

// NewService creates a new GitHub Service instance with the given client ID and client secret.
func NewService(clientID, clientSecret string) *Service {
	return &Service{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		OAuthAPIURL:  "https://github.com/login/oauth/access_token",
		UserAPIURL:   "https://api.github.com/user",
	}
}

// ExchangeCodeForToken exchanges the given code from
// https://github.com/login/oauth/authorize?client_id=&redirect_uri=
// for an access token.
func (g *Service) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	formData := url.Values{
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"code":          {code},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		g.OAuthAPIURL,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose
	if err != nil {
		return "", fmt.Errorf("error doing request: %w", err)
	}
	defer func(b io.ReadCloser) {
		e := b.Close()
		if e != nil {
			fmt.Printf("error closing response body: %v\n", e)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error parsing response body: %w", err)
	}

	result := exchangeCodeResponse{
		AccessToken:      values.Get("access_token"),
		Error:            values.Get("error"),
		ErrorDescription: values.Get("error_description"),
		ErrorURI:         values.Get("error_uri"),
	}

	if result.Error != "" {
		return "", fmt.Errorf("error from github: %s", result.Error)
	}

	return result.AccessToken, nil
}

// GetUserInfo returns user info from GitHub using the given access token from Service.ExchangeCodeForToken.
func (g *Service) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		g.UserAPIURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer func(b io.ReadCloser) {
		e := b.Close()
		if e != nil {
			fmt.Printf("error closing response body: %v\n", e)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from github: %s", body)
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return &userInfo, nil
}

// UserInfo is a struct to hold GitHub user data.
type UserInfo struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type exchangeCodeResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}

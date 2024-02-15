package github

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_NewGitHub(t *testing.T) {
	g := NewGitHub("dummyClientID", "dummyClientSecret")

	assert.Equal(t, "dummyClientID", g.ClientID)
	assert.Equal(t, "dummyClientSecret", g.ClientSecret)
	assert.Equal(t, "https://github.com/login/oauth/access_token", g.OAuthAPIURL)
	assert.Equal(t, "https://api.github.com/user", g.UserAPIURL)
}

func Test_GetUserInfo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Send response to be tested
			_, err := rw.Write([]byte(`{
			"login": "octocat",
			"id": 1,
			"node_id": "MDQ6VXNlcjE=",
			"avatar_url": "https://github.com/images/error/octocat_happy.gif",
			"gravatar_id": "",
			"url": "https://api.github.com/users/octocat",
			"html_url": "https://github.com/octocat",
			"followers_url": "https://api.github.com/users/octocat/followers",
			"following_url": "https://api.github.com/users/octocat/following{/other_user}",
			"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
			"organizations_url": "https://api.github.com/users/octocat/orgs",
			"repos_url": "https://api.github.com/users/octocat/repos",
			"events_url": "https://api.github.com/users/octocat/events{/privacy}",
			"received_events_url": "https://api.github.com/users/octocat/received_events",
			"type": "User",
			"site_admin": false,
			"name": "monalisa octocat",
			"company": "GitHub",
			"blog": "https://github.com/blog",
			"location": "San Francisco",
			"email": "octocat@github.com",
			"hireable": false,
			"bio": "There once was...",
			"twitter_username": "monatheoctocat",
			"public_repos": 2,
			"public_gists": 1,
			"followers": 20,
			"following": 0,
			"created_at": "2008-01-14T04:33:35Z",
			"updated_at": "2008-01-14T04:33:35Z"
		}`))
			if err != nil {
				return
			}
		}))
		// Close the server when test finishes
		defer server.Close()

		// Create a new GitHub instance
		g := GitHub{
			ClientID:     "dummyClientID",
			ClientSecret: "dummyClientSecret",
			UserAPIURL:   server.URL,
		}

		// Call GetUserInfo
		userInfo, err := g.GetUserInfo(context.Background(), "dummyToken")

		// Assert no error
		assert.NoError(t, err)

		// Assert the returned user info
		// Assert the returned user info
		assert.Equal(t, "octocat", userInfo.Login)
		assert.Equal(t, 1, userInfo.ID)
		assert.Equal(t, "MDQ6VXNlcjE=", userInfo.NodeID)
		assert.Equal(t, "https://github.com/images/error/octocat_happy.gif", userInfo.AvatarURL)
		assert.Equal(t, "", userInfo.GravatarID)
		assert.Equal(t, "https://api.github.com/users/octocat", userInfo.URL)
		assert.Equal(t, "https://github.com/octocat", userInfo.HTMLURL)
		assert.Equal(t, "https://api.github.com/users/octocat/followers", userInfo.FollowersURL)
		assert.Equal(t, "https://api.github.com/users/octocat/following{/other_user}", userInfo.FollowingURL)
		assert.Equal(t, "https://api.github.com/users/octocat/gists{/gist_id}", userInfo.GistsURL)
		assert.Equal(t, "https://api.github.com/users/octocat/starred{/owner}{/repo}", userInfo.StarredURL)
		assert.Equal(t, "https://api.github.com/users/octocat/subscriptions", userInfo.SubscriptionsURL)
		assert.Equal(t, "https://api.github.com/users/octocat/orgs", userInfo.OrganizationsURL)
		assert.Equal(t, "https://api.github.com/users/octocat/repos", userInfo.ReposURL)
		assert.Equal(t, "https://api.github.com/users/octocat/events{/privacy}", userInfo.EventsURL)
		assert.Equal(t, "https://api.github.com/users/octocat/received_events", userInfo.ReceivedEventsURL)
		assert.Equal(t, "User", userInfo.Type)
		assert.Equal(t, false, userInfo.SiteAdmin)
		assert.Equal(t, "monalisa octocat", userInfo.Name)
		assert.Equal(t, "GitHub", userInfo.Company)
		assert.Equal(t, "https://github.com/blog", userInfo.Blog)
		assert.Equal(t, "San Francisco", userInfo.Location)
		assert.Equal(t, "octocat@github.com", userInfo.Email)
		assert.Equal(t, false, userInfo.Hireable)
		assert.Equal(t, "There once was...", userInfo.Bio)
		assert.Equal(t, "monatheoctocat", userInfo.TwitterUsername)
		assert.Equal(t, 2, userInfo.PublicRepos)
		assert.Equal(t, 1, userInfo.PublicGists)
		assert.Equal(t, 20, userInfo.Followers)
		assert.Equal(t, 0, userInfo.Following)
		assert.Equal(t, time.Date(2008, 1, 14, 4, 33, 35, 0, time.UTC), userInfo.CreatedAt)
		assert.Equal(t, time.Date(2008, 1, 14, 4, 33, 35, 0, time.UTC), userInfo.UpdatedAt)
	})

	t.Run("Error 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusUnauthorized)
			_, err := rw.Write([]byte(`{
				"message": "Bad credentials",
				"documentation_url": "https://docs.github.com/rest"
		}`))
			if err != nil {
				return
			}
		}))
		// Close the server when test finishes
		defer server.Close()

		// Create a new GitHub instance
		g := GitHub{
			ClientID:     "dummyClientID",
			ClientSecret: "dummyClientSecret",
			UserAPIURL:   server.URL,
		}

		// Call GetUserInfo
		userInfo, err := g.GetUserInfo(context.Background(), "dummyToken")
		if err == nil {
			t.Fatal("expected an error")
		}

		// Assert the error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error response from github")
		assert.Nil(t, userInfo)
	})
}

func Test_ExchangeCodeForToken(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, err := rw.Write([]byte(`access_token=gho_1fskFhk1sh2kfs3dh4fh2&scope=&token_type=bearer`))

			if err != nil {
				return
			}
		}))
		defer server.Close()

		g := GitHub{
			ClientID:     "dummyClientID",
			ClientSecret: "dummyClientSecret",
			OAuthAPIURL:  server.URL,
		}

		token, err := g.ExchangeCodeForToken(context.Background(), "dummyCode")

		assert.NoError(t, err)
		assert.Equal(t, "gho_1fskFhk1sh2kfs3dh4fh2", token)
	})

	t.Run("error bad_verification_code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, err := rw.Write([]byte(`error=bad_verification_code&error_description=The+code+passed+is+incorrect+or+expired.&error_uri=https%3A%2F%2Fdocs.github.com`))
			if err != nil {
				return
			}
		}))
		defer server.Close()

		g := GitHub{
			ClientID:     "dummyClientID",
			ClientSecret: "dummyClientSecret",
			OAuthAPIURL:  server.URL,
		}

		token, err := g.ExchangeCodeForToken(context.Background(), "dummyCode")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "error from github: bad_verification_code")
	})
}

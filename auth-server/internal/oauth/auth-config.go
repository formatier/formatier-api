package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

var GoogleAuthConfig = oauth2.Config{
	RedirectURL:  os.Getenv("API_URL"),
	ClientID:     os.Getenv("GOOGLE_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET"),
	Endpoint:     google.Endpoint,
}

var GithubAuthConfig = oauth2.Config{
	RedirectURL:  os.Getenv("API_URL"),
	ClientID:     os.Getenv("GITHUB_ID"),
	ClientSecret: os.Getenv("GITHUB_SECRET"),
	Endpoint:     github.Endpoint,
}

var MicrosoftAuthConfig = oauth2.Config{
	RedirectURL:  os.Getenv("API_URL"),
	ClientID:     os.Getenv("MICROSOFT_ID"),
	ClientSecret: os.Getenv("MICROSOFT_SECRET"),
	Endpoint:     microsoft.LiveConnectEndpoint,
}

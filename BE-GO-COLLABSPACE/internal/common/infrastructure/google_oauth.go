package infrastructure

import (
	"Go-CollabSpace/internal/common/apperror"
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	oauthConfig *oauth2.Config
}

func NewGoogleOAuth(clientID, clientSecret, redirectURL string) OAuthProvider {
	return &GoogleOAuth{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
			//Endpoint: oauth2.Endpoint{
			//	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			//	TokenURL: "https://oauth2.googleapis.com/token",
			//},
		},
	}
}

func (g GoogleOAuth) GetProfile(ctx context.Context, code string) (*SocialProfile, error) {
	token, err := g.oauthConfig.Exchange(ctx, code) // Exchange the authorization code for an access token
	if err != nil {
		return nil, apperror.ExchangeTokenFailed
	}

	client := g.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, apperror.ErrFetchProfileFailed
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperror.ErrFetchProfileFailed
	}

	var googleData struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleData); err != nil {
		return nil, apperror.ErrFetchProfileFailed
	}

	return &SocialProfile{
		ID:     googleData.ID,
		Email:  googleData.Email,
		Name:   googleData.Name,
		Avatar: googleData.Picture,
	}, nil
}

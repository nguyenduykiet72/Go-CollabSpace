package infrastructure

import "context"

type SocialProfile struct {
	ID     string
	Email  string
	Name   string
	Avatar string
}

type OAuthProvider interface {
	GetProfile(ctx context.Context, code string) (*SocialProfile, error)
}

package client

import "context"

type ClientInterface interface {
	GetLogins(ctx context.Context, url string) ([]Login, error)
	Associate(ctx context.Context) (AssociateResponse, error)
}

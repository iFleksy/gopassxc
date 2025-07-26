package client

type ClientInterface interface {
	GetLogins(url string) ([]Login, error)
}

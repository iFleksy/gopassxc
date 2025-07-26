package storage

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound       = errors.New("profile not found")
	ErrAlreadyExists  = errors.New("profile already exists")
	ErrInvalidProfile = errors.New("invalid profile")
	ErrInvalidStorage = errors.New("invalid storage")
)

type Storage interface {
	AddProfile(p Profile) error
	ExtractProfile(name string) (Profile, error)
	DeleteProfile(name string) error
}

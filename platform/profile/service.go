package profile

import (
	"context"
)

type Service interface {
	ApplyProfile(ctx context.Context, p *Profile) error
	GetProfiles(ctx context.Context, opt GetProfilesOption) ([]Profile, error)
	RemoveProfiles(ctx context.Context, ids []string) error
}

type GetProfilesOption struct {
	Identifier string `json:"id"`
}

type Store interface {
	ProfileById(id string) (*Profile, error)
	Save(p *Profile) error
	List() ([]Profile, error)
	Delete(id string) error
}

func New(store Store) *ProfileService {
	return &ProfileService{store: store}
}

type ProfileService struct {
	store Store
}

func IsNotFound(err error) bool {
	type notFoundError interface {
		error
		NotFound() bool
	}

	_, ok := err.(notFoundError)
	return ok
}

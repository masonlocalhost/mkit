package service

import (
	"context"
	"mkit/example/depinjection/model"
	"mkit/example/depinjection/repository"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{
		userRepo,
	}
}

func (s *UserService) ListSomeUsers(ctx context.Context) ([]*model.User, error) {
	return s.userRepo.FindUsersByFilters(ctx, nil, nil)
}

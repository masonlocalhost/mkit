package main

import (
	"context"
	"database/sql"
	"mkit/example/depinjection/model"
	"mkit/example/depinjection/repository"
	"mkit/example/depinjection/service"
)

func main() {
	var (
		db  *sql.DB // example db
		ctx = context.Background()
	)

	// use user repo
	userRepo := repository.NewUserRepo(db)
	userService1 := service.NewUserService(userRepo)
	userService1.ListSomeUsers(ctx) // get result from real db

	// use custom user repo
	userService2 := service.NewUserService(&customUserRepo{})

	userService2.ListSomeUsers(ctx) // get result from mock user repo
}

// custom user repo implementation
type customUserRepo struct{}

func (r *customUserRepo) FindUsersByFilters(ctx context.Context, names []string, ids []string) ([]*model.User, error) {
	// Logic for custom filtering would go here
	return []*model.User{
		{ID: "101", Name: "Custom_Alpha"},
		{ID: "102", Name: "Custom_Beta"},
	}, nil
}

func (r *customUserRepo) FirstUserByID(ctx context.Context, id string) (*model.User, error) {
	// Return a user specific to this custom implementation
	return &model.User{ID: id, Name: "CustomUser_" + id}, nil
}

func (r *customUserRepo) Create(ctx context.Context, name string) (*model.User, error) {
	// In a real app, you'd generate a UUID or use a DB auto-increment
	return &model.User{ID: "999", Name: name}, nil
}

func (r *customUserRepo) UpdateByID(ctx context.Context, id string, newName string) (*model.User, error) {
	// Mocking an update success
	return &model.User{ID: id, Name: newName}, nil
}

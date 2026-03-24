package repository

import (
	"context"
	"database/sql"
	"mkit/example/depinjection/model"
)

type UserRepo interface {
	FindUsersByFilters(ctx context.Context, names []string, ids []string) ([]*model.User, error)
	FirstUserByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, name string) (*model.User, error)
	UpdateByID(ctx context.Context, id string, newName string) (*model.User, error)
}

func NewUserRepo(db *sql.DB) UserRepo { // return interface instead of implementation struct
	return &impl{ // impl satisfy UserRepo interface
		db: db, // inject db as a dependency
	}
}

type impl struct {
	db *sql.DB
}

func (r *impl) FindUsersByFilters(ctx context.Context, names []string, ids []string) ([]*model.User, error) {
	// query and return users by filters
	return []*model.User{{ID: "1", Name: "Mot"}, {ID: "2", Name: "Hai"}}, nil
}

func (r *impl) FirstUserByID(ctx context.Context, id string) (*model.User, error) {
	// mock: searching for a specific user
	return &model.User{ID: id, Name: "Mock User"}, nil
}

func (r *impl) Create(ctx context.Context, name string) (*model.User, error) {
	// mock: creating a new user with a generated ID
	return &model.User{ID: "3", Name: name}, nil
}

func (r *impl) UpdateByID(ctx context.Context, id string, newName string) (*model.User, error) {
	// mock: updating an existing user
	return &model.User{ID: id, Name: newName}, nil
}

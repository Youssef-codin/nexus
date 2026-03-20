package users

import (
	"context"

	"github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
)

type iuserRepo interface {
	GetUserByName(ctx context.Context, fullName string) ([]repo.User, error)
}

type UserRepo struct {
	db *db.DB
}

func NewUserRepo(database *db.DB) iuserRepo {
	return &UserRepo{db: database}
}

func (r *UserRepo) GetUserByName(ctx context.Context, fullName string) ([]repo.User, error) {
	return db.Queries(ctx, r.db).GetUserByName(ctx, fullName)
}

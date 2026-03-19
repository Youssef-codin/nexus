package users

import (
	"context"

	"github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type iuserRepo interface {
	GetUserByName(ctx context.Context, fullName string) ([]repo.User, error)
}

type UserRepo struct {
	pool *pgx.Conn
}

func NewUserRepo(pool *pgx.Conn) iuserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetUserByName(ctx context.Context, fullName string) ([]repo.User, error) {
	return db.Queries(ctx, r.pool).GetUserByName(ctx, fullName)
}

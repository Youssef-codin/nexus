package auth

import (
	"context"

	"github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type iauthRepo interface {
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.CreateUserRow, error)
	GetUserByRefreshToken(ctx context.Context, token pgtype.Text) (repo.User, error)
	UpdateRefreshToken(ctx context.Context, arg repo.UpdateRefreshTokenParams) error
	GetUserById(ctx context.Context, id pgtype.UUID) (repo.User, error)
	RevokeRefreshToken(ctx context.Context, id pgtype.UUID) error
}

type AuthRepo struct {
	pool *pgx.Conn
}

func NewAuthRepo(pool *pgx.Conn) iauthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return db.Queries(ctx, r.pool).GetUserByEmail(ctx, email)
}

func (r *AuthRepo) CreateUser(
	ctx context.Context,
	arg repo.CreateUserParams,
) (repo.CreateUserRow, error) {
	return db.Queries(ctx, r.pool).CreateUser(ctx, arg)
}

func (r *AuthRepo) GetUserByRefreshToken(
	ctx context.Context,
	token pgtype.Text,
) (repo.User, error) {
	return db.Queries(ctx, r.pool).GetUserByRefreshToken(ctx, token)
}

func (r *AuthRepo) UpdateRefreshToken(
	ctx context.Context,
	arg repo.UpdateRefreshTokenParams,
) error {
	return db.Queries(ctx, r.pool).UpdateRefreshToken(ctx, arg)
}

func (r *AuthRepo) GetUserById(ctx context.Context, id pgtype.UUID) (repo.User, error) {
	return db.Queries(ctx, r.pool).GetUserById(ctx, id)
}

func (r *AuthRepo) RevokeRefreshToken(ctx context.Context, id pgtype.UUID) error {
	return db.Queries(ctx, r.pool).RevokeRefreshToken(ctx, id)
}

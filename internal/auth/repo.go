package auth

import (
	"context"

	"github.com/Youssef-codin/NexusPay/internal/db"
	repo "github.com/Youssef-codin/NexusPay/internal/db/postgresql/sqlc"
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
	db *db.DB
}

func NewAuthRepo(database *db.DB) iauthRepo {
	return &AuthRepo{db: database}
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return db.Queries(ctx, r.db).GetUserByEmail(ctx, email)
}

func (r *AuthRepo) CreateUser(
	ctx context.Context,
	arg repo.CreateUserParams,
) (repo.CreateUserRow, error) {
	return db.Queries(ctx, r.db).CreateUser(ctx, arg)
}

func (r *AuthRepo) GetUserByRefreshToken(
	ctx context.Context,
	token pgtype.Text,
) (repo.User, error) {
	return db.Queries(ctx, r.db).GetUserByRefreshToken(ctx, token)
}

func (r *AuthRepo) UpdateRefreshToken(
	ctx context.Context,
	arg repo.UpdateRefreshTokenParams,
) error {
	return db.Queries(ctx, r.db).UpdateRefreshToken(ctx, arg)
}

func (r *AuthRepo) GetUserById(ctx context.Context, id pgtype.UUID) (repo.User, error) {
	return db.Queries(ctx, r.db).GetUserById(ctx, id)
}

func (r *AuthRepo) RevokeRefreshToken(ctx context.Context, id pgtype.UUID) error {
	return db.Queries(ctx, r.db).RevokeRefreshToken(ctx, id)
}

package repository

import (
	"context"
	"errors"

	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (ur *UserRepository) GetUserByOAuthId(ctx context.Context, oauthId string) (*domain.User, error) {
	query := `SELECT id, oauth_id, picture, name, email, custom_url FROM users WHERE oauth_id = $1`

	user := new(domain.User)

	err := ur.pool.QueryRow(ctx, query, oauthId).Scan(
		&user.ID,
		&user.OAuthId,
		&user.Picture,
		&user.Name,
		&user.Email,
		&user.CustomUrl,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (ur *UserRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, oauth_id, picture, name, email, custom_url FROM users WHERE id = $1`

	user := new(domain.User)

	err := ur.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.OAuthId,
		&user.Picture,
		&user.Name,
		&user.Email,
		&user.CustomUrl,
	)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (oauth_id, picture, name, email, custom_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := ur.pool.QueryRow(ctx, query, user.OAuthId, user.Picture, user.Name, user.Email, user.CustomUrl).Scan(&user.ID)

	return err
}

func (ur *UserRepository) CreateOrGetUser(ctx context.Context, dUser *domain.User) (*domain.User, error) {
	user, err := ur.GetUserByOAuthId(ctx, dUser.OAuthId)

	if err == nil {
		return user, nil
	}

	err = ur.CreateUser(ctx, dUser)

	if err != nil {
		return nil, err
	}

	return dUser, nil
}

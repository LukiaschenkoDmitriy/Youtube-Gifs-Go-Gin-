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

func (ur *UserRepository) GetUserByOAuthId(oauthId string) (*domain.User, error) {
	query := `SELECT id, oauth_id, picture, name, email, custom_url FROM users WHERE oauth_id = $1`

	user := new(domain.User)

	err := ur.pool.QueryRow(context.Background(), query, oauthId).Scan(
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

func (ur *UserRepository) GetUserById(id string) (*domain.User, error) {
	query := `SELECT id, oauth_id, picture, name, email, custom_url FROM users WHERE id = $1`

	user := new(domain.User)

	err := ur.pool.QueryRow(context.Background(), query, id).Scan(
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

func (ur *UserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (oauth_id, picture, name, email, custom_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := ur.pool.QueryRow(context.Background(), query, user.OAuthId, user.Picture, user.Name, user.Email, user.CustomUrl).Scan(&user.ID)

	return err
}

func (ur *UserRepository) CreateOrGetUser(dUser *domain.User) (*domain.User, error) {
	user, err := ur.GetUserByOAuthId(dUser.OAuthId)

	if err == nil {
		return user, nil
	}

	err = ur.CreateUser(dUser)

	if err != nil {
		return nil, err
	}

	return dUser, nil
}

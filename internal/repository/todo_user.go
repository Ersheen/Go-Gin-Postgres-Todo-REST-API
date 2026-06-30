package repository

import (
	"context"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, mail string, password string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at, updated_at
	`

	var user models.User

	err := pool.QueryRow(ctx, query, mail, password).Scan(
		&user.ID,
		&user.Mail,
		&user.Password,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByMail(pool *pgxpool.Pool, mail string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at FROM users WHERE email=$1
	`

	var UserDetails models.User
	err := pool.QueryRow(ctx, query, mail).Scan(
		&UserDetails.ID,
		&UserDetails.Mail,
		&UserDetails.Password,
		&UserDetails.Created_at,
		&UserDetails.Updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &UserDetails, nil
}

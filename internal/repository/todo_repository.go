package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, UserId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO todo_user (title, completed, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, title, completed, created_at, updated_at, user_id
	`
	var todo models.Todo

	err := pool.QueryRow(ctx, query, title, completed, UserId).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool, status bool, UserId string) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, completed, created_at, updated_at, user_id from todo_user
	WHERE user_id=$1 ORDER BY created_at DESC`

	var rows, err = pool.Query(ctx, query, UserId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserId,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil

}

func GetTodo(pool *pgxpool.Pool, id int, UserId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id FROM todo_user WHERE id = $1 AND user_id=$2
	`

	var todo models.Todo

	err := pool.QueryRow(ctx, query, id, UserId).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func buildTodoUpdateQuery(title *string, completed *bool) (string, []any, error) {
	if title == nil && completed == nil {
		return "", nil, errors.New("no fields to update")
	}

	fields := []string{}
	args := []any{}
	paramIndex := 1

	if title != nil {
		fields = append(fields, fmt.Sprintf("title=$%d", paramIndex))
		args = append(args, *title)
		paramIndex++
	}

	if completed != nil {
		fields = append(fields, fmt.Sprintf("completed=$%d", paramIndex))
		args = append(args, *completed)
		paramIndex++
	}

	return strings.Join(fields, ", "), args, nil
}

func UpdateTodo(pool *pgxpool.Pool, title *string, completed *bool, id int, userID string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setClauses, args, err := buildTodoUpdateQuery(title, completed)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		UPDATE todo_user
		SET %s, updated_at=CURRENT_TIMESTAMP
		WHERE id=$%d AND user_id=$%d
		RETURNING id, title, completed, created_at, updated_at, user_id
	`, setClauses, len(args)+1, len(args)+2)

	args = append(args, id, userID)

	var todo models.Todo

	err = pool.QueryRow(ctx, query, args...).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, UserId string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM todo_user WHERE id=$1 AND user_id=$2
	`
	// var todo models.Todo

	result, err := pool.Exec(ctx, query, id, UserId)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil

}

func GetTodoByQuery(pool *pgxpool.Pool, status bool, UserId string) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, title, completed, created_at, updated_at, user_id FROM todo_user WHERE completed=$1 AND user_id=$2
	`

	var rows, err = pool.Query(ctx, query, status, UserId)

	if err != nil {
		return nil, err
	}

	var todos []models.Todo = []models.Todo{}

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserId,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return nil, errors.New("no todo found")
	}

	return todos, nil
}

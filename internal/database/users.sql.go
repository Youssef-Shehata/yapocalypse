// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
returning id, created_at, updated_at, email, password, premuim, username
`

type CreateUserParams struct {
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Premuim,
		&i.Username,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, password, premuim, username from users where email= $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Premuim,
		&i.Username,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, email, password, premuim, username from users where id= $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Premuim,
		&i.Username,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, created_at, updated_at, email, password, premuim, username from users where username= $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Premuim,
		&i.Username,
	)
	return i, err
}

const resetUser = `-- name: ResetUser :exec
delete  from  users
`

func (q *Queries) ResetUser(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetUser)
	return err
}

const subscribeToPremuim = `-- name: SubscribeToPremuim :exec
UPDATE users set premium = true where id = $1
`

func (q *Queries) SubscribeToPremuim(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, subscribeToPremuim, id)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users set email =$1 ,password = $2 WHERE id = $3 returning id, created_at, updated_at, email, password, premuim, username
`

type UpdateUserParams struct {
	Email    string
	Password string
	ID       uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.Email, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Premuim,
		&i.Username,
	)
	return i, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: yaps.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteYap = `-- name: DeleteYap :exec
DELETE from  Yaps where id = $1 and user_id = $2
`

type DeleteYapParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteYap(ctx context.Context, arg DeleteYapParams) error {
	_, err := q.db.ExecContext(ctx, deleteYap, arg.ID, arg.UserID)
	return err
}

const getYapById = `-- name: GetYapById :one
select id, created_at, updated_at, body, user_id from Yaps where id = $1
`

func (q *Queries) GetYapById(ctx context.Context, id uuid.UUID) (Yap, error) {
	row := q.db.QueryRowContext(ctx, getYapById, id)
	var i Yap
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getYaps = `-- name: GetYaps :many
select id, created_at, updated_at, body, user_id from Yaps
`

func (q *Queries) GetYaps(ctx context.Context) ([]Yap, error) {
	rows, err := q.db.QueryContext(ctx, getYaps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Yap
	for rows.Next() {
		var i Yap
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getYapsByUserId = `-- name: GetYapsByUserId :many
select id, created_at, updated_at, body, user_id from Yaps where user_id = $1
`

func (q *Queries) GetYapsByUserId(ctx context.Context, userID uuid.UUID) ([]Yap, error) {
	rows, err := q.db.QueryContext(ctx, getYapsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Yap
	for rows.Next() {
		var i Yap
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newYap = `-- name: NewYap :one
INSERT INTO Yaps(id, created_at, updated_at, user_id ,body )

VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, body, user_id
`

type NewYapParams struct {
	UserID uuid.UUID
	Body   string
}

func (q *Queries) NewYap(ctx context.Context, arg NewYapParams) (Yap, error) {
	row := q.db.QueryRowContext(ctx, newYap, arg.UserID, arg.Body)
	var i Yap
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const resetYaps = `-- name: ResetYaps :exec
DELETE from  Yaps
`

func (q *Queries) ResetYaps(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetYaps)
	return err
}

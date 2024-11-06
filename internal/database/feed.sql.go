// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getInitialFeed = `-- name: GetInitialFeed :many
select y.id, y.created_at, y.updated_at, y.body, y.user_id from yaps y join feed f on y.user_id = f.user_id
where y.user_id = $1 order by y.created_at Desc LIMIT 20
`

func (q *Queries) GetInitialFeed(ctx context.Context, userID uuid.UUID) ([]Yap, error) {
	rows, err := q.db.QueryContext(ctx, getInitialFeed, userID)
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

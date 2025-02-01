// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id, email, hashed_password, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5 ) RETURNING id, email, hashed_password, created_at, updated_at
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.HashedPassword,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllUser = `-- name: DeleteAllUser :exec
truncate table users CASCADE
`

func (q *Queries) DeleteAllUser(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUser)
	return err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, email, hashed_password, created_at, updated_at FROM users where lower(email) = lower($1)
`

func (q *Queries) FindUserByEmail(ctx context.Context, lower string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, lower)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT id, email, hashed_password, created_at, updated_at FROM users where id = $1
`

func (q *Queries) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserByID = `-- name: UpdateUserByID :exec
UPDATE users
    SET email = $1, hashed_password = $2, updated_at = $4
    WHERE id = $3
`

type UpdateUserByIDParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
	UpdatedAt      time.Time
}

func (q *Queries) UpdateUserByID(ctx context.Context, arg UpdateUserByIDParams) error {
	_, err := q.db.ExecContext(ctx, updateUserByID,
		arg.Email,
		arg.HashedPassword,
		arg.ID,
		arg.UpdatedAt,
	)
	return err
}

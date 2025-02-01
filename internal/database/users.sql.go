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
) VALUES ( $1, $2, $3, $4, $5 ) RETURNING id, email, hashed_password, created_at, updated_at, is_chirpy_red
`

type CreateUserParams struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
		&i.IsChirpyRed,
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
SELECT id, email, hashed_password, created_at, updated_at, is_chirpy_red FROM users where lower(email) = lower($1)
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
		&i.IsChirpyRed,
	)
	return i, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT id, email, hashed_password, created_at, updated_at, is_chirpy_red FROM users where id = $1
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
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserByID = `-- name: UpdateUserByID :exec
UPDATE users
    SET email = $1, hashed_password = $2, updated_at = $4
    WHERE id = $3
`

type UpdateUserByIDParams struct {
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	ID             uuid.UUID `json:"id"`
	UpdatedAt      time.Time `json:"updated_at"`
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

const updateUsertoChirpRed = `-- name: UpdateUsertoChirpRed :exec
UPDATE users SET is_chirpy_red = true WHERE id = $1
`

func (q *Queries) UpdateUsertoChirpRed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateUsertoChirpRed, id)
	return err
}

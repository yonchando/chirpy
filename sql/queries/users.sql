-- name: FindUserByEmail :one
SELECT * FROM users where lower(email) = lower($1);

-- name: FindUserByID :one
SELECT * FROM users where id = $1;

-- name: UpdateUserByID :exec
UPDATE users
    SET email = $1, hashed_password = $2, updated_at = $4
    WHERE id = $3;

-- name: UpdateUsertoChirpRed :exec
UPDATE users SET is_chirpy_red = true WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (
    id, email, hashed_password, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5 ) RETURNING *;

-- name: DeleteAllUser :exec
truncate table users CASCADE;

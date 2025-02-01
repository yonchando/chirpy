-- name: FindRefreshToken :one
SELECT * FROM refresh_tokens where token = $1 and revoked_at is null;


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
    SET revoked_at = $1, updated_at = $3
    WHERE token = $2;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token, user_id, expires_at, revoked_at, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6 ) RETURNING *;

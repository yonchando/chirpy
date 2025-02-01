-- name: GetAllChirps :many
SELECT * FROM chirps order by created_at;

-- name: FindChirpByID :one
SELECT * FROM chirps where id = $1;

-- name: CreateChirp :one
INSERT INTO chirps (
    id, body, user_id, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;

-- name: DeleteAllChirp :exec
truncate table chirps cascade;

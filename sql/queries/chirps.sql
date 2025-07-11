-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps;

-- name: ResetChirps :exec
SELECT FROM chirps;

-- name: GetChirpsByUser :many
SELECT * FROM chirps WHERE user_id = $1;

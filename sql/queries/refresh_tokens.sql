-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
	$1,
	NOW(),
	NOW(),
	$2,
	$3,
	NULL
)
RETURNING *;

-- name: ResetRefreshToken :exec
DELETE FROM refresh_tokens;

-- name: GetRefreshTokenByUser :one
SELECT token FROM refresh_tokens WHERE user_id = $1;

-- name: GetUserFromRefreshToken :one
SELECT user_id FROM refresh_tokens WHERE token = $1;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

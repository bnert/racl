-- name: CreateUser :one
INSERT INTO racl_users (
  name, api_key, api_secret
) VALUES (
  $1, $2, $3
) RETURNING *;


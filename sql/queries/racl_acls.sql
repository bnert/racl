-- name: CreateAcl :one
INSERT INTO racl_acls (
  resource_id, entity, capabilities
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: CreateDefaultAcl :one
INSERT INTO racl_acls (
  resource_id, entity, capabilities
) VALUES (
  $1, $2, '{"c", "r", "u", "d", "a"}'
) RETURNING *;

-- name: UpdateAclCapabilities :one
UPDATE racl_acls
SET capabilities = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAcl :one
DELETE FROM racl_acls
WHERE id = $1
RETURNING *;


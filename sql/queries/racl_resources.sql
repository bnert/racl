-- name: GetResource :one
SELECT *
FROM racl_resources
WHERE id = $1;

-- name: CreateResource :one
INSERT INTO racl_resources (
  id  
) VALUES (
  $1
) RETURNING *;

-- name: DeleteResource :one
DELETE FROM racl_resources
WHERE id = $1
RETURNING *;

-- name: GetAclForResourceByEntity :one
SELECT
  acl.capabilities
FROM
   racl_acls as acl
WHERE
  acl.entity = $1
  AND acl.resource_id = $2;


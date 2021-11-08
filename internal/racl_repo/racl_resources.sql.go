// Code generated by sqlc. DO NOT EDIT.
// source: racl_resources.sql

package racl_repo

import (
	"context"
)

const createResource = `-- name: CreateResource :one
INSERT INTO racl_resources (
  id  
) VALUES (
  $1
) RETURNING id, created_at, updated_at
`

func (q *Queries) CreateResource(ctx context.Context, id string) (RaclResource, error) {
	row := q.db.QueryRow(ctx, createResource, id)
	var i RaclResource
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteResource = `-- name: DeleteResource :one
DELETE FROM racl_resources
WHERE id = $1
RETURNING id, created_at, updated_at
`

func (q *Queries) DeleteResource(ctx context.Context, id string) (RaclResource, error) {
	row := q.db.QueryRow(ctx, deleteResource, id)
	var i RaclResource
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const getAclForResourceByEntity = `-- name: GetAclForResourceByEntity :one
SELECT
  id, created_at, updated_at, resource_id, entity, capabilities
FROM
  racl_acls as acl
WHERE
  acl.entity = $1
  AND acl.resource_id = $2
`

type GetAclForResourceByEntityParams struct {
	Entity     string `json:"entity"`
	ResourceID string `json:"resourceID"`
}

func (q *Queries) GetAclForResourceByEntity(ctx context.Context, arg GetAclForResourceByEntityParams) (RaclAcl, error) {
	row := q.db.QueryRow(ctx, getAclForResourceByEntity, arg.Entity, arg.ResourceID)
	var i RaclAcl
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ResourceID,
		&i.Entity,
		&i.Capabilities,
	)
	return i, err
}

const getResource = `-- name: GetResource :one
SELECT id, created_at, updated_at
FROM racl_resources
WHERE id = $1
`

func (q *Queries) GetResource(ctx context.Context, id string) (RaclResource, error) {
	row := q.db.QueryRow(ctx, getResource, id)
	var i RaclResource
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

-- name: CreateAcl :one
insert into racl_acls (
  resource_id, entity, capabilities
) values (
  $1, $2, $3
) returning *;

-- name: CreateDefaultAcl :one
insert into racl_acls (
  resource_id, entity, capabilities
) values (
  $1, $2, '{"c", "r", "u", "d", "a"}'
) returning *;

-- name: GetAclByEntity :one
select
  *
from
  racl_acls
where
  entity = $1;

-- name: GetAclByEntityAndResource :one
select
  *
from
  racl_acls
where
  entity = $1
  and resource_id = $2;

-- name: UpdateAclCapabilities :one
update
  racl_acls
set
  capabilities = $3
where
  entity = $1
  and resource_id = $2
returning *;

-- name: DeleteAcl :one
delete from racl_acls
where
  entity = $1
  and resource_id = $2
returning *;


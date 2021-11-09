-- name: GetResource :one
select *
from racl_resources
where id = $1;

-- name: CreateResource :one
insert into racl_resources (
  id  
) values (
  $1
) returning *;

-- name: DeleteResource :one
delete from racl_resources
where id = $1
returning *;


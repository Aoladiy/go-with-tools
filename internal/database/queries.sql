-- name: GetAllBrands :many
select *
from brands
where deleted_at is null;

-- name: GetBrand :one
select *
from brands
where id = $1
  and deleted_at is null
limit 1;

-- name: CreateBrand :one
insert into brands (name, slug)
VALUES ($1, $2)
returning *;

-- name: UpdateBrand :one
update brands
SET name       = $2,
    slug       = $3,
    updated_at = now()
where id = $1
  and deleted_at is null
returning *;

-- name: DeleteBrand :execrows
update brands
set deleted_at = now(),
    updated_at = now()
where id = $1
  and deleted_at is null;

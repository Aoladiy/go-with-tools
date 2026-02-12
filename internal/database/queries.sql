-- name: GetAllBrands :many
select id, name, slug, created_at, updated_at
from brands
where deleted_at is null;

-- name: GetBrand :one
select id, name, slug, created_at, updated_at
from brands
where id = $1
  and deleted_at is null
limit 1;

-- name: CreateBrand :one
insert into brands (name, slug)
VALUES ($1, $2)
returning id, name, slug, created_at, updated_at;

-- name: UpdateBrand :one
update brands
SET name       = $2,
    slug       = $3,
    updated_at = now()
where id = $1
  and deleted_at is null
returning id, name, slug, created_at, updated_at;

-- name: DeleteBrand :execrows
update brands
set deleted_at = now(),
    updated_at = now()
where id = $1
  and deleted_at is null;

-- name: GetAllCategories :many
select id, name, slug, parent_id, created_at, updated_at
from categories
where deleted_at is null;

-- name: GetCategory :one
select id, name, slug, parent_id, created_at, updated_at
from categories
where id = $1
  and deleted_at is null
limit 1;

-- name: CreateCategory :one
insert into categories (name, slug, parent_id)
VALUES ($1, $2, $3)
returning id, name, slug, parent_id, created_at, updated_at;

-- name: UpdateCategory :one
update categories
SET name       = $2,
    slug       = $3,
    parent_id  = $4,
    updated_at = now()
where id = $1
  and deleted_at is null
returning id, name, slug, parent_id, created_at, updated_at;

-- name: DeleteCategory :execrows
update categories
set deleted_at = now(),
    updated_at = now()
where id = $1
  and deleted_at is null;

-- name: GetAllProducts :many
select id,
       brand_id,
       category_id,
       name,
       slug,
       description,
       price_kopeck,
       is_active,
       created_at,
       updated_at
from products
where deleted_at is null;

-- name: GetProduct :one
select id,
       brand_id,
       category_id,
       name,
       slug,
       description,
       price_kopeck,
       is_active,
       created_at,
       updated_at
from products
where id = $1
  and deleted_at is null
limit 1;

-- name: CreateProduct :one
insert into products (brand_id, category_id, name, slug, description, price_kopeck, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
returning id, brand_id, category_id, name, slug, description, price_kopeck, is_active, created_at, updated_at;

-- name: UpdateProduct :one
update products
SET brand_id     = $2,
    category_id  = $3,
    name         = $4,
    slug         = $5,
    description  = $6,
    price_kopeck = $7,
    is_active    = $8,
    updated_at   = now()
where id = $1
  and deleted_at is null
returning id, brand_id, category_id, name, slug, description, price_kopeck, is_active, created_at, updated_at;

-- name: DeleteProduct :execrows
update products
set deleted_at = now(),
    updated_at = now()
where id = $1
  and deleted_at is null;
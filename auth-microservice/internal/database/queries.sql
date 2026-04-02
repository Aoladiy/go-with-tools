-- name: CreateAdminUser :one
insert into admin_users (email, password_hash)
VALUES ($1, $2)
returning id, email, created_at, updated_at;

-- name: GetAdminUser :one
select id, email, password_hash, created_at, updated_at
from admin_users
where email = $1
  and deleted_at is null;

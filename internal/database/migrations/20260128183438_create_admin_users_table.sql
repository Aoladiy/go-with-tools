-- +goose Up
-- +goose StatementBegin

create table admin_users
(
    id            bigint generated always as identity primary key,
    email         text        not null unique,
    password_hash text        not null,
    created_at    timestamptz not null default now(),
    updated_at    timestamptz not null default now(),
    deleted_at    timestamptz          default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists admin_users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
create table brands
(
    id         bigint generated always as identity primary key,
    name       text        not null unique,
    slug       text        not null unique,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz          default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists brands;
-- +goose StatementEnd

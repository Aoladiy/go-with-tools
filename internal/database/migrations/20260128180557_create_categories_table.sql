-- +goose Up
-- +goose StatementBegin
create table categories
(
    id         bigint generated always as identity primary key,
    name       text        not null,
    slug       text        not null unique,
    parent_id  bigint,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz          default null,

    constraint fk_categories_parent_id
        foreign key (parent_id)
            references categories (id)
            on delete SET NULL
);
create index idx_categories_parent_id on categories (parent_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop index if exists idx_categories_parent_id;
drop table if exists categories;
-- +goose StatementEnd

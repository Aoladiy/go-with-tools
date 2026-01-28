-- +goose Up
-- +goose StatementBegin

create table products
(
    id           bigint generated always as identity primary key,
    brand_id     bigint      not null,
    category_id  bigint      not null,
    name         text        not null,
    slug         text        not null unique,
    description  text        not null default '',
    price_kopeck int         not null check ( price_kopeck >= 0 ),
    is_active    boolean     not null default true,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now(),
    deleted_at   timestamptz          default null,

    constraint fk_products_brand_id
        foreign key (brand_id)
            references brands (id)
            on delete restrict,
    constraint fk_products_category_id
        foreign key (category_id)
            references categories (id)
            on delete restrict
);
create index idx_products_is_active on products (is_active);
create index idx_products_brand_id on products (brand_id);
create index idx_products_category_id on products (category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_products_is_active;
drop index if exists idx_products_brand_id;
drop index if exists idx_products_category_id;
drop table if exists products;
-- +goose StatementEnd

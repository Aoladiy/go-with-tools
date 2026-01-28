-- +goose Up
-- +goose StatementBegin

create table product_price_history
(
    id               bigint generated always as identity primary key,
    product_id       bigint      not null,
    old_price_kopeck int         not null check ( old_price_kopeck >= 0 ),
    new_price_kopeck int         not null check ( new_price_kopeck >= 0 ),
    created_at       timestamptz not null default now(),
    updated_by       bigint      not null,

    constraint fk_product_price_history_product_id
        foreign key (product_id)
            references products (id)
            on delete restrict,

    constraint fk_product_price_history_updated_by
        foreign key (updated_by)
            references admin_users (id)
            on delete restrict
);
create index idx_product_price_history_product_id on product_price_history (product_id);
create index idx_product_price_history_updated_by on product_price_history (updated_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_product_price_history_product_id;
drop index if exists idx_product_price_history_updated_by;
drop table if exists product_price_history;
-- +goose StatementEnd

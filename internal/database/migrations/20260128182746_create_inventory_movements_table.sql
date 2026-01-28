-- +goose Up
-- +goose StatementBegin

create table inventory_movements
(
    id          bigint generated always as identity primary key,
    product_id  bigint      not null,
    delta       int         not null check ( delta <> 0 ),
    description text        not null,
    created_at  timestamptz not null default now(),

    constraint fk_inventory_movements_product_id
        foreign key (product_id)
            references products (id)
            on delete restrict
);
create index idx_inventory_movements_product_id on inventory_movements (product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_inventory_movements_product_id;
drop table if exists inventory_movements;
-- +goose StatementEnd

-- +goose Up
-- Make unique constraints tenant-scoped (multi-tenancy)
-- +goose StatementBegin
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_order_number_key;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_tenant_number ON orders (tenant_id, order_number);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_code_key;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS products_code_key;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_products_code;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_tenant_number;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE orders ADD CONSTRAINT orders_order_number_key UNIQUE (order_number);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS products_code_key ON products (code);
-- +goose StatementEnd

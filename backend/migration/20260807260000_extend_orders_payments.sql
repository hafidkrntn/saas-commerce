-- +goose Up
-- Extends orders/order_items (created earlier) with the full checkout model
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS customer_id UUID REFERENCES customers (id),
    ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS payment_method_id UUID,
    ADD COLUMN IF NOT EXISTS shipping_method VARCHAR(100),
    ADD COLUMN IF NOT EXISTS shipping_address JSONB,
    ADD COLUMN IF NOT EXISTS subtotal NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS shipping_cost NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tax NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE order_items
    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS product_name VARCHAR(200) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS product_sku VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS product_image VARCHAR(500);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_orders_tenant_id ON orders (tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders (payment_status);
CREATE INDEX IF NOT EXISTS idx_order_items_tenant_id ON order_items (tenant_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_timeline (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    user_name VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_order_timeline_order_id ON order_timeline (order_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payment_methods_tenant_id ON payment_methods (tenant_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    payment_method_id UUID REFERENCES payment_methods (id),
    provider VARCHAR(50) NOT NULL DEFAULT 'mock',
    provider_transaction_id VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payments_tenant_id ON payments (tenant_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_provider_txn ON payments (provider_transaction_id);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO payment_methods (id, tenant_id, name, type, is_enabled, description, config) VALUES
('00000000-0000-0000-0000-000000000401', '00000000-0000-0000-0000-000000000101', 'Bank Transfer (BCA)', 'bank_transfer', TRUE, 'Transfer manual melalui BCA', '{"bank":"BCA","account":"1234567890"}'::jsonb),
('00000000-0000-0000-0000-000000000402', '00000000-0000-0000-0000-000000000101', 'Kartu Kredit', 'card', TRUE, 'Pembayaran dengan kartu kredit', '{"provider":"mock"}'::jsonb),
('00000000-0000-0000-0000-000000000403', '00000000-0000-0000-0000-000000000101', 'COD', 'cod', TRUE, 'Bayar di tempat', '{}'::jsonb),
('00000000-0000-0000-0000-000000000404', '00000000-0000-0000-0000-000000000101', 'GoPay', 'e_wallet', FALSE, 'Dompet digital GoPay', '{"provider":"mock"}'::jsonb)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS payment_methods;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS order_timeline;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE order_items
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS product_name,
    DROP COLUMN IF EXISTS product_sku,
    DROP COLUMN IF EXISTS product_image;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS customer_id,
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS payment_method_id,
    DROP COLUMN IF EXISTS shipping_method,
    DROP COLUMN IF EXISTS shipping_address,
    DROP COLUMN IF EXISTS subtotal,
    DROP COLUMN IF EXISTS shipping_cost,
    DROP COLUMN IF EXISTS tax,
    DROP COLUMN IF EXISTS discount,
    DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd

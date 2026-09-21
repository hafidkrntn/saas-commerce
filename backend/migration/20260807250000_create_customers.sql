-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200),
    phone VARCHAR(50),
    avatar VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    total_orders INTEGER NOT NULL DEFAULT 0,
    total_spent NUMERIC(15, 2) NOT NULL DEFAULT 0,
    lifetime_value NUMERIC(15, 2) NOT NULL DEFAULT 0,
    average_order_value NUMERIC(15, 2) NOT NULL DEFAULT 0,
    last_order_at TIMESTAMPTZ,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_customers_tenant_id ON customers (tenant_id);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers (email);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers (phone);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customer_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    label VARCHAR(50) NOT NULL DEFAULT 'home',
    line1 VARCHAR(255) NOT NULL,
    line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    zip VARCHAR(20),
    country VARCHAR(100) NOT NULL DEFAULT 'Indonesia',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_customer_addresses_customer_id ON customer_addresses (customer_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customer_activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_customer_activities_customer_id ON customer_activities (customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS customer_activities;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS customer_addresses;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS customers;
-- +goose StatementEnd

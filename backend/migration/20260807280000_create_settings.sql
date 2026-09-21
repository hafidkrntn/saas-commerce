-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS store_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL DEFAULT '',
    description TEXT,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',
    order_prefix VARCHAR(20) NOT NULL DEFAULT 'ORD',
    low_stock_threshold INTEGER NOT NULL DEFAULT 5,
    enable_reviews BOOLEAN NOT NULL DEFAULT TRUE,
    enable_wishlist BOOLEAN NOT NULL DEFAULT TRUE,
    enable_gift_cards BOOLEAN NOT NULL DEFAULT FALSE,
    default_weight_unit VARCHAR(10) NOT NULL DEFAULT 'kg',
    default_dimension_unit VARCHAR(10) NOT NULL DEFAULT 'cm',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS company_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL DEFAULT '',
    legal_name VARCHAR(200),
    logo VARCHAR(500),
    email VARCHAR(200),
    phone VARCHAR(50),
    website VARCHAR(200),
    address JSONB NOT NULL DEFAULT '{}'::jsonb,
    tax_id VARCHAR(50),
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',
    date_format VARCHAR(20) NOT NULL DEFAULT 'd MMM yyyy',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tax_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    rate NUMERIC(5, 2) NOT NULL DEFAULT 0,
    region VARCHAR(100) NOT NULL DEFAULT 'all',
    type VARCHAR(20) NOT NULL DEFAULT 'vat',
    applies_to VARCHAR(20) NOT NULL DEFAULT 'all',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_tax_rates_tenant_id ON tax_rates (tenant_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipping_zones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    regions JSONB NOT NULL DEFAULT '[]'::jsonb,
    method VARCHAR(20) NOT NULL DEFAULT 'flat_rate',
    rate NUMERIC(15, 2) NOT NULL DEFAULT 0,
    free_above NUMERIC(15, 2),
    estimated_days VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_shipping_zones_tenant_id ON shipping_zones (tenant_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    event VARCHAR(100) NOT NULL,
    email BOOLEAN NOT NULL DEFAULT FALSE,
    sms BOOLEAN NOT NULL DEFAULT FALSE,
    in_app BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, event)
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO store_settings (tenant_id, name, description, currency, timezone, order_prefix) VALUES
('00000000-0000-0000-0000-000000000101', 'LAKU', 'Toko demo LAKU', 'IDR', 'Asia/Jakarta', 'ORD')
ON CONFLICT (tenant_id) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO company_settings (tenant_id, name, legal_name, email, phone, website, currency, timezone) VALUES
('00000000-0000-0000-0000-000000000101', 'LAKU', 'PT Laku Digital Nusantara', 'halo@laku.com', '+62 21 1234 5678', 'https://laku.com', 'IDR', 'Asia/Jakarta')
ON CONFLICT (tenant_id) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO notification_settings (tenant_id, event, email, sms, in_app, description) VALUES
('00000000-0000-0000-0000-000000000101', 'order_created', TRUE, FALSE, TRUE, 'Notifikasi saat order baru dibuat'),
('00000000-0000-0000-0000-000000000101', 'order_shipped', TRUE, TRUE, TRUE, 'Notifikasi saat order dikirim'),
('00000000-0000-0000-0000-000000000101', 'order_delivered', TRUE, FALSE, TRUE, 'Notifikasi saat order selesai'),
('00000000-0000-0000-0000-000000000101', 'order_cancelled', TRUE, FALSE, TRUE, 'Notifikasi saat order dibatalkan'),
('00000000-0000-0000-0000-000000000101', 'payment_received', TRUE, FALSE, TRUE, 'Notifikasi saat pembayaran diterima'),
('00000000-0000-0000-0000-000000000101', 'low_stock', TRUE, FALSE, TRUE, 'Notifikasi saat stok menipis')
ON CONFLICT (tenant_id, event) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_settings;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS shipping_zones;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS tax_rates;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS company_settings;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS store_settings;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    price NUMERIC(15, 2) NOT NULL DEFAULT 0,
    description TEXT,
    features JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    logo VARCHAR(500),
    plan_id UUID REFERENCES plans (id),
    subscription_status VARCHAR(20) NOT NULL DEFAULT 'trial',
    trial_ends_at TIMESTAMPTZ,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants (slug);
CREATE INDEX IF NOT EXISTS idx_tenants_plan_id ON tenants (plan_id);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO plans (id, name, code, price, description, features) VALUES
('00000000-0000-0000-0000-000000000001', 'Free', 'free', 0, 'Untuk toko baru', '["1 staff","100 produk"]'::jsonb),
('00000000-0000-0000-0000-000000000002', 'Pro', 'pro', 199000, 'Untuk toko yang berkembang', '["10 staff","1000 produk","Analitik lanjutan"]'::jsonb),
('00000000-0000-0000-0000-000000000003', 'Business', 'business', 499000, 'Untuk bisnis besar', '["Unlimited staff","Produk unlimited","API access"]'::jsonb)
ON CONFLICT (code) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS plans;
-- +goose StatementEnd

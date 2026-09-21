-- +goose Up
-- Permission bits: 1=products 2=orders 4=customers 8=inventory 16=settings
--                  32=users 64=analytics 128=reviews 256=payments
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions BIGINT NOT NULL DEFAULT 0,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_roles_tenant_id ON roles (tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_tenant_name ON roles (tenant_id, name) WHERE tenant_id IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users (tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash VARCHAR(128) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens (token_hash);
-- +goose StatementEnd

-- +goose StatementBegin
-- Seed: demo tenant, system roles, and admin users
INSERT INTO tenants (id, name, slug, logo, plan_id, subscription_status) VALUES
('00000000-0000-0000-0000-000000000101', 'LAKU Demo Store', 'laku-demo', 'https://placehold.co/200x200?text=LAKU', '00000000-0000-0000-0000-000000000002', 'active')
ON CONFLICT (slug) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO roles (id, tenant_id, name, description, permissions, is_system) VALUES
('00000000-0000-0000-0000-000000000201', NULL, 'platform-admin', 'Platform administrator', 511, TRUE),
('00000000-0000-0000-0000-000000000202', '00000000-0000-0000-0000-000000000101', 'admin', 'Store administrator', 511, TRUE),
('00000000-0000-0000-0000-000000000203', '00000000-0000-0000-0000-000000000101', 'staff', 'Store staff', 79, TRUE)
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO users (id, tenant_id, name, email, password_hash, status) VALUES
('00000000-0000-0000-0000-000000000301', NULL, 'Platform Admin', 'platform@laku.id', '$2a$10$RVyD/NM4rCZOtqKyT4pFD.Si1L4AvFB5AAliiS34fAmdMMp3Meqz.', 'active'),
('00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000101', 'Store Admin', 'admin@laku.id', '$2a$10$RVyD/NM4rCZOtqKyT4pFD.Si1L4AvFB5AAliiS34fAmdMMp3Meqz.', 'active'),
('00000000-0000-0000-0000-000000000303', '00000000-0000-0000-0000-000000000101', 'Staff User', 'staff@laku.id', '$2a$10$RVyD/NM4rCZOtqKyT4pFD.Si1L4AvFB5AAliiS34fAmdMMp3Meqz.', 'active')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO user_roles (user_id, role_id) VALUES
('00000000-0000-0000-0000-000000000301', '00000000-0000-0000-0000-000000000201'),
('00000000-0000-0000-0000-000000000302', '00000000-0000-0000-0000-000000000202'),
('00000000-0000-0000-0000-000000000303', '00000000-0000-0000-0000-000000000203')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS user_roles;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd

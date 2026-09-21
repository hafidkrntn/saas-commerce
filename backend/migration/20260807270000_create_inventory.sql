-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    address TEXT,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_warehouses_tenant_id ON warehouses (tenant_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id),
    warehouse_id UUID REFERENCES warehouses (id),
    type VARCHAR(20) NOT NULL,
    quantity INTEGER NOT NULL,
    reference VARCHAR(100),
    note TEXT,
    user_name VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_stock_movements_tenant_id ON stock_movements (tenant_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements (product_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements (created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incoming_shipments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id),
    quantity INTEGER NOT NULL,
    expected_date DATE,
    supplier VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    notes TEXT,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_incoming_shipments_tenant_id ON incoming_shipments (tenant_id);
CREATE INDEX IF NOT EXISTS idx_incoming_shipments_product_id ON incoming_shipments (product_id);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO warehouses (id, tenant_id, name, address, is_default) VALUES
('00000000-0000-0000-0000-000000000501', '00000000-0000-0000-0000-000000000101', 'Gudang Utama', 'Jakarta Selatan, DKI Jakarta', TRUE),
('00000000-0000-0000-0000-000000000502', '00000000-0000-0000-0000-000000000101', 'Gudang Surabaya', 'Surabaya, Jawa Timur', FALSE)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incoming_shipments;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS stock_movements;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS warehouses;
-- +goose StatementEnd

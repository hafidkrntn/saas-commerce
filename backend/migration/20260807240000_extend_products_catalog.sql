-- +goose Up
-- Extends products (created earlier) with the full SaaS catalog fields
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS sku VARCHAR(100),
    ADD COLUMN IF NOT EXISTS compare_at_price NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cost_price NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reserved INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS incoming INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS low_stock_threshold INTEGER NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS rating NUMERIC(3, 2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS review_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_products_tenant_id ON products (tenant_id);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_tenant_sku ON products (tenant_id, sku) WHERE sku IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_tenant_name ON products (tenant_id, name);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt VARCHAR(200),
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images (product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    customer_id UUID,
    customer_name VARCHAR(200) NOT NULL,
    customer_avatar VARCHAR(500),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title VARCHAR(200),
    content TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_reviews_tenant_id ON reviews (tenant_id);
CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews (product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviews;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS product_images;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE products
    DROP COLUMN IF EXISTS tenant_id,
    DROP COLUMN IF EXISTS sku,
    DROP COLUMN IF EXISTS compare_at_price,
    DROP COLUMN IF EXISTS cost_price,
    DROP COLUMN IF EXISTS reserved,
    DROP COLUMN IF EXISTS incoming,
    DROP COLUMN IF EXISTS low_stock_threshold,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS review_count,
    DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd

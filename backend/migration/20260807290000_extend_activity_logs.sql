-- +goose Up
-- +goose StatementBegin
ALTER TABLE activity_logs
    ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES tenants (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_activity_logs_tenant_id ON activity_logs (tenant_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activity_logs
    DROP COLUMN IF EXISTS tenant_id;
-- +goose StatementEnd

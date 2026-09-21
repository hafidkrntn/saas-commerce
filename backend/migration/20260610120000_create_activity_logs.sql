-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    method VARCHAR(10) NOT NULL,
    module VARCHAR(100),
    endpoint VARCHAR(255) NOT NULL,
    record_id UUID,
    payload JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_module ON activity_logs(module);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activity_logs;
-- +goose StatementEnd

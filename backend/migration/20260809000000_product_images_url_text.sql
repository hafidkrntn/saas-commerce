-- +goose Up
-- Data URL gambar (base64) jauh melebihi VARCHAR(500)
-- +goose StatementBegin
ALTER TABLE product_images ALTER COLUMN url TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_images ALTER COLUMN url TYPE VARCHAR(500);
-- +goose StatementEnd

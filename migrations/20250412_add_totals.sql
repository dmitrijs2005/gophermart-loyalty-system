-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN accrued_total NUMERIC(15, 2) DEFAULT 0;
ALTER TABLE users ADD COLUMN withdrawn_total NUMERIC(15, 2) DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN accrued_total;
ALTER TABLE users DROP COLUMN withdrawn_total;
-- +goose StatementEnd


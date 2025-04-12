-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdrawals ALTER COLUMN amount TYPE INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdrawals ALTER COLUMN amount TYPE NUMERIC(15, 2);
-- +goose StatementEnd


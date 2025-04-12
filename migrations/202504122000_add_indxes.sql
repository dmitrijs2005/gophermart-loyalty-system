-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_withdrawals_user_id ON withdrawals (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_orders_user_id;
DROP INDEX idx_withdrawals_user_id;
-- +goose StatementEnd


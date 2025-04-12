-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(), 
    login TEXT NOT NULL,
    password TEXT NOT NULL,  -- e.g., "int64" or "float64"

    PRIMARY KEY (id)  -- PK
);

CREATE UNIQUE INDEX unique_user_login ON users (login);

CREATE TABLE orders (
    id uuid DEFAULT gen_random_uuid(), 
    user_id uuid NOT NULL,
    number TEXT NOT NULL,
    status TEXT NOT NULL,
    accrual NUMERIC(15, 2) DEFAULT 0,
    uploaded_at TIMESTAMPTZ DEFAULT now(),

    PRIMARY KEY (id)  -- PK
);

CREATE UNIQUE INDEX unique_order_number ON orders (number);

CREATE TABLE withdrawals (
    id uuid DEFAULT gen_random_uuid(), 
    user_id uuid NOT NULL,
    "order" TEXT NOT NULL,
    amount NUMERIC(15, 2) DEFAULT 0,
    uploaded_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (id)  -- PK
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE orders;
DROP TABLE withdrawals;
-- +goose StatementEnd

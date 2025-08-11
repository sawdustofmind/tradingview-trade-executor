CREATE TABLE customers
(
    id                    SERIAL PRIMARY KEY,
    username              TEXT                     NOT NULL UNIQUE,
    password              TEXT                     NOT NULL,
    legal_name            TEXT                     NOT NULL,
    gender                TEXT                     NOT NULL,
    country               TEXT                     NOT NULL,
    phone_number          TEXT                     NOT NULL,
    image_base64          TEXT                     NOT NULL,
    bybit_api_key         TEXT,
    bybit_test_api_key    TEXT,
    bybit_api_secret      TEXT,
    bybit_test_api_secret TEXT,
    created_at            timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at            timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON customers
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
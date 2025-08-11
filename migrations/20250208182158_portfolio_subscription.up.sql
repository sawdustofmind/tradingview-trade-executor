CREATE TABLE portfolio_subscription
(
    id           SERIAL PRIMARY KEY,
    portfolio_id INTEGER                  NOT NULL,
    customer_id  INTEGER                  NOT NULL,
    is_test      BOOLEAN                  NOT NULL DEFAULT true,
    amount       DECIMAL(18, 8)           NOT NULL,
    exchange     TEXT                     NOT NULL,
    status       TEXT                     NOT NULL,
    pnl          DECIMAL(18, 8)           NOT NULL,
    created_at   timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at   timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON portfolio_subscription
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
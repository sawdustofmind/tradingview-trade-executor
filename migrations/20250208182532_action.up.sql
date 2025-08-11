CREATE TABLE action
(
    id                   SERIAL PRIMARY KEY,
    corr_id              UUID                     NOT NULL,
    customer_id          INT                      NOT NULL,
    sub_id               INT                      NOT NULL,
    portfolio_id         INT                      NOT NULL,
    action_type          TEXT                     NOT NULL,
    details              jsonb                    NOT NULL,
    need_to_fetch_trades BOOLEAN                  NOT NULL,
    error                TEXT                     NOT NULL,
    created_at           timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE action_trades
(
    id           SERIAL PRIMARY KEY,
    corr_id      UUID                     NOT NULL,
    customer_id  INT                      NOT NULL,
    sub_id       INT                      NOT NULL,
    portfolio_id INT                      NOT NULL,
    exchange     TEXT                     NOT NULL,
    side         TEXT                     NOT NULL,
    symbol       TEXT                     NOT NULL,
    quantity     DECIMAL(18, 8)           NOT NULL,
    price        DECIMAL(18, 8)           NOT NULL,
    commission   DECIMAL(18, 8)           NOT NULL,
    created_at   timestamp with time zone NOT NULL DEFAULT NOW()
);

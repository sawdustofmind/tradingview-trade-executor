CREATE TABLE fren
(
    id           SERIAL PRIMARY KEY,
    name         TEXT                     NOT NULL,
    description  TEXT                     NOT NULL,
    image_base64 TEXT                     NOT NULL,
    created_at   timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at   timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE fren_portfolio
(
    fren_id      INT NOT NULL REFERENCES fren (id) ON DELETE cascade,
    portfolio_id INT NOT NULL REFERENCES portfolio (id) ON DELETE cascade,
    PRIMARY KEY (fren_id, portfolio_id)
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON fren
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
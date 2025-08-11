CREATE TABLE portfolio
(
    id                       SERIAL PRIMARY KEY,
    name                     TEXT                     NOT NULL,
    description              TEXT                     NOT NULL,
    image_base64             TEXT                     NOT NULL,
    year_pnl                 DECIMAL(18, 2)           NOT NULL,
    avg_delay                TEXT                     NOT NULL,
    risk_level               TEXT                     NOT NULL,
    strategy_type            TEXT                     NOT NULL,
    dca_levels               INT                      NOT NULL,
    leverage                 INT                      NOT NULL,
    cycle_investment_percent DECIMAL(18, 2)           NOT NULL,
    holdings                 jsonb                    NOT NULL,
    created_at               timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at               timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON portfolio
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

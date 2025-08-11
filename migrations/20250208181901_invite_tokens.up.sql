CREATE TABLE invite_tokens
(
    id         SERIAL PRIMARY KEY,
    token      BYTEA UNIQUE             NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);
-- +goose Up
CREATE TABLE todos(
    id                  BIGSERIAL       PRIMARY KEY,
    title               VARCHAR(20)     NOT NULL,
    description         TEXT            NOT NULL,
    completed           BOOLEAN         NOT NULL        DEFAULT FALSE,
    completed_at        TIMESTAMPTZ     NULL,
    created_at          TIMESTAMPTZ     NOT NULL        DEFAULT clock_timestamp(),
    updated_at          TIMESTAMPTZ     NOT NULL        DEFAULT clock_timestamp()
);

-- +goose Down
DROP TABLE IF EXISTS todos;

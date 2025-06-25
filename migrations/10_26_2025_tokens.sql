-- +goose up
-- +goose statementbegin

CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL references users(id) ON DELETE CASCADE,
    expired TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);
-- +goose statementend

-- +goose down
-- +goose statementbegin
DROP TABLE tokens;
-- +goose statementend
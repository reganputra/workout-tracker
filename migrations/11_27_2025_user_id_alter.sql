-- +goose up
-- +goose statementbegin

ALTER TABLE workout ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose statementend


-- +goose down
-- +goose statementbegin
ALTER TABLE workout DROP COLUMN user_id;
-- +goose statementend
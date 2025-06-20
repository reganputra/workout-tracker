-- +goose up
-- +goose statementbegin
CREATE table if not exists users(
                                    id bigserial primary key,
                                    username varchar(255) unique not null,
    email varchar(255) unique not null,
    password_hash varchar(255) not null,
    bio text,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp  -- No comma here
                             )
-- +goose statementend


-- +goose down
-- +goose statementbegin
DROP TABLE users;
-- +goose statementend
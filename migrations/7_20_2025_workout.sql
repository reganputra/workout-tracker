-- +goose up
-- +goose statementbegin
CREATE table if not exists workout(
    id bigserial primary key,
    title varchar(255) not null,
    description text,
    duration integer not null,
    calories_burned integer not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
 )
-- +goose statementend


-- +goose down
-- +goose statementbegin
DROP TABLE workout;
-- +goose statementend
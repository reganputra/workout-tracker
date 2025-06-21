-- +goose up
-- +goose statementbegin
CREATE table if not exists workout_entries(
                                              id bigserial primary key,
                                              workout_id bigint not null references workout(id) on delete cascade,
    exercise_name varchar(255) not null,
    sets integer not null,
    reps integer,
    duration_seconds integer,
    weight decimal(5,2),
    notes text,
    order_index integer not null,
    created_at timestamp with time zone default current_timestamp,
                                                                                                    constraint valid_workout_entry check (
                                                                                                    (reps is not null or duration_seconds is not null) AND
(reps is null or duration_seconds is null)
    )
    )
-- +goose statementend


-- +goose down
-- +goose statementbegin
DROP TABLE workout_entries;  -- Changed from "users" to "workout_entries"
-- +goose statementend

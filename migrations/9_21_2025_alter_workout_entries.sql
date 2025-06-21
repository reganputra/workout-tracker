-- +goose up
-- +goose statementbegin
ALTER TABLE workout_entries
    ALTER COLUMN reps DROP NOT NULL,
    ALTER COLUMN duration_seconds DROP NOT NULL,
    ALTER COLUMN weight DROP NOT NULL;
-- +goose statementend

-- +goose down
-- +goose statementbegin
ALTER TABLE workout_entries
    ALTER COLUMN reps SET NOT NULL,
    ALTER COLUMN duration_seconds SET NOT NULL,
    ALTER COLUMN weight SET NOT NULL;
-- +goose statementend
package store

import "database/sql"

type WorkoutEntry struct {
	Id              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type Workout struct {
	Id              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{
		db: db,
	}
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
}

func (ws *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := ws.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := "INSERT INTO workout (title, description, duration, calories_burned) VALUES ($1, $2, $3, $4) RETURNING id"

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.Id)
	if err != nil {
		return nil, err
	}
	for _, entry := range workout.Entries {
		query := "INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
		err = tx.QueryRow(query, workout.Id, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.Id)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (ws *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	query := "SELECT id, title, description, duration, calories_burned FROM workout WHERE id = $1"
	workout := &Workout{}
	err := ws.db.QueryRow(query, id).Scan(&workout.Id, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err != nil {
		return nil, err
	}

	query = "SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index FROM workout_entries WHERE workout_id = $1 ORDER BY order_index"
	rows, err := ws.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entry := WorkoutEntry{}
		err = rows.Scan(&entry.Id, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

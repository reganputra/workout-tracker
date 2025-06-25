package testing

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"workout-tracker/store"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	// run migrations for testing
	err = store.Migrate(db, "../migrations")
	if err != nil {
		t.Fatalf("Failed to run test db migrations: %v", err)
	}

	_, err = db.Exec("TRUNCATE workout, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	init := store.NewWorkoutStore(db)

	tests := []struct {
		name    string
		workout *store.Workout
		wantErr bool
	}{
		{
			name: "Create valid workout without entries",
			workout: &store.Workout{
				Title:           "Morning Workout",
				Description:     "A quick morning workout routine",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries:         []store.WorkoutEntry{},
			},
			wantErr: false,
		},
		{
			name: "Create valid workout with entries",
			workout: &store.Workout{
				Title:           "Full Body Workout",
				Description:     "Complete full body routine",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []store.WorkoutEntry{
					{
						ExerciseName: "Push-ups",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       Float64Ptr(135.5),
						Notes:        "Good for full body",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		// can add more test cases here, such as invalid workouts
		{
			name: "Create invalid workout",
			workout: &store.Workout{
				Title:           "Push up day",
				Description:     "Complete upper body day",
				DurationMinutes: 90,
				CaloriesBurned:  300,
				Entries: []store.WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         2,
						Reps:         IntPtr(60),
						Notes:        "Keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            4,
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(60),
						Weight:          Float64Ptr(140.5),
						Notes:           "Focus on depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true, // This should fail because the duration is too long for a workout with only 2 entries
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := init.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, got.Title)
			assert.Equal(t, tt.workout.Description, got.Description)
			assert.Equal(t, tt.workout.DurationMinutes, got.DurationMinutes)

			retrieved, err := init.GetWorkoutById(int64(got.Id))
			require.NoError(t, err)

			assert.Equal(t, tt.workout.Title, retrieved.Title)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

package set

import (
	"testing"
	"time"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

func TestNewSet(t *testing.T) {
	t.Parallel()
	id, err := NewSetIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new set id: %v", err)
	}
	userID, err := user.NewUserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	workoutID, err := workout.NewWorkoutIDFromString("a1a1a1a1-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new workout id: %v", err)
	}
	exerciseID, err := exercise.NewExerciseIDFromString("f1f538e5-4a37-409c-be99-09ee7bfefc50")
	if err != nil {
		t.Errorf("failed to new exercise id: %v", err)
	}
	rep, err := NewRep(10)
	if err != nil {
		t.Errorf("failed to new rep: %v", err)
	}
	weight, err := NewWeight(60.0)
	if err != nil {
		t.Errorf("failed to new weight: %v", err)
	}
	trainedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		success    bool
		id         SetID
		userID     user.UserID
		workoutID  workout.WorkoutID
		exerciseID exercise.ExerciseID
		rep        Rep
		weight     Weight
		trainedAt  time.Time
		createdAt  time.Time
	}{
		{"success new set", true, id, userID, workoutID, exerciseID, rep, weight, trainedAt, time.Now()},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewSet(tt.id, tt.userID, tt.workoutID, tt.exerciseID, tt.rep, tt.weight, tt.trainedAt, tt.createdAt)
			if tt.success && s.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", s.ID(), tt.id)
			}
			if tt.success && s.UserID() != tt.userID {
				t.Errorf("UserID() = %v, want %v", s.UserID(), tt.userID)
			}
			if tt.success && s.WorkoutID() != tt.workoutID {
				t.Errorf("WorkoutID() = %v, want %v", s.WorkoutID(), tt.workoutID)
			}
			if tt.success && s.ExerciseID() != tt.exerciseID {
				t.Errorf("ExerciseID() = %v, want %v", s.ExerciseID(), tt.exerciseID)
			}
			if tt.success && s.Rep() != tt.rep {
				t.Errorf("Rep() = %v, want %v", s.Rep(), tt.rep)
			}
			if tt.success && s.Weight() != tt.weight {
				t.Errorf("Weight() = %v, want %v", s.Weight(), tt.weight)
			}
			if tt.success && !s.TrainedAt().Equal(tt.trainedAt) {
				t.Errorf("TrainedAt() = %v, want %v", s.TrainedAt(), tt.trainedAt)
			}
			if tt.success && !s.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", s.CreatedAt(), tt.createdAt)
			}
		})
	}
}

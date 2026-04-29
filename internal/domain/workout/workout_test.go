package workout

import (
	"testing"
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
)

func TestNewWorkout(t *testing.T) {
	t.Parallel()
	id, err := NewWorkoutIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new workout id: %v", err)
	}
	userID, err := user.NewUserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	startedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	finishedAt := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		success    bool
		id         WorkoutID
		userID     user.UserID
		startedAt  time.Time
		finishedAt *time.Time
		createdAt  time.Time
		isFinished bool
	}{
		{"success in-progress workout", true, id, userID, startedAt, nil, time.Now(), false},
		{"success finished workout", true, id, userID, startedAt, &finishedAt, time.Now(), true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := NewWorkout(tt.id, tt.userID, tt.startedAt, tt.finishedAt, tt.createdAt)
			if tt.success && w.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", w.ID(), tt.id)
			}
			if tt.success && w.UserID() != tt.userID {
				t.Errorf("UserID() = %v, want %v", w.UserID(), tt.userID)
			}
			if tt.success && !w.StartedAt().Equal(tt.startedAt) {
				t.Errorf("StartedAt() = %v, want %v", w.StartedAt(), tt.startedAt)
			}
			if tt.success && w.IsFinished() != tt.isFinished {
				t.Errorf("IsFinished() = %v, want %v", w.IsFinished(), tt.isFinished)
			}
			if tt.isFinished {
				got := w.FinishedAt()
				if got == nil || !got.Equal(*tt.finishedAt) {
					t.Errorf("FinishedAt() = %v, want %v", got, tt.finishedAt)
				}
			} else if w.FinishedAt() != nil {
				t.Errorf("FinishedAt() should be nil, got %v", w.FinishedAt())
			}
			if tt.success && !w.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", w.CreatedAt(), tt.createdAt)
			}
		})
	}
}

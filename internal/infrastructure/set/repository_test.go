package set

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
	"github.com/qkitzero/workout-service/testutil"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, s set.Set)
	}{
		{
			name:    "success create set",
			success: true,
			setup: func(mock sqlmock.Sqlmock, s set.Set) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "sets" ("id","user_id","workout_id","exercise_id","rep","weight","trained_at","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
					WithArgs(s.ID(), s.UserID(), s.WorkoutID(), s.ExerciseID(), s.Rep(), s.Weight(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create set error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, s set.Set) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "sets" ("id","user_id","workout_id","exercise_id","rep","weight","trained_at","created_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`)).
					WithArgs(s.ID(), s.UserID(), s.WorkoutID(), s.ExerciseID(), s.Rep(), s.Weight(), testutil.AnyTime{}, testutil.AnyTime{}).
					WillReturnError(errors.New("create set error"))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSet := mocksset.NewMockSet(ctrl)
			mockSet.EXPECT().ID().Return(set.SetID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().UserID().Return(user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")).AnyTimes()
			mockSet.EXPECT().WorkoutID().Return(workout.WorkoutID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().ExerciseID().Return(exercise.ExerciseID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().Rep().Return(set.Rep(10)).AnyTimes()
			mockSet.EXPECT().Weight().Return(set.Weight(60.0)).AnyTimes()
			mockSet.EXPECT().TrainedAt().Return(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).AnyTimes()
			mockSet.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockSet)

			repo := NewSetRepository(gormDB)

			err = repo.Create(context.Background(), mockSet)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, s set.Set)
	}{
		{
			name:    "success update set",
			success: true,
			setup: func(mock sqlmock.Sqlmock, s set.Set) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sets" SET "user_id"=$1,"workout_id"=$2,"exercise_id"=$3,"rep"=$4,"weight"=$5,"trained_at"=$6,"created_at"=$7 WHERE "id" = $8`)).
					WithArgs(s.UserID(), s.WorkoutID(), s.ExerciseID(), s.Rep(), s.Weight(), testutil.AnyTime{}, testutil.AnyTime{}, s.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, s set.Set) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sets" SET "user_id"=$1,"workout_id"=$2,"exercise_id"=$3,"rep"=$4,"weight"=$5,"trained_at"=$6,"created_at"=$7 WHERE "id" = $8`)).
					WithArgs(s.UserID(), s.WorkoutID(), s.ExerciseID(), s.Rep(), s.Weight(), testutil.AnyTime{}, testutil.AnyTime{}, s.ID()).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to new sqlmock: %s", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSet := mocksset.NewMockSet(ctrl)
			mockSet.EXPECT().ID().Return(set.SetID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().UserID().Return(user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")).AnyTimes()
			mockSet.EXPECT().WorkoutID().Return(workout.WorkoutID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().ExerciseID().Return(exercise.ExerciseID{UUID: uuid.New()}).AnyTimes()
			mockSet.EXPECT().Rep().Return(set.Rep(10)).AnyTimes()
			mockSet.EXPECT().Weight().Return(set.Weight(60.0)).AnyTimes()
			mockSet.EXPECT().TrainedAt().Return(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).AnyTimes()
			mockSet.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockSet)

			repo := NewSetRepository(gormDB)
			err = repo.Update(context.Background(), mockSet)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	id := set.SetID{UUID: uuid.New()}
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, id set.SetID)
	}{
		{
			name:    "success delete set",
			success: true,
			setup: func(mock sqlmock.Sqlmock, id set.SetID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sets" WHERE id = $1`)).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure delete error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, id set.SetID) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sets" WHERE id = $1`)).
					WithArgs(id).
					WillReturnError(errors.New("delete error"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to new sqlmock: %s", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.setup(mock, id)

			repo := NewSetRepository(gormDB)
			err = repo.Delete(context.Background(), id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	id := set.SetID{UUID: uuid.New()}
	tests := []struct {
		name     string
		success  bool
		notFound bool
		setup    func(mock sqlmock.Sqlmock, id set.SetID)
	}{
		{
			name:    "success find by id",
			success: true,
			setup: func(mock sqlmock.Sqlmock, id set.SetID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "workout_id", "exercise_id", "rep", "weight", "trained_at", "created_at"}).
					AddRow(id, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", uuid.New().String(), "f1f538e5-4a37-409c-be99-09ee7bfefc50", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE id = $1 ORDER BY "sets"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:     "failure not found",
			success:  false,
			notFound: true,
			setup: func(mock sqlmock.Sqlmock, id set.SetID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE id = $1 ORDER BY "sets"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure other error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, id set.SetID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE id = $1 ORDER BY "sets"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to new sqlmock: %s", err)
			}
			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to open gorm: %s", err)
			}

			tt.setup(mock, id)

			repo := NewSetRepository(gormDB)
			_, err = repo.FindByID(context.Background(), id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.notFound && !errors.Is(err, set.ErrSetNotFound) {
				t.Errorf("expected ErrSetNotFound, got %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		userID  user.UserID
		setup   func(mock sqlmock.Sqlmock, userID user.UserID)
	}{
		{
			name:    "success find sets by user id",
			success: true,
			userID:  user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425"),
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "workout_id", "exercise_id", "rep", "weight", "trained_at", "created_at"}).
					AddRow(uuid.New().String(), userID, uuid.New().String(), "f1f538e5-4a37-409c-be99-09ee7bfefc50", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "success find sets empty result",
			success: true,
			userID:  user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425"),
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "workout_id", "exercise_id", "rep", "weight", "trained_at", "created_at"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure find sets error",
			success: false,
			userID:  user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425"),
			setup: func(mock sqlmock.Sqlmock, userID user.UserID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE user_id = $1`)).
					WithArgs(userID).
					WillReturnError(errors.New("find sets error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock, tt.userID)

			repo := NewSetRepository(gormDB)

			_, err = repo.FindByUserID(context.Background(), tt.userID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByWorkoutID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		workoutID workout.WorkoutID
		setup     func(mock sqlmock.Sqlmock, workoutID workout.WorkoutID)
	}{
		{
			name:      "success find sets by workout id",
			success:   true,
			workoutID: workout.WorkoutID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, workoutID workout.WorkoutID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "workout_id", "exercise_id", "rep", "weight", "trained_at", "created_at"}).
					AddRow(uuid.New().String(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", workoutID, "f1f538e5-4a37-409c-be99-09ee7bfefc50", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE workout_id = $1`)).
					WithArgs(workoutID).
					WillReturnRows(rows)
			},
		},
		{
			name:      "failure find sets by workout id error",
			success:   false,
			workoutID: workout.WorkoutID{UUID: uuid.New()},
			setup: func(mock sqlmock.Sqlmock, workoutID workout.WorkoutID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sets" WHERE workout_id = $1`)).
					WithArgs(workoutID).
					WillReturnError(errors.New("find sets error"))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock, tt.workoutID)

			repo := NewSetRepository(gormDB)

			_, err = repo.FindByWorkoutID(context.Background(), tt.workoutID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

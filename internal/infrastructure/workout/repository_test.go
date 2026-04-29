package workout

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

	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksworkout "github.com/qkitzero/workout-service/mocks/domain/workout"
	"github.com/qkitzero/workout-service/testutil"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, w workout.Workout)
	}{
		{
			name:    "success create workout",
			success: true,
			setup: func(mock sqlmock.Sqlmock, w workout.Workout) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "workouts" ("id","user_id","started_at","finished_at","created_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(w.ID(), w.UserID(), testutil.AnyTime{}, nil, testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure create error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, w workout.Workout) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "workouts" ("id","user_id","started_at","finished_at","created_at") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(w.ID(), w.UserID(), testutil.AnyTime{}, nil, testutil.AnyTime{}).
					WillReturnError(errors.New("create error"))
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

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().ID().Return(workout.WorkoutID{UUID: uuid.New()}).AnyTimes()
			mockWorkout.EXPECT().UserID().Return(user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")).AnyTimes()
			mockWorkout.EXPECT().StartedAt().Return(time.Now()).AnyTimes()
			mockWorkout.EXPECT().FinishedAt().Return(nil).AnyTimes()
			mockWorkout.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockWorkout)

			repo := NewWorkoutRepository(gormDB)
			err = repo.Create(context.Background(), mockWorkout)
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
	finishedAt := time.Now()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock, w workout.Workout)
	}{
		{
			name:    "success update workout",
			success: true,
			setup: func(mock sqlmock.Sqlmock, w workout.Workout) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "workouts" SET "user_id"=$1,"started_at"=$2,"finished_at"=$3,"created_at"=$4 WHERE "id" = $5`)).
					WithArgs(w.UserID(), testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}, w.ID()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "failure update error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, w workout.Workout) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "workouts" SET "user_id"=$1,"started_at"=$2,"finished_at"=$3,"created_at"=$4 WHERE "id" = $5`)).
					WithArgs(w.UserID(), testutil.AnyTime{}, testutil.AnyTime{}, testutil.AnyTime{}, w.ID()).
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

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().ID().Return(workout.WorkoutID{UUID: uuid.New()}).AnyTimes()
			mockWorkout.EXPECT().UserID().Return(user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")).AnyTimes()
			mockWorkout.EXPECT().StartedAt().Return(time.Now()).AnyTimes()
			mockWorkout.EXPECT().FinishedAt().Return(&finishedAt).AnyTimes()
			mockWorkout.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()

			tt.setup(mock, mockWorkout)

			repo := NewWorkoutRepository(gormDB)
			err = repo.Update(context.Background(), mockWorkout)
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
	id := workout.WorkoutID{UUID: uuid.New()}
	tests := []struct {
		name     string
		success  bool
		notFound bool
		setup    func(mock sqlmock.Sqlmock, id workout.WorkoutID)
	}{
		{
			name:    "success find by id",
			success: true,
			setup: func(mock sqlmock.Sqlmock, id workout.WorkoutID) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "started_at", "finished_at", "created_at"}).
					AddRow(id, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", time.Now(), nil, time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE id = $1 ORDER BY "workouts"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:     "failure not found",
			success:  false,
			notFound: true,
			setup: func(mock sqlmock.Sqlmock, id workout.WorkoutID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE id = $1 ORDER BY "workouts"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure other error",
			success: false,
			setup: func(mock sqlmock.Sqlmock, id workout.WorkoutID) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE id = $1 ORDER BY "workouts"."id" LIMIT $2`)).
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

			repo := NewWorkoutRepository(gormDB)
			_, err = repo.FindByID(context.Background(), id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.notFound && !errors.Is(err, workout.ErrWorkoutNotFound) {
				t.Errorf("expected ErrWorkoutNotFound, got %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestFindByUserID(t *testing.T) {
	t.Parallel()
	uid := user.UserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")

	t.Run("success no filters", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to new sqlmock: %s", err)
		}
		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm: %s", err)
		}

		rows := sqlmock.NewRows([]string{"id", "user_id", "started_at", "finished_at", "created_at"}).
			AddRow(uuid.New().String(), uid, time.Now(), nil, time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE user_id = $1`)).
			WithArgs(uid).
			WillReturnRows(rows)

		repo := NewWorkoutRepository(gormDB)
		_, err = repo.FindByUserID(context.Background(), uid, nil, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("success with from/to", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to new sqlmock: %s", err)
		}
		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm: %s", err)
		}

		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "user_id", "started_at", "finished_at", "created_at"})
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE user_id = $1 AND started_at >= $2 AND started_at < $3`)).
			WithArgs(uid, from, to).
			WillReturnRows(rows)

		repo := NewWorkoutRepository(gormDB)
		_, err = repo.FindByUserID(context.Background(), uid, &from, &to)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("failure db error", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to new sqlmock: %s", err)
		}
		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm: %s", err)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "workouts" WHERE user_id = $1`)).
			WithArgs(uid).
			WillReturnError(errors.New("db error"))

		repo := NewWorkoutRepository(gormDB)
		_, err = repo.FindByUserID(context.Background(), uid, nil, nil)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestExists(t *testing.T) {
	t.Parallel()
	id := workout.WorkoutID{UUID: uuid.New()}

	t.Run("success exists", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to new sqlmock: %s", err)
		}
		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm: %s", err)
		}

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "workouts" WHERE id = $1`)).
			WithArgs(id).
			WillReturnRows(rows)

		repo := NewWorkoutRepository(gormDB)
		ok, err := repo.Exists(context.Background(), id)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !ok {
			t.Errorf("expected exists=true")
		}
	})

	t.Run("failure exists error", func(t *testing.T) {
		t.Parallel()
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to new sqlmock: %s", err)
		}
		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to open gorm: %s", err)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "workouts" WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(errors.New("db error"))

		repo := NewWorkoutRepository(gormDB)
		_, err = repo.Exists(context.Background(), id)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

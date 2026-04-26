package exercise

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

func TestFindAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success find all exercises",
			success: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.MatchExpectationsInOrder(false)
				rows := sqlmock.NewRows([]string{"id", "code", "category", "created_at"}).
					AddRow("f1f538e5-4a37-409c-be99-09ee7bfefc50", "bench_press", "compound", nil)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercises"`)).
					WillReturnRows(rows)
				transRows := sqlmock.NewRows([]string{"exercise_id", "lang", "name"}).
					AddRow("f1f538e5-4a37-409c-be99-09ee7bfefc50", "ja", "ベンチプレス")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercise_translations"`)).
					WillReturnRows(transRows)
				joinRows := sqlmock.NewRows([]string{"exercise_id", "muscle_id"}).
					AddRow("f1f538e5-4a37-409c-be99-09ee7bfefc50", "4b5a784a-3333-4721-a071-2e3fbd570c7f")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercise_muscle"`)).
					WillReturnRows(joinRows)
				mgRows := sqlmock.NewRows([]string{"id", "code", "created_at"}).
					AddRow("4b5a784a-3333-4721-a071-2e3fbd570c7f", "chest", nil)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscles"`)).
					WillReturnRows(mgRows)
				mgTransRows := sqlmock.NewRows([]string{"muscle_id", "lang", "name"}).
					AddRow("4b5a784a-3333-4721-a071-2e3fbd570c7f", "ja", "胸")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscle_translations"`)).
					WillReturnRows(mgTransRows)
			},
		},
		{
			name:    "failure find all error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercises"`)).
					WillReturnError(errors.New("find all error"))
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

			tt.setup(mock)

			repo := NewExerciseRepository(gormDB)

			_, err = repo.FindAll(context.Background(), i18n.LanguageJa)
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
	idStr := "f1f538e5-4a37-409c-be99-09ee7bfefc50"
	id, err := exercise.NewExerciseIDFromString(idStr)
	if err != nil {
		t.Fatalf("failed to parse id: %v", err)
	}

	tests := []struct {
		name    string
		success bool
		wantErr error
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success find by id",
			success: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.MatchExpectationsInOrder(false)
				rows := sqlmock.NewRows([]string{"id", "code", "category", "created_at"}).
					AddRow(idStr, "bench_press", "compound", nil)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercises" WHERE id = $1 ORDER BY "exercises"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnRows(rows)
				transRows := sqlmock.NewRows([]string{"exercise_id", "lang", "name"}).
					AddRow(idStr, "ja", "ベンチプレス")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercise_translations"`)).
					WillReturnRows(transRows)
				joinRows := sqlmock.NewRows([]string{"exercise_id", "muscle_id"}).
					AddRow(idStr, "4b5a784a-3333-4721-a071-2e3fbd570c7f")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercise_muscle"`)).
					WillReturnRows(joinRows)
				mgRows := sqlmock.NewRows([]string{"id", "code", "created_at"}).
					AddRow("4b5a784a-3333-4721-a071-2e3fbd570c7f", "chest", nil)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscles"`)).
					WillReturnRows(mgRows)
				mgTransRows := sqlmock.NewRows([]string{"muscle_id", "lang", "name"}).
					AddRow("4b5a784a-3333-4721-a071-2e3fbd570c7f", "ja", "胸")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscle_translations"`)).
					WillReturnRows(mgTransRows)
			},
		},
		{
			name:    "failure not found",
			success: false,
			wantErr: exercise.ErrExerciseNotFound,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercises" WHERE id = $1 ORDER BY "exercises"."id" LIMIT $2`)).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "failure db error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exercises" WHERE id = $1 ORDER BY "exercises"."id" LIMIT $2`)).
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
				t.Errorf("failed to new sqlmock: %s", err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
			if err != nil {
				t.Errorf("failed to open gorm: %s", err)
			}

			tt.setup(mock)

			repo := NewExerciseRepository(gormDB)

			_, err = repo.FindByID(context.Background(), id, i18n.LanguageJa)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	t.Parallel()
	id, err := exercise.NewExerciseIDFromString("f1f538e5-4a37-409c-be99-09ee7bfefc50")
	if err != nil {
		t.Fatalf("failed to parse id: %v", err)
	}

	tests := []struct {
		name    string
		success bool
		want    bool
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success exists true",
			success: true,
			want:    true,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "exercises" WHERE id = $1`)).
					WithArgs(id).
					WillReturnRows(rows)
			},
		},
		{
			name:    "success exists false",
			success: true,
			want:    false,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "exercises" WHERE id = $1`)).
					WithArgs(id).
					WillReturnRows(rows)
			},
		},
		{
			name:    "failure exists error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "exercises" WHERE id = $1`)).
					WithArgs(id).
					WillReturnError(errors.New("exists error"))
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

			tt.setup(mock)

			repo := NewExerciseRepository(gormDB)

			got, err := repo.Exists(context.Background(), id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

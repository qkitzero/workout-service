package muscle

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

func TestFindAll(t *testing.T) {
	t.Parallel()
	idStr := "4b5a784a-3333-4721-a071-2e3fbd570c7f"
	tests := []struct {
		name    string
		success bool
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name:    "success find all muscles",
			success: true,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "code", "created_at"}).
					AddRow(idStr, "chest", nil)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscles"`)).
					WillReturnRows(rows)
				transRows := sqlmock.NewRows([]string{"muscle_id", "lang", "name"}).
					AddRow(idStr, "ja", "胸")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscle_translations" WHERE "muscle_translations"."muscle_id" = $1`)).
					WillReturnRows(transRows)
			},
		},
		{
			name:    "failure find all error",
			success: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "muscles"`)).
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

			repo := NewMuscleRepository(gormDB)

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

package muscle

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

type MuscleModel struct {
	ID           muscle.MuscleID
	Code         muscle.Code
	CreatedAt    time.Time
	Translations []MuscleTranslationModel `gorm:"foreignKey:MuscleID"`
}

func (MuscleModel) TableName() string {
	return "muscles"
}

func (m MuscleModel) ToDomain(lang i18n.Language) muscle.Muscle {
	translations := make([]muscle.Translation, len(m.Translations))
	for i, t := range m.Translations {
		translations[i] = muscle.NewTranslation(t.Lang, t.Name)
	}
	name := muscle.ResolveName(translations, lang, m.Code)
	return muscle.NewMuscle(m.ID, m.Code, name)
}

type MuscleTranslationModel struct {
	MuscleID muscle.MuscleID
	Lang     i18n.Language
	Name     muscle.Name
}

func (MuscleTranslationModel) TableName() string {
	return "muscle_translations"
}

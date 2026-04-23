package exercise

type Exercise interface {
	ID() ExerciseID
	Code() Code
	Category() Category
	Translations() []Translation
	Name(lang Language) Name
}

type exercise struct {
	id           ExerciseID
	code         Code
	category     Category
	translations []Translation
}

func (e exercise) ID() ExerciseID               { return e.id }
func (e exercise) Code() Code                   { return e.code }
func (e exercise) Category() Category           { return e.category }
func (e exercise) Translations() []Translation  { return e.translations }

func (e exercise) Name(lang Language) Name {
	for _, t := range e.translations {
		if t.Language() == lang {
			return t.Name()
		}
	}
	for _, t := range e.translations {
		if t.Language() == LanguageJa {
			return t.Name()
		}
	}
	if len(e.translations) > 0 {
		return e.translations[0].Name()
	}
	return Name(e.code.String())
}

func NewExercise(id ExerciseID, code Code, category Category, translations []Translation) Exercise {
	return &exercise{
		id:           id,
		code:         code,
		category:     category,
		translations: translations,
	}
}

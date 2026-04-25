package muscle

import (
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type Translation struct {
	language i18n.Language
	name     Name
}

func (t Translation) Language() i18n.Language { return t.language }
func (t Translation) Name() Name              { return t.name }

func NewTranslation(language i18n.Language, name Name) Translation {
	return Translation{language: language, name: name}
}

func ResolveName(translations []Translation, lang i18n.Language, fallbackCode Code) Name {
	for _, t := range translations {
		if t.Language() == lang {
			return t.Name()
		}
	}
	for _, t := range translations {
		if t.Language() == i18n.LanguageJa {
			return t.Name()
		}
	}
	return Name(fallbackCode.String())
}

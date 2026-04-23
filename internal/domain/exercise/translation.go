package exercise

type Translation struct {
	language Language
	name     Name
}

func (t Translation) Language() Language { return t.language }
func (t Translation) Name() Name         { return t.name }

func NewTranslation(language Language, name Name) Translation {
	return Translation{language: language, name: name}
}

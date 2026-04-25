package muscle

type Muscle interface {
	ID() MuscleID
	Code() Code
	Name() Name
}

type muscle struct {
	id   MuscleID
	code Code
	name Name
}

func (m muscle) ID() MuscleID { return m.id }
func (m muscle) Code() Code   { return m.code }
func (m muscle) Name() Name   { return m.name }

func NewMuscle(id MuscleID, code Code, name Name) Muscle {
	return &muscle{
		id:   id,
		code: code,
		name: name,
	}
}

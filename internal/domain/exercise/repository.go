package exercise

type ExerciseRepository interface {
	FindAll() ([]Exercise, error)
	FindByID(id ExerciseID) (Exercise, error)
	Exists(id ExerciseID) (bool, error)
}

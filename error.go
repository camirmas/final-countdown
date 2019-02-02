package countdown

type TimerNotFoundError struct{}
type TimerExistsError struct{}

func (e TimerNotFoundError) Error() string {
	return "Timer not found"
}

func (e TimerExistsError) Error() string {
	return "Timer already exists"
}

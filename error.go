package countdown

type TimerNotFoundError struct{}

func (e TimerNotFoundError) Error() string {
	return "Timer not found"
}

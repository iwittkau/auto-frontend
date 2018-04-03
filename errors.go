package auto_frontend

// Button related errors
const (
	ErrButtonAlreadyRegistered = Error("button already registered")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

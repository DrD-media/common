package errors

import "fmt"

type CustomError struct {
	Message string
	Err     error
}

func New(message string) error {
	return &CustomError{Message: message}
}
func Is() {

}
func Wrap(err error, message string) error {
	return &CustomError{Message: message, Err: err}
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

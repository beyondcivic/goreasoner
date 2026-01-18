package reasoner

import "fmt"

type AppError struct {
	// Message to show the user.
	Message string
	// Value to include with message
	Value any
}

func (e AppError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Value)
	} else {
		return e.Message
	}
}

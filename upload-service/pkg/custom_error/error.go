package customerror

import "fmt"

type SystemCommunicationError struct {
	Reason string
}

func (re *SystemCommunicationError) Error() string {
	return fmt.Sprintf("System Error: %s", re.Reason)
}

type SafeError struct {
	// Used for HTTP Responses
	OutwardMessage string
	// Can be used for detailed logging
	InternalError error
}

func (se *SafeError) Error() string {
	return fmt.Sprintf("Encountered Error %s, will be returned to user as %s", se.InternalError.Error(), se.OutwardMessage)
}

func NewSafeErrorFromError(err error) *SafeError {
	return &SafeError{InternalError: err, OutwardMessage: err.Error()}
}

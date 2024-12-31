package customerror

import "fmt"

type SystemCommunicationError struct {
	Reason string
}

func (re *SystemCommunicationError) Error() string {
	return fmt.Sprintf("System Error: %s", re.Reason)
}

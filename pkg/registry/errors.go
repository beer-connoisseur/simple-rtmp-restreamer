package registry

import "fmt"

var (
	ErrStreamNotExist = "ErrStreamNotExist"
)

type StreamNotFoundError struct{}

func (e StreamNotFoundError) Error() string {
	return fmt.Sprintf("%s", ErrStreamNotExist)
}

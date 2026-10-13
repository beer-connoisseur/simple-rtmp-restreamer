package registry

import "fmt"

var (
	ErrStreamNotExist = "StreamNotExist"
)

type StreamNotFoundError struct{}

func (e StreamNotFoundError) Error() string {
	return fmt.Sprintf("%s", ErrStreamNotExist)
}

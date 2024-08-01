package exception

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type ConnectionException struct {
	Message    string
	StackTrace string
	Timestamp  time.Time
}

func NewConnectionException(message string, err error) *ConnectionException {
	return &ConnectionException{
		Message:    message,
		StackTrace: fmt.Sprintf("%+v", errors.WithStack(err)),
		Timestamp:  time.Now(),
	}
}

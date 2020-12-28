package errors

import (
	"fmt"
)

type CDKRuntimeError struct {
	Err       error
	CustomMsg string
}

func New(text string) error {
	return &CDKRuntimeError{nil, text}
}

func (e *CDKRuntimeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s:\n%s", e.CustomMsg, e.Err)
	} else {
		return e.CustomMsg
	}
}

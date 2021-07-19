package cerrors

import (
	"fmt"
	"io"
	"strings"
)

type withOwnMessage interface {
	OwnMessage() string
}

type withFrame interface {
	Frame() Frame
}

type withUnwrap interface {
	Unwrap() error
}

func Format(err error, s io.Writer, withFunctions bool) {
	var prefix string

	for err != nil {
		next := getNext(err)
		fmt.Fprintln(s, prefix+getOwnMessage(err, next))

		if wFrame, ok := err.(withFrame); ok {
			wFrame.Frame().Format(s, withFunctions)
		}

		prefix = "  - "
		err = next
	}
}

func getOwnMessage(err, next error) string {
	if wMessage, ok := err.(withOwnMessage); ok {
		return wMessage.OwnMessage()
	} else if next != nil {
		return unwrapMessageStr(err.Error(), next.Error())
	} else {
		return err.Error()
	}
}

func getNext(err error) error {
	if wUnwrap, ok := err.(withUnwrap); ok {
		return wUnwrap.Unwrap()
	}
	return nil
}

// UnwrapMessage returns the portion of the error message which does not repeat
// in the wrapped error.
//
// Examples:
//
// err := errors.New("couldn't")
//
// UnwrapMessage(fmt.Errorf("doing foo: %w", err)) // => "doing foo"
// UnwrapMessage(fmt.Errorf("doing foo %w", err)) // => "doing foo"
// UnwrapMessage(fmt.Errorf("failed (%w) while doing foo", err)))  // => "failed (*) while doing foo"
func UnwrapMessage(err error) string {
	msg := err.Error()
	if wrapper, ok := err.(withUnwrap); ok {
		wrapped := wrapper.Unwrap()
		if wrapped != nil {
			if withOwnMessage, ok := wrapped.(withOwnMessage); ok {
				return withOwnMessage.OwnMessage()
			}
			wrappedMsg := wrapped.Error()
			return unwrapMessageStr(msg, wrappedMsg)
		}
	}

	return msg
}

func unwrapMessageStr(parent, child string) string {
	msg := parent
	if idx := strings.Index(parent, child); idx > 0 {
		if idx+len(child) == len(parent) {
			msg = strings.TrimRight(parent[:idx], ": ")
		} else {
			msg = parent[:idx] + "*" + parent[idx+len(child):]
		}
	}
	return msg
}

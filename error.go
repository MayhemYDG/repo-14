package cerrors

import (
	"fmt"
)

func New(msg string) error {
	return &cerror{
		msg:   msg,
		frame: Caller(1),
	}
}

func Newf(msg string, args ...interface{}) error {
	return &cerror{
		msg:   fmt.Sprintf(msg, args...),
		frame: Caller(1),
	}
}

func Wrap(err error, msg string) error {
	return &cerror{
		msg:   msg,
		next:  err,
		frame: Caller(1),
	}
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return &cerror{
		msg:   fmt.Sprintf(msg, args...),
		next:  err,
		frame: Caller(1),
	}
}

type cerror struct {
	msg   string
	next  error
	frame Frame
}

func (c *cerror) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') || s.Flag('#') {
			Format(c, s, s.Flag('#'))
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, c.Error())
	default:
		fmt.Fprintf(s, "%%!%s(cerror)", string(verb))
	}
}

func (c *cerror) OwnMessage() string {
	return c.msg
}

func (c *cerror) Frame() Frame {
	return c.frame
}

func (c *cerror) Error() string {
	if c.next != nil {
		return c.msg + ": " + c.next.Error()
	}
	return c.msg
}

func (c *cerror) Unwrap() error {
	return c.next
}

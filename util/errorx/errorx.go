package errorx

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type Error struct {
	file string
	line int
	msg  string
}

func New(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	_, file, line, _ := runtime.Caller(1)
	e := &Error{
		file: file,
		line: line,
	}
	if err != nil {
		e.msg = err.Error()
	} else {
		e.msg = "unknown error"
	}
	fmt.Printf("Error: %s\n", e.msg)
	debug.PrintStack()
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d %s", e.file, e.line, e.msg)
}

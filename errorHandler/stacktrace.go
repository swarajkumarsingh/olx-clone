package errorHandler

import (
	"fmt"
	"runtime"
)

// MiddlewareError represet error codes
type MiddlewareError struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// NoEligibilityError User is not eligible for loan
var NoEligibilityError = MiddlewareError{
	"FLM601",
	"No data to compute loan eligibility for user.",
}

const maxStackDepth = 64

type stack []uintptr
type traceableError struct {
	err   error
	stack *stack
}

type traceable interface {
	// Cause returns original error.
	Cause() error
	// Stack returns stack trace of this error.
	Stack() *stack
}

func (t traceableError) Error() string {
	if t.err == nil {
		return ""
	}
	return t.err.Error()
}

func (t traceableError) Cause() error {
	return t.err
}

func (t traceableError) Stack() *stack {
	return t.stack
}

// StackTrace returns the stack trace with given error
// If the error is not traceable, empty string is returned
func StackTrace(err error) string {
	if tracer, ok := err.(traceable); ok {
		return tracer.Stack().Format()
	}
	return ""
}

// Cause returns the original error of a traceable error
// If the error is not traceable, it returns itself
func Cause(err error) error {
	if tracer, ok := err.(traceable); ok {
		return tracer.Cause()
	}
	return err
}

// Format formats the each trace frame into a string
func (s *stack) Format() (ret string) {
	if s == nil {
		return
	}
	v := *s
	if len(v) == 0 {
		return
	}
	frames := runtime.CallersFrames(v)
	for {
		frame, more := frames.Next()
		ret = ret + fmt.Sprintf("%s(%s:%s); ", frame.Function, frame.File, fmt.Sprint(frame.Line))
		if !more {
			return
		}
	}
}

// Wrap wraps an error into a traceable error.
// If the wrapped error is traceable, do nothing.
// If the wrapped error is not traceable, record the stack trace.
func Wrap(err error) error {
	if _, ok := err.(traceable); ok {
		return err
	}
	return traceableError{
		err:   err,
		stack: recordStack(),
	}
}

func recordStack() *stack {
	s := make(stack, maxStackDepth)
	n := runtime.Callers(3, s) // skipping first three frames to skip calling functions of stacktrace.go from appearing in the stack trace
	s = s[:n]
	return &s
}

package wraperr

import (
	"errors"
	"runtime"
	"strings"
)

type WrapErr struct {
	curErr     error
	prevErr    *WrapErr
	stackTrace StackTrace
}

type StackTrace []Frame
type Frame uintptr

func Wrap(curErr error, prevErr error) *WrapErr {
	if curErr == nil {
		panic("wrap_err: curErr cannot be nil")
	}
	res := &WrapErr{
		curErr: curErr,
	}
	var e *WrapErr
	if errors.As(prevErr, &e) {
		res.prevErr = e
		res.stackTrace = e.stackTrace
	} else {
		if prevErr == nil {
			prevErr = errors.New("nil")
		}
		res.prevErr = &WrapErr{
			curErr:  prevErr,
			prevErr: nil,
		}

		pc := make([]uintptr, 10)
		n := runtime.Callers(2, pc)
		if n != 0 {
			pc = pc[:n]
			frames := runtime.CallersFrames(pc)
			for {
				frame, more := frames.Next()
				res.stackTrace = append(res.stackTrace, Frame(frame.PC))
				if !more {
					break
				}
			}
		}
	}
	return res
}

func (e WrapErr) Is(target error) bool {
	return e.curErr.Error() == target.Error()
}

func (e *WrapErr) As(target interface{}) bool {
	return errors.As(e.curErr, target)
}

func (e WrapErr) ContainsError(target error) bool {
	return errors.Is(e.curErr, target)
}

func (e WrapErr) ContainsType(target interface{}) bool {
	return errors.As(e.curErr, target)
}

func (e WrapErr) String() string {
	return e.curErr.Error()
}

func (e WrapErr) Unwrap() error {
	if e.prevErr != nil {
		return e.prevErr
	} else {
		return nil
	}
}

func (e WrapErr) Error() string {
	if e.curErr == nil {
		return ""
	}
	var str strings.Builder
	str.WriteString(e.curErr.Error())
	err := e.prevErr
	for err != nil {
		str.WriteString(": " + err.String())
		err = err.prevErr
	}
	return str.String()
}

func (e WrapErr) StackTrace() StackTrace {
	return e.stackTrace
}

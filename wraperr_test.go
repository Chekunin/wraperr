package wraperr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrappingError(t *testing.T) {
	err := errors.New("error1")
	err = NewWrapErr(fmt.Errorf("error2"), err)
	err = NewWrapErr(fmt.Errorf("error3"), err)
	err = NewWrapErr(fmt.Errorf("error4"), err)
	assert.Equal(t, "error4: error3: error2: error1", err.Error())

	err2 := NewWrapErr(fmt.Errorf("error1"), nil)
	assert.Equal(t, "error1: nil", err2.Error())

	err3 := NewWrapErr(fmt.Errorf("error2"), fmt.Errorf("error1"))
	err3 = NewWrapErr(fmt.Errorf("error3"), err3)
	assert.Equal(t, "error2: error1", err3.Unwrap().Error())
	assert.Equal(t, "error3", err3.String())
}

func TestWrapEmptyError(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotEmpty(t, err)
	}()
	NewWrapErr(nil, errors.New("error1"))
	assert.Equal(t, true, false)
}

type customErr struct {
	Code    int
	Message string
}

func (err customErr) Error() string {
	return fmt.Sprintf("%d %s", err.Code, err.Message)
}

func TestAs(t *testing.T) {
	var err error
	err = customErr{
		Code:    1,
		Message: "Test message",
	}
	err = NewWrapErr(fmt.Errorf("wrapper"), err)
	var custErr customErr
	if errors.As(err, &custErr) {
		assert.Equal(t, customErr{
			Code:    1,
			Message: "Test message",
		}.Error(), custErr.Error())
	} else {
		assert.Fail(t, "err is not type if customErr")
	}
}

func TestIs(t *testing.T) {
	custErr := customErr{
		Code:    1,
		Message: "Test message",
	}
	var err error
	err = NewWrapErr(fmt.Errorf("error1"), nil)
	err = NewWrapErr(custErr, err)
	err = NewWrapErr(fmt.Errorf("error2"), err)
	if !errors.Is(err, custErr) {
		assert.Fail(t, "err does not contain custErr")
	}
}

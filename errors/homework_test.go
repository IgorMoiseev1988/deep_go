package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return ""
	}

	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d errors occured:\n", len(e.errs))
	for _, err := range e.errs {
		fmt.Fprintf(&sb, "\t* %v", err)
	}
	fmt.Fprintf(&sb, "\n")
	return sb.String()
}

func Append(err error, errs ...error) *MultiError {
	var multiErr *MultiError
	if err != nil {
		me, ok := err.(*MultiError)
		if ok {
			multiErr = me
		} else {
			multiErr = &MultiError{errs: []error{err}}
		}
	} else {
		multiErr = &MultiError{}
	}

	for _, e := range errs {
		if e != nil {
			multiErr.errs = append(multiErr.errs, e)
		}
	}
	return multiErr
}

func (e *MultiError) As(target interface{}) bool {
	if e == nil {
		return false
	}

	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Is(target error) bool {
	if e == nil {
		return false
	}

	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *MultiError) Unwrap() error {
	return nil
}

type EmptyError struct {
	i int
}

type EmptyError2 struct {
	i int
}

func (e *EmptyError) Error() string {
	return fmt.Sprint(e.i)
}

func (e *EmptyError2) Error() string {
	return fmt.Sprint(e.i)
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	err = Append(err, &EmptyError{i: 5})
	var ee *EmptyError
	var ee2 *EmptyError2
	as_ok := errors.As(err, &ee)
	assert.True(t, as_ok)
	assert.EqualError(t, ee, "5")

	as_ok = errors.As(err, &ee2)
	assert.False(t, as_ok)

	is_ok := errors.Is(err, ee)
	assert.True(t, is_ok)
	assert.EqualError(t, ee, "5")
	is_ok = errors.Is(err, ee2)
	assert.False(t, is_ok)

}

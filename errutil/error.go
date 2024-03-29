// Package errutil implements some useful extensions to the stdlib
// Go errors package, in the same spirit of packages like ioutil/httputil.
//
// Utilities include:
//
// - An error type that makes it easy to work with const error sentinels.
// - An easy way to wrap a list of errors together.
// - An easy way to reduce a list of errors.
//
// Flexible enough that you can do your own wrapping/merging logic
// but in a functional/simple way.
package errutil

import "errors"

// Error implements the Go's error interface in the simplest
// way possible, allowing initialization error sentinels to be done
// at compile time as constants. It does so by using a string
// as it's base type.
type Error string

// Reducer reduces 2 errors into one
type Reducer func(error, error) error

// Error return a string representation of the error.
func (e Error) Error() string {
	return string(e)
}

// Chain creates a chain of errors suitable to be used
// with Go's Unwrap interface through functions like
// errors.Is and errors.As.
// Chaining order will be the same as the order of the
// arguments, the first error is the head wrapping up
// the next one, and so goes on.
//
// An empty list of errors will return a nil error.
func Chain(errs ...error) error {
	errs = removeNils(errs)

	if len(errs) == 0 {
		return nil
	}

	return errorChain{
		head: errs[0],
		tail: Chain(errs[1:]...),
	}
}

// Reduce will reduce all errors to a single one using the
// provided reduce function.
//
// If errs is empty it returns nil, if errs has a single err
// (len(errs) == 1) it will return the err itself.
//
// Nil errors on the errs args will be filtered out initially,
// before reducing, so you can expect errors passed to the reducer
// to be always non-nil.
//
// But if the reducer function itself returns nil, then the returned nil
// won't be filtered and will be passed as an argument on the next
// reducing step.
func Reduce(r Reducer, errs ...error) error {
	errs = removeNils(errs)
	return reduce(r, errs...)
}

type errorChain struct {
	head error
	tail error
}

// Error return a string representation of the chain of errors.
func (e errorChain) Error() string {
	if e.tail == nil {
		return e.head.Error()
	}
	return e.head.Error() + ": " + e.tail.Error()
}

func (e errorChain) Unwrap() error {
	return e.tail
}

func (e errorChain) Is(target error) bool {
	return errors.Is(e.head, target)
}

func (e errorChain) As(target interface{}) bool {
	return errors.As(e.head, target)
}

func removeNils(errs []error) []error {
	res := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}

func reduce(r Reducer, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	err1, err2 := errs[0], errs[1]
	err := r(err1, err2)
	return reduce(r, append([]error{err}, errs[2:]...)...)
}

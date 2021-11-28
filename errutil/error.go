// Package errutil implements some useful extensions to the stdlib
// Go errors package, in the same spirit of packages like ioutil/httputil.
//
// Utilities include:
//
// - An error type that makes it easy to work with const error sentinels.
// - An easy way to wrap a list of errors together.
// - An easy way to merge a list of errors opaquely.
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
	// TODO(katcipis): should we do something when
	// in the middle of the errs slice we have nils ?
	// prone to filtering nils out, or they will break the chain anyway.
	if len(errs) == 0 {
		return nil
	}
	return errorChain{
		head: errs[0],
		tail: Chain(errs[1:]...),
	}
}

// Merge reduces all given errors to a single opaque error
// instance containing the original information from all
// errors aggregated but in an opaque way (opposed to Chain).
//
// If the len(errs) == 0 it returns nil, if len(errs) == 1 it
// returns the error itself and it will filter out nil errs.
//
// It is specially useful to manage lists of errors and
// return an error only if any of the errors failed.
func Merge(errs ...error) error {
	return Reduce(func(err1, err2 error) error {
		// TODO: add nil filtering

		return errors.New(err1.Error() + ": " + err2.Error())
	}, errs...)
}

// Reduce will reduce all errors to a single one using the
// provided reduce function.
//
// If errs is empty it returns nil, if errs has a single err
// (len(errs) == 1) it will return the err itself.
//
// It won't assume anything else about the given errs, always
// calling the reduce function, so nil errs on will be passed
// to the reduce function so it can deal with them.
//
// Reduce will panic if the given reduce function panics.
func Reduce(reduce Reducer, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	err1, err2 := errs[0], errs[1]
	err := reduce(err1, err2)
	return Reduce(reduce, append([]error{err}, errs[2:]...)...)
}

type errorChain struct {
	head error
	tail error
}

// Error return a string representation of the chain of errors.
func (e errorChain) Error() string {
	if e.head == nil {
		return ""
	}
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

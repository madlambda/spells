// Package errutil implements some useful extensions to the stdlib
// Go errors package, in the same spirit of packages like ioutil/httputil.
//
// Utilities include:
//
// - An error type that makes it easy to work with const error sentinels.
package errutil

// Error implements the Go's error interface in the simplest
// way possible, allowing initialization error sentinels to be done
// at compile time as constants. It does so by using a string
// as it's base type.
type Error string

// ErrorChain implements the Go's Unwrap interface, representing
// a chain of errors. It is usually built by using the Chain
// function and you rarely need to manipulate it directly,
// just use errors.Is, errors.As functions from go stdlib.
type ErrorChain struct {
	// Head is the head error of the chain.
	Head error
	// Tail is the rest of the chain, nil if this is the last err.
	Tail error
}

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
	return ErrorChain{
		Head: errs[0],
		Tail: Chain(errs[1:]...),
	}
}

// Error return a string representation of the chain of errors.
func (e ErrorChain) Error() string {
	return "TODO"
}

// Unwrap return the wrapped error or nil if there is none.
func (e ErrorChain) Unwrap() error {
	return e.Tail
}

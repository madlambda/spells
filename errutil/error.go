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

// Error return a string representation of the error.
func (e Error) Error() string {
	return string(e)
}

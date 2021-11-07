// Package errutil implements some useful extensions to the stdlib
// Go errors package, in the same spirit of packages like ioutil/httputil.
//
// Utilities include:
//
// - An error type that makes it easy to work with const error sentinels.
// - Improved chaining of error sentinels (wrapping multiple errors).
package errutil

type Error string

func (e Error) Error() string {
	return string(e)
}

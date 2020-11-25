// Package iotest implements io Read/Writers that are useful for testing.
// It is inspired on the stdlib iotest: https://golang.org/pkg/testing/iotest/
// It adds on it, instead of being a replacement.
package iotest

import "io"

// RepeaterReader is an io.Reader that repeats a given io.Reader
type RepeaterReader struct {
}

// NewRepeater creates RepeaterReader that will repeat the
// given reader "n" times. If n is 0 it will repeat the
// given stream forever (and infinite stream).
func NewRepeater(r io.Reader, n uint) *RepeaterReader {
	return nil
}

func (r *RepeaterReader) Read(d []byte) (int, error) {
	return 0, nil
}

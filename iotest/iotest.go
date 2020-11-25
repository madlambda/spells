// Package iotest implements io Read/Writers that are useful for testing.
// It is inspired on the stdlib iotest: https://golang.org/pkg/testing/iotest/
// It adds on it, instead of being a replacement.
package iotest

import "io"

// RepeaterReader is an io.Reader that repeats a given io.Reader
type RepeaterReader struct {
	reader      io.Reader
	repeatCount uint
}

// NewRepeater creates RepeaterReader that will repeat the
// given reader "n" times. If n=0 it won't repeat it and will only
// provide the contents of the given reader once. If n=1 it will repeat
// once, duplicating the input.
func NewRepeater(r io.Reader, n uint) *RepeaterReader {
	return &RepeaterReader{
		reader:      r,
		repeatCount: n,
	}
}

func (r *RepeaterReader) Read(d []byte) (int, error) {
	return r.reader.Read(d)
}

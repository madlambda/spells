// Package iotest implements io Read/Writers that are useful for testing.
// It is inspired on the stdlib iotest: https://golang.org/pkg/testing/iotest/
// It adds on it, instead of being a replacement.
package iotest

import "io"

// RepeaterReader is an io.Reader that repeats a given io.Reader
type RepeaterReader struct {
	reader      io.Reader
	readData    []byte
	readIndex   int
	err         error
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
	if r.err == nil {
		n, err := r.reader.Read(d)
		r.err = err
		r.readData = append(r.readData, d[:n]...)
		if err == io.EOF && r.repeatCount > 0 {
			return n, nil
		}
		return n, err
	}
	// TODO: handle err is not EOF

	if r.repeatCount == 0 {
		return 0, r.err
	}

	n := copy(d, r.readData[r.readIndex:])
	r.readIndex += n

	if r.readIndex >= len(r.readData) {
		r.readIndex = 0
		r.repeatCount -= 1
	}
	return n, nil
}

// Package iotest implements io Read/Writers that are useful for testing.
// It is inspired on the stdlib iotest: https://golang.org/pkg/testing/iotest/
// It adds on it, instead of being a replacement.
package iotest

import (
	"errors"
	"io"
)

// RepeaterReader is an io.Reader that repeats a given io.Reader
// It is NOT safe to use a RepeaterReader concurrently.
type RepeaterReader struct {
	reader      io.Reader
	readData    []byte
	readIndex   int
	err         error
	repeatCount uint
}

// BrokenReader is an io.Reader that always fails
// It is safe to use a BrokenReader concurrently.
type BrokenReader struct {
	// Err is the error you want to inject on Read calls.
	// If it is nil a default error is going to be returned on the Read call.
	Err error
}

// NewRepeater creates RepeaterReader that will repeat the
// given reader "n" times. If n=0 it won't repeat it and will only
// provide the contents of the given reader once. If n=1 it will repeat
// once, duplicating the input.
//
// It will use O(N) memory where N is the size of the contents read from
// the given reader (all contents are kept in memory and looped over).
//
// It can be a cheap way to generate gigantic inputs by repeating a very
// small input.
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

	if r.err != io.EOF {
		return 0, r.err
	}

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

func (r BrokenReader) Read(d []byte) (int, error) {
	if r.Err == nil {
		r.Err = errors.New("BrokenReaderError")
	}
	return 0, r.Err
}

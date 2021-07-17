// Package iotest implements io Read/Writers that are useful for testing.
// It is inspired on the stdlib iotest: https://golang.org/pkg/testing/iotest/
// It adds on it, instead of being a replacement.
package iotest

import (
	"fmt"
	"io"
)

// RepeatReader is an io.Reader that repeats a given io.Reader
// It is NOT safe to use a RepeatReader concurrently.
type RepeatReader struct {
	reader      io.Reader
	readData    []byte
	readIndex   int
	err         error
	repeatCount int
}

// NewRepeatReader creates RepeatReader that will repeat the
// given reader "n" times. If n=0 it won't repeat it and will only
// provide the contents of the given reader once. If n=1 it will repeat
// once, duplicating the input.
//
// It will use O(N) memory where N is the size of the contents read from
// the given reader (all contents are kept in memory and looped over).
//
// It can be a cheap way to generate gigantic inputs by repeating a very
// small input.
func NewRepeatReader(r io.Reader, n int) *RepeatReader {
	return &RepeatReader{
		reader:      r,
		repeatCount: n,
	}
}

func (r *RepeatReader) Read(d []byte) (int, error) {
	if r.repeatCount == 0 {
		if r.err == nil {
			r.err = io.EOF
		}
		return 0, r.err
	}

	if r.err == nil {
		n, err := r.reader.Read(d)
		r.err = err
		r.readData = append(r.readData, d[:n]...)

		fmt.Println(err)
		if err == io.EOF {
			r.repeatCount -= 1
			return n, nil
		}
		return n, err
	}

	if r.err != io.EOF {
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

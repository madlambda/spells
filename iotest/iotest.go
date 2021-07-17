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

// RepeatReaderError is an enumeration of errors that are part of
// RepeatReader interface. The sentinel errors are wrapped
// with dynamic context, so always check using errors.Is.
type RepeatReaderError string

const (
	RepeatReaderInvalidCount RepeatReaderError = "Repeat reader count must be >= 0"
)

// NewRepeatReader creates RepeatReader that will repeat the
// given io.Reader "n" times.
//
// If n=0 it won't provide any data and will return an io.EOF immediately (empty stream).
// If n=1 it will repeat once, just like reading the original stream.
// If n=2 it will repeat the underlying stream twice, doubling it.
//
// It will use O(N) memory where N is the size of the contents read from
// the given reader (all contents are kept in memory and looped over).
//
// It can be a cheap way to generate gigantic inputs by repeating a very
// small input (it will only use O(N) where N= small input size).
//
// It is a severe programming error to pass a negative count, resulting in a panic.
func NewRepeatReader(r io.Reader, n int) *RepeatReader {
	if n < 0 {
		panic(fmt.Errorf("%w : count is %d", RepeatReaderInvalidCount, n))
	}
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

func (r RepeatReaderError) Error() string {
	return string(r)
}

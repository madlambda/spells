package runes

import (
	"fmt"
	"io"
	"unicode"
)

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(data) runes into data. It returns the number of runes
// read (0 <= n <= len(data)) and any error encountered. Even if Read
// returns n < len(data), it may use all of data as scratch space during the
// call. If some data is available but not len(data) runes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 runes, it returns the number of
// runes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of runes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 runes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some runes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero rune count with a nil error, except when len(data) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain data.
type Reader interface {
	Read(data []rune) (int, error)
}

// UnicodeReader wraps an `io.RuneReader` in a `Reader` interface.
type UnicodeReader struct {
	r io.RuneReader
}

// NewUnicodeReader implements a runes.Reader interface from a io.RuneReader
// interface.
func NewUnicodeReader(r io.RuneReader) *UnicodeReader {
	return &UnicodeReader{
		r: r,
	}
}

// Read implements the runes.Reader interface.
func (r *UnicodeReader) Read(data []rune) (int, error) {
	offset := 0
	for i := 0; i < len(data); i++ {
		r, size, err := r.r.ReadRune()
		if err != nil {
			return i, err
		}
		if r == unicode.ReplacementChar {
			return i, fmt.Errorf("invalid rune at offset %d", offset)
		}
		offset += size
		data[i] = r
	}

	return len(data), nil
}

// ReadAll reads from r until an error or EOF and returns the data.
func ReadAll(r Reader) ([]rune, error) {
	b := make([]rune, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

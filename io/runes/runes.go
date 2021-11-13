package runes

import (
	"fmt"
	"io"
	"unicode"
)

// Reader reads Unicode encoded bytes into data.
// For the details about the usage of such readers, see the stdlib io.Reader
// documentation.
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

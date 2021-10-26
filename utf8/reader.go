package utf8

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type (
	RuneReader struct {
		r io.RuneReader
	}

	ReaderReader struct {
		r io.Reader
	}

	Reader interface {
		Read(data []rune) (int, error)
	}
)

// NewReaderReader creates a utf8.Reader implementation from an io.Reader
// interface.
func NewReaderReader(r io.Reader) *ReaderReader {
	return &ReaderReader{
		r: r,
	}
}

// Read implements utf8.Reader interface.
//
// It will use O(N) space where N is the size of the data slice provided.
// For ASCII data it allocates exactly len(data) bytes (or sizeof(data)/4) but
// in case of multi-byte code points it allocates a maximum of sizeof(data)
// or len(data)*4.
//
// The first call to the source reader uses a full buffer of len(data) bytes
// but subsequent calls read byte by byte. This is a tradeoff optimization for
// the kind of use cases this function was designed for: parsing languages.
// Configuration files and programming languages source code are mainly composed
// of ASCII data, so allocating len(data)*4 seems to be wasteful.
func (rr *ReaderReader) Read(data []rune) (int, error) {
	var (
		nrunes int
		start  int
	)

	end := len(data)
	b := make([]byte, 0, end)

	for nrunes < len(data) {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
			end = len(b) + 1
		}

		begin := len(b)
		n, err := rr.r.Read(b[begin:end])
		b = b[:len(b)+n]

		if n > 0 {
			count := lastPartialCount(b[start : begin+n])
			for start+count < begin+n {
				r, size := utf8.DecodeRune(b[start : begin+n])
				if r == utf8.RuneError {
					return nrunes, fmt.Errorf("invalid rune")
				}

				start += size
				data[nrunes] = r
				nrunes++
			}

			if start+count == begin+n && len(b) == end {
				end++
			}
		}

		if err != nil {
			return nrunes, err
		}
	}

	return nrunes, nil
}

// lastPartialCount count how many bytes of partial runes are in the end of the
// input slice. In the worst case it iterates a maximum of 4 times (utf8.UTF8Max).
func lastPartialCount(p []byte) int {
	// Look for final start of rune.
	for i := 0; i < len(p) && i < utf8.UTFMax; {
		i++
		if utf8.RuneStart(p[len(p)-i]) {
			if utf8.FullRune(p[len(p)-i:]) {
				return 0
			}
			// Found i bytes of partial rune.
			return i
		}
	}
	// Did not find start of final rune - invalid or empty input.
	return 0
}

// NewReader implements a utf8.Reader interface from a io.RuneReader interface.
func NewReader(r io.RuneReader) *RuneReader {
	return &RuneReader{
		r: r,
	}
}

// Read implements the utf8.Reader interface.
func (rd *RuneReader) Read(data []rune) (int, error) {
	for i := 0; i < len(data); i++ {
		r, _, err := rd.r.ReadRune()
		if err != nil {
			return i, err
		}

		if r == utf8.RuneError {
			return i, fmt.Errorf("invalid rune")
		}

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

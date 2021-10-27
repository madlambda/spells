package utf8

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/madlambda/spells/io/bytes"
)

type (
	RuneReader struct {
		r io.RuneReader
	}

	// ReaderReader implements a runes.Reader from a bytes.Reader interface.
	ReaderReader struct {
		r bytes.Reader
	}
)

// NewReaderReader creates a runes.Reader implementation from an bytes.Reader
// interface.
func NewReaderReader(r bytes.Reader) *ReaderReader {
	return &ReaderReader{
		r: r,
	}
}

// Read implements runes.Reader interface.
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
		nrunes    int
		runeStart int // decoding starts at this offset
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

		lastRead := begin + n

		if n > 0 {
			count := lastPartialRuneCount(b[runeStart:lastRead])

			// decode all runes until the partial in the end.
			for runeStart+count < lastRead {
				r, size := utf8.DecodeRune(b[runeStart:lastRead])
				if r == utf8.RuneError {
					return nrunes, fmt.Errorf("invalid rune")
				}

				runeStart += size
				data[nrunes] = r
				nrunes++
			}

			if len(b) == end {
				end++
			}
		}

		if err != nil {
			return nrunes, err
		}
	}

	return nrunes, nil
}

// lastPartialRuneCount count how many bytes of a partial rune are in the end of
// the input slice. In the worst case it iterates a maximum of 4 times
// (utf8.UTF8Max).
func lastPartialRuneCount(p []byte) int {
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

// NewReader implements a runes.Reader interface from a io.RuneReader interface.
func NewReader(r io.RuneReader) *RuneReader {
	return &RuneReader{
		r: r,
	}
}

// Read implements the runes.Reader interface.
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

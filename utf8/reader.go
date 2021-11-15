package utf8

import (
	"fmt"
	"unicode/utf8"

	"github.com/madlambda/spells/io/bytes"
)

// Decoder decodes UTF-8 encoded bytes in a `bytes.Reader` stream.
type Decoder struct {
	r bytes.Reader
}

// Error represents UTF-8 decoding error.
type Error struct {
	offset    int
	undecoded []byte
}

// NewDecoder creates an UTF-8 decoder from a `bytes.Reader` byte stream.
// By calling the Read method, the decoder reads and decodes runes from the byte
// stream. It's an unbuffered decoder, which means every data read from r is
// returned either decoded or in the `Error.Bytes()`, no state is kept in the
// decoder for future calls.
//
// Example:
//   dec := utf8.NewDecoder(bytestream)
//   data := make([]rune, 512)
//   n, err := dec.Read(data[:])
//   if err != nil && err != io.EOF {
//	     undecodedBytes := utf8.ErrBytes(err)
//       ...
//   }
func NewDecoder(r bytes.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Read reads the underlying byte stream and decodes it to runes.
// For details about how to use Read(), please see the documentation of the
// runes.Reader interface.
//
// It will use O(N) space where N is the size of the data slice provided. For
// ASCII data it allocates exactly len(data) bytes (or sizeof(data)/4) but in
// case of multi-byte code points it allocates a maximum of sizeof(data) or
// len(data)*4.
//
// The first call to the source reader uses a full buffer of len(data) bytes but
// subsequent calls read byte by byte. This is a tradeoff optimization for the
// kind of use cases this function was designed for: parsing languages.
// Configuration files and programming languages source code are mainly composed
// of ASCII data, so allocating len(data)*4 seems to be wasteful.
func (d *Decoder) Read(data []rune) (int, error) {
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
		n, err := d.r.Read(b[begin:end])
		lastRead := begin + n
		b = b[:lastRead]

		if n > 0 {
			count := lastPartialRuneCount(b[runeStart:lastRead])

			// decode all runes until the partial in the end.
			for runeStart+count < lastRead {
				r, size := utf8.DecodeRune(b[runeStart:lastRead])
				if r == utf8.RuneError {
					return nrunes, &Error{
						offset:    runeStart,
						undecoded: b[runeStart:],
					}
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

// Offset returns the offset of the offended byte.
func (e *Error) Offset() int { return e.offset }

// Undecoded returns the undecoded bytes read from the byte stream.
// The first byte is the offended invalid rune.
func (e *Error) Undecoded() []byte { return e.undecoded }

func (e *Error) Error() string {
	return fmt.Sprintf("invalid rune at offset %d", e.offset)
}

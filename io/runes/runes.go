package runes

import "io"

// Reader reads Unicode encoded bytes into data. It has the same semantics as
// io.Reader.
type Reader interface {
	Read(data []rune) (int, error)
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

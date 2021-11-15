package utf8_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	stdiotest "testing/iotest"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/io/runes"
	"github.com/madlambda/spells/iotest"
	"github.com/madlambda/spells/utf8"
)

type decodeError struct {
	offset    int
	undecoded []byte
}

type testResult struct {
	err     *decodeError
	partial []rune
}

type testcase struct {
	name   string
	input  string
	repeat int
	want   testResult
}

const socraticParadox = "Εν οίδα ότι ουδέν οίδα"

var testcases = []testcase{
	{
		name:   "empty input",
		input:  "",
		repeat: 1,
	},
	{
		name:   "invalid code point",
		input:  string([]byte{0x80}),
		repeat: 1,
		want: testResult{
			err: &decodeError{
				offset:    0,
				undecoded: []byte{0x80},
			},
		},
	},
	{
		name:   "simple ascii",
		input:  "test",
		repeat: 1,
	},
	{
		name:   "single rune - 2 bytes",
		input:  "ί",
		repeat: 1,
	},
	{
		name:   "1 ascii - start with multibyte",
		input:  "ίA",
		repeat: 1,
	},
	{
		name:   "1 ascii - end with multibyte",
		input:  "Aί",
		repeat: 1,
	},
	{
		name:   "single rune - 3 bytes",
		input:  "ࠀ",
		repeat: 1,
	},
	{
		name:   "decoding ࠀࠆࠉࠌ",
		input:  "ࠀࠆࠉࠌ",
		repeat: 1,
	},
	{
		name:   "decoding mixed 3-byte and ASCII ࠀࠆࠉASCIIࠌ",
		input:  "ࠀࠆࠉASCIIࠌ",
		repeat: 1,
	},
	{
		name:   "single rune - 4 bytes",
		input:  "𒁘",
		repeat: 1,
	},
	{
		name:   "decoding 𓁺𒂞𒆙𒈙𒌦𓁙",
		input:  "𓁺𒂞𒆙𒈙𒌦𓁙",
		repeat: 1,
	},
	{
		name:   "decoding mixed ASCII and 4-byte - 𓁺𒂞𒆙ASCII𒈙𒌦𓁙",
		input:  "𓁺𒂞𒆙ASCII𒈙𒌦𓁙",
		repeat: 1,
	},
	{
		name:   "decoding " + socraticParadox,
		input:  socraticParadox,
		repeat: 1,
	},
	{
		name:   "decoding multibyte code points + ascii",
		input:  "Εν οίδα ότι TESTE ουδέν οίδα",
		repeat: 1,
	},
	{
		name:   "ascii with newline",
		input:  "test\n",
		repeat: 1024,
	},
	{
		name:   "multibyte - 1024 bytes",
		input:  socraticParadox,
		repeat: 1024,
	},
	{
		name:   "multibyte - 1025 bytes",
		input:  socraticParadox,
		repeat: 1025,
	},
	{
		name:   "multibyte - 2048 bytes",
		input:  socraticParadox,
		repeat: 2048,
	},
	{
		name:   "multibyte - 2049 bytes",
		input:  socraticParadox,
		repeat: 2049,
	},
	{
		name:   "invalid rune sequence",
		input:  string([]byte{206, 149, 206, 206, 206}),
		repeat: 1,
		want: testResult{
			err: &decodeError{
				offset:    2,
				undecoded: []byte{0xce, 0xce, 0xce},
			},
			partial: []rune{'Ε'},
		},
	},
}

func TestUTF8Decoder(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			reader := utf8.NewDecoder(iotest.NewRepeatReader(
				strings.NewReader(tc.input),
				tc.repeat))

			got, err := runes.ReadAll(reader)

			var utf8err *utf8.Error

			if err != nil && !errors.As(err, &utf8err) {
				t.Fatalf("got unexpected error: %v", err)
			}

			var expected []rune

			if tc.want.err == nil {
				repeater := iotest.NewRepeatReader(strings.NewReader(tc.input),
					tc.repeat)

				expectedBytes, err := io.ReadAll(repeater)
				assert.NoError(t, err, "repeating expected bytes")

				expected = []rune(string(expectedBytes))
			} else {
				assert.EqualInts(t, tc.want.err.offset, utf8err.Offset(),
					"error offset mismatch")

				wantUndecoded := tc.want.err.undecoded
				gotUndecoded := utf8err.Undecoded()
				assert.EqualInts(t, len(wantUndecoded), len(gotUndecoded),
					"length of undecoded data mismatch")

				for i := 0; i < len(wantUndecoded); i++ {
					if wantUndecoded[i] != gotUndecoded[i] {
						t.Fatalf("undecoded byte %x != %x ", gotUndecoded[i],
							wantUndecoded[i])
					}
				}

				expected = tc.want.partial
			}

			assertRunesEqual(t, expected, got)
		})
	}

	// TODO(i4k): implement stdiotest.TestReader for utf8.Decoder
}

func TestUTF8ReaderMultipleSizes(t *testing.T) {
	input := iotest.NewRepeatReader(strings.NewReader(socraticParadox), 100)
	inputBytes, err := io.ReadAll(input)
	assert.NoError(t, err, "failed to repeat input")

	for i := 0; i < 10; i++ {
		buf1 := bytes.NewBuffer(inputBytes)
		buf2 := bytes.NewBuffer(inputBytes)
		repeater1 := iotest.NewRepeatReader(buf1, i)
		repeater2 := iotest.NewRepeatReader(buf2, i)

		runes, err := runes.ReadAll(utf8.NewDecoder(repeater1))
		assert.NoError(t, err, "reading runes")

		expectedBytes, err := io.ReadAll(repeater2)
		assert.NoError(t, err, "reading expected runes")

		expected := []rune(string(expectedBytes))
		assertRunesEqual(t, expected, runes)
	}
}

func TestUTF8ReaderHalfRead(t *testing.T) {
	input := iotest.NewRepeatReader(strings.NewReader(socraticParadox), 100)
	inputBytes, err := io.ReadAll(input)
	assert.NoError(t, err, "failed to repeat input")

	for i := 0; i < 10; i++ {
		buf1 := bytes.NewBuffer(inputBytes)
		buf2 := bytes.NewBuffer(inputBytes)
		repeater1 := stdiotest.HalfReader(iotest.NewRepeatReader(buf1, i))
		repeater2 := iotest.NewRepeatReader(buf2, i)

		runes, err := runes.ReadAll(utf8.NewDecoder(repeater1))
		assert.NoError(t, err, "reading runes")

		expectedBytes, err := io.ReadAll(repeater2)
		assert.NoError(t, err, "reading expected runes")

		expected := []rune(string(expectedBytes))
		assertRunesEqual(t, expected, runes)
	}
}

func TestUTF8ReaderFromFileMultipleBufferSizes(t *testing.T) {
	expectedRunes := []rune(socraticParadox)
	expectedLen := len(expectedRunes)

	temp, err := ioutil.TempFile("", "spells-utf8")
	assert.NoError(t, err, "creating temp file")

	defer os.Remove(temp.Name())

	n, err := temp.WriteString(socraticParadox)
	assert.NoError(t, err, "writing utf8 file")
	assert.EqualInts(t, len(socraticParadox), n, "written len mismatch")

	temp.Close()

	for i := 0; i < expectedLen; i++ {
		t.Run(fmt.Sprintf("using buffer size: %d", i), func(t *testing.T) {
			temp, err := os.Open(temp.Name())
			assert.NoError(t, err, "open file for read")

			defer temp.Close()

			r := utf8.NewDecoder(temp)

			runes := make([]rune, i)
			n, err = r.Read(runes)

			assert.NoError(t, err, "failed reading runes")
			assert.EqualInts(t, i, n, "read wrong number of runes")
			assert.EqualStrings(t, string(expectedRunes[0:i]), string(runes[0:n]),
				"wrong read utf8 string")
		})
	}
}

func TestUTF8ReaderNonEOF(t *testing.T) {
	buf := stdiotest.DataErrReader(strings.NewReader("test"))
	reader := utf8.NewDecoder(buf)

	data := make([]rune, 10)
	n, err := reader.Read(data)
	assert.EqualInts(t, 4, n, "read size mismatch")
	assert.EqualErrs(t, io.EOF, err, "error is not EOF")
	assert.EqualStrings(t, "test", string(data[:n]), "string mismatch")
}

func TestUTF8ReaderError(t *testing.T) {
	expected := fmt.Errorf("some error")
	buf := stdiotest.ErrReader(expected)
	reader := utf8.NewDecoder(buf)

	data := make([]rune, 10)
	n, err := reader.Read(data)
	assert.EqualErrs(t, expected, err, "err mismatch")
	assert.EqualInts(t, 0, n, "read size mismatch")
}

func TestUTF8ReaderOneByteReader(t *testing.T) {
	buf1 := strings.NewReader(socraticParadox)
	buf2 := strings.NewReader(socraticParadox)
	repeater1 := stdiotest.HalfReader(iotest.NewRepeatReader(buf1, 100))
	repeater2 := iotest.NewRepeatReader(buf2, 100)

	runes, err := runes.ReadAll(utf8.NewDecoder(repeater1))
	assert.NoError(t, err, "reading runes")

	expectedBytes, err := io.ReadAll(repeater2)
	assert.NoError(t, err, "reading expected runes")

	expected := []rune(string(expectedBytes))

	assertRunesEqual(t, expected, runes)
}

func assertRunesEqual(t *testing.T, expected, got []rune) {
	assert.EqualInts(t, len(expected), len(got), "length mismatch")

	for i, r := range expected {
		if r != got[i] {
			t.Errorf("want[%d = %c] but got[%d = %c]", r, r, got[i], got[i])
		}
	}
}

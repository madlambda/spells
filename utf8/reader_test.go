package utf8_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	stdiotest "testing/iotest"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/io/runes"
	"github.com/madlambda/spells/iotest"
	"github.com/madlambda/spells/utf8"
)

const expectedString = "Εν οίδα ότι ουδέν οίδα"

func TestUTF8Reader(t *testing.T) {
	type testcase struct {
		input  string
		repeat int
		err    error
		errOut []rune
	}

	testcases := []testcase{
		{
			input:  "",
			repeat: 1,
		},
		{
			input:  string([]byte{0x80}),
			repeat: 1,
			err:    fmt.Errorf("invalid rune"),
		},
		{
			input:  "test",
			repeat: 1,
		},
		{
			input:  "Εν οίδα ότι ουδέν οίδα",
			repeat: 1,
		},
		{
			input:  "Εν οίδα ότι TESTE ουδέν οίδα",
			repeat: 1,
		},
		{
			input:  "test\n",
			repeat: 1024,
		},
		{
			input:  "Εν οίδα ότι ουδέν οίδα",
			repeat: 1024,
		},
		{
			input:  "Εν οίδα ότι ουδέν οίδα",
			repeat: 1025,
		},
		{
			input:  "Εν οίδα ότι ουδέν οίδα",
			repeat: 2048,
		},
		{
			input:  "Εν οίδα ότι ουδέν οίδα",
			repeat: 2049,
		},
		{
			input:  string([]byte{206, 149, 206, 206, 206}),
			repeat: 1,
			err:    fmt.Errorf("invalid rune"),
			errOut: []rune{'Ε'},
		},
	}

	// test NewReaderReader
	for _, tc := range testcases {
		reader := utf8.NewReaderReader(iotest.NewRepeatReader(
			bytes.NewBuffer([]byte(tc.input)),
			tc.repeat))

		got, readErr := runes.ReadAll(reader)
		assert.EqualErrs(t, tc.err, readErr, "Read() error")

		var expected []rune

		if readErr == nil {
			repeater := iotest.NewRepeatReader(bytes.NewBuffer([]byte(tc.input)),
				tc.repeat)

			expectedBytes, err := io.ReadAll(repeater)
			assert.NoError(t, err, "repeating expected")

			expected = []rune(string(expectedBytes))
		} else {
			expected = tc.errOut
		}

		assert.EqualInts(t, len(expected), len(got), "rune slice len mismatch: %s", string(got))

		for i, r := range expected {
			if r != got[i] {
				t.Errorf("want[%d = %c] but got[%d = %c]", r, r, got[i], got[i])
			}
		}
	}

	// test NewReader
	for _, tc := range testcases {
		inputReader := iotest.NewRepeatReader(
			bytes.NewBuffer([]byte(tc.input)),
			tc.repeat)

		data, err := io.ReadAll(inputReader)
		assert.NoError(t, err, "reading repeated string")

		reader := utf8.NewReader(bytes.NewBuffer(data))

		got, readErr := runes.ReadAll(reader)
		assert.EqualErrs(t, tc.err, readErr, "read() error")

		var expected []rune

		if readErr == nil {
			repeater := iotest.NewRepeatReader(bytes.NewBuffer([]byte(tc.input)),
				tc.repeat)

			expectedBytes, err := io.ReadAll(repeater)
			assert.NoError(t, err, "repeating expected")

			expected = []rune(string(expectedBytes))
		} else {
			expected = tc.errOut
		}

		assert.EqualInts(t, len(expected), len(got), "rune slice len mismatch")

		for i, r := range expected {
			if r != got[i] {
				t.Errorf("want[%d = %c] but got[%d = %c]", r, r, got[i], got[i])
			}
		}
	}

	// TODO(i4k): implement stdiotest.TestReader for utf8.Reader
}

func TestUTF8ReaderMultipleSizes(t *testing.T) {
	input := iotest.NewRepeatReader(bytes.NewBuffer([]byte(expectedString)), 100)
	inputBytes, err := io.ReadAll(input)
	assert.NoError(t, err, "failed to repeat input")

	for i := 0; i < 10; i++ {
		buf1 := bytes.NewBuffer(inputBytes)
		buf2 := bytes.NewBuffer(inputBytes)
		repeater1 := iotest.NewRepeatReader(buf1, i)
		repeater2 := iotest.NewRepeatReader(buf2, i)

		runes, err := runes.ReadAll(utf8.NewReaderReader(repeater1))
		assert.NoError(t, err, "reading runes")

		expectedBytes, err := io.ReadAll(repeater2)
		assert.NoError(t, err, "reading expected runes")

		expected := []rune(string(expectedBytes))

		assert.EqualInts(t, len(expected), len(runes), "length mismatch")

		for j := 0; j < len(expected); j++ {
			if expected[j] != runes[j] {
				t.Errorf("rune %d are not equal: %c != %c",
					j, expected[j], runes[j])
			}
		}
	}
}

func TestUTF8ReaderHalfRead(t *testing.T) {
	input := iotest.NewRepeatReader(bytes.NewBuffer([]byte(expectedString)), 100)
	inputBytes, err := io.ReadAll(input)
	assert.NoError(t, err, "failed to repeat input")

	for i := 0; i < 10; i++ {
		buf1 := bytes.NewBuffer(inputBytes)
		buf2 := bytes.NewBuffer(inputBytes)
		repeater1 := stdiotest.HalfReader(iotest.NewRepeatReader(buf1, i))
		repeater2 := iotest.NewRepeatReader(buf2, i)

		runes, err := runes.ReadAll(utf8.NewReaderReader(repeater1))
		assert.NoError(t, err, "reading runes")

		expectedBytes, err := io.ReadAll(repeater2)
		assert.NoError(t, err, "reading expected runes")

		expected := []rune(string(expectedBytes))

		assert.EqualInts(t, len(expected), len(runes), "length mismatch")

		for j := 0; j < len(expected); j++ {
			if expected[j] != runes[j] {
				t.Errorf("rune %d are not equal: %c != %c",
					j, expected[j], runes[j])
			}
		}
	}
}

func TestUTF8ReaderFromFileMultipleBufferSizes(t *testing.T) {
	expectedRunes := []rune(expectedString)
	expectedLen := len(expectedRunes)

	temp, err := ioutil.TempFile("", "spells-utf8")
	assert.NoError(t, err, "creating temp file")

	defer os.Remove(temp.Name())

	n, err := temp.WriteString(expectedString)
	assert.NoError(t, err, "writing utf8 file")
	assert.EqualInts(t, len(expectedString), n, "written len mismatch")

	for i := 0; i < expectedLen; i++ {
		t.Run(fmt.Sprintf("using buffer size: %d", i), func(t *testing.T) {
			off, err := temp.Seek(0, 0)
			assert.NoError(t, err, "seek to file offset 0")
			assert.EqualInts(t, 0, int(off), "invalid offset")

			r := utf8.NewReaderReader(temp)

			runes := make([]rune, i)
			n, err = r.Read(runes[:])

			assert.NoError(t, err, "failed reading runes")
			assert.EqualInts(t, i, n, "read wrong number of runes")
			assert.EqualStrings(t, string(expectedRunes[0:i]), string(runes[0:n]),
				"wrong read utf8 string")
		})
	}
}

func TestUTF8ReaderNonEOF(t *testing.T) {
	buf := stdiotest.DataErrReader(bytes.NewBuffer([]byte("test")))
	reader := utf8.NewReaderReader(buf)

	data := make([]rune, 10)
	n, err := reader.Read(data[:])
	assert.EqualInts(t, 4, n, "read size mismatch")
	assert.EqualErrs(t, io.EOF, err, "error is not EOF")
	assert.EqualStrings(t, "test", string(data[:n]), "string mismatch")
}

func TestUTF8ReaderError(t *testing.T) {
	expected := fmt.Errorf("some error")
	buf := stdiotest.ErrReader(expected)
	reader := utf8.NewReaderReader(buf)

	data := make([]rune, 10)
	n, err := reader.Read(data[:])
	assert.EqualErrs(t, expected, err, "err mismatch")
	assert.EqualInts(t, 0, n, "read size mismatch")
}

func TestUTF8ReaderOneByteReader(t *testing.T) {
	buf1 := bytes.NewBuffer([]byte(expectedString))
	buf2 := bytes.NewBuffer([]byte(expectedString))
	repeater1 := stdiotest.HalfReader(iotest.NewRepeatReader(buf1, 100))
	repeater2 := iotest.NewRepeatReader(buf2, 100)

	runes, err := runes.ReadAll(utf8.NewReaderReader(repeater1))
	assert.NoError(t, err, "reading runes")

	expectedBytes, err := io.ReadAll(repeater2)
	assert.NoError(t, err, "reading expected runes")

	expected := []rune(string(expectedBytes))

	assert.EqualInts(t, len(expected), len(runes), "length mismatch")

	for j := 0; j < len(expected); j++ {
		if expected[j] != runes[j] {
			t.Errorf("rune %d are not equal: %c != %c",
				j, expected[j], runes[j])
		}
	}
}

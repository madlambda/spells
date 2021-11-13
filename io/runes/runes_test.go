package runes_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/io/runes"
	"github.com/madlambda/spells/iotest"
)

type testResult struct {
	err     error
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
			err: fmt.Errorf("invalid rune at offset 0"),
		},
	},
	{
		name:   "simple ascii",
		input:  "test",
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
			err:     fmt.Errorf("invalid rune at offset 2"),
			partial: []rune{'Ε'},
		},
	},
}

func TestUnicodeDecoder(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := iotest.NewRepeatReader(
				bytes.NewBuffer([]byte(tc.input)),
				tc.repeat)

			data, err := io.ReadAll(inputReader)
			assert.NoError(t, err, "reading repeated string")

			reader := runes.NewUnicodeReader(bytes.NewBuffer(data))

			got, err := runes.ReadAll(reader)
			assert.EqualErrs(t, tc.want.err, err, "readAll error mismatch")

			var expected []rune

			if tc.want.err == nil {
				repeater := iotest.NewRepeatReader(bytes.NewBuffer([]byte(tc.input)),
					tc.repeat)

				expectedBytes, err := io.ReadAll(repeater)
				assert.NoError(t, err, "repeating expected")

				expected = []rune(string(expectedBytes))
			} else {
				expected = tc.want.partial
			}

			assert.EqualInts(t, len(expected), len(got), "rune slice len mismatch")

			for i, r := range expected {
				if r != got[i] {
					t.Errorf("want[%d = %c] but got[%d = %c]", r, r, got[i], got[i])
				}
			}
		})
	}
}

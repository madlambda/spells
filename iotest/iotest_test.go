package iotest_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	stdiotest "testing/iotest"

	"github.com/madlambda/spells/iotest"
)

// TODO: Test behavior with negative indexes
// Behavior could be:
// - Just once, as 0
// - Repeat forever...which may be useful and is not unheard off: https://pkg.go.dev/strings#SplitAfterN
// - Panic, like: https://pkg.go.dev/strings#Repeat
//
// Overall, maybe it is best to keep same behavior as strings.Repeat (0 = none, 1 = repeat once).

func TestRepeatReader(t *testing.T) {

	type Test struct {
		name   string
		data   []byte
		repeat int
		want   []byte
	}

	tests := []Test{
		{
			name:   "NoRepeat",
			data:   []byte("test"),
			repeat: 0,
			want:   []byte("test"),
		},
		{
			name:   "RepeatOnce",
			data:   []byte("test"),
			repeat: 1,
			want:   []byte("testtest"),
		},
		{
			name:   "RepeatTwice",
			data:   []byte("test"),
			repeat: 2,
			want:   []byte("testtesttest"),
		},
		{
			name:   "RepeatSingleByte",
			data:   []byte("t"),
			repeat: 1,
			want:   []byte("tt"),
		},
		{
			name:   "RepeatEmpty",
			data:   []byte{},
			repeat: 1,
			want:   []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := bytes.NewBuffer(test.data)
			repeater := iotest.NewRepeatReader(input, test.repeat)

			got, err := ioutil.ReadAll(repeater)
			if err != nil {
				t.Fatal(err)
			}

			assertEqualText(t, got, test.want)
		})
	}
}

func TestRepeatReaderReadBytePerByte(t *testing.T) {
	inputData := []byte("bytePerByte")
	want := append(inputData, inputData...)
	repeater := iotest.NewRepeatReader(bytes.NewBuffer(inputData), 1)

	got, err := ioutil.ReadAll(stdiotest.OneByteReader(repeater))
	if err != nil {
		t.Fatal(err)
	}
	assertEqualText(t, got, want)
}

func TestRepeatReaderCornerCasesOnUnderlyingReader(t *testing.T) {
	scenarios := map[string]func(io.Reader) io.Reader{
		"SingleByteReading": stdiotest.OneByteReader,
		"EOFWithData":       stdiotest.DataErrReader,
	}

	for testname, newReader := range scenarios {
		t.Run(testname, func(t *testing.T) {
			inputData := []byte("cornercases")
			want := append(inputData, inputData...)
			repeater := iotest.NewRepeatReader(newReader(bytes.NewBuffer(inputData)), 1)

			got, err := ioutil.ReadAll(repeater)
			if err != nil {
				t.Fatal(err)
			}
			assertEqualText(t, got, want)
		})
	}
}

func TestRepeatReaderNonEOFErr(t *testing.T) {
	want := errors.New("TestRepeatReaderNonEOFErr")
	repeater := iotest.NewRepeatReader(stdiotest.ErrReader(want), 666)
	data := make([]byte, 10)

	for i := 0; i < 10; i++ {
		n, err := repeater.Read(data)
		if n != 0 {
			t.Errorf("got n=%d; want=0", n)
		}
		if err != want {
			t.Errorf("got %v; want %v", err, want)
		}
	}
}

func assertEqualText(t *testing.T, got []byte, want []byte) {
	t.Helper()

	// WHY: we usually test with text
	gotstr := string(got)
	wantstr := string(want)
	if gotstr != wantstr {
		t.Errorf("got %q; want %q", gotstr, wantstr)
	}
}

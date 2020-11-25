package iotest_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/madlambda/spells/iotest"
)

func TestRepeatReader(t *testing.T) {

	type Test struct {
		name   string
		data   []byte
		repeat uint
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := bytes.NewBuffer(test.data)
			repeater := iotest.NewRepeater(input, test.repeat)

			got, err := ioutil.ReadAll(repeater)
			if err != nil {
				t.Fatal(err)
			}

			// WHY: we usually test with text
			gotstr := string(got)
			wantstr := string(test.want)

			if gotstr != wantstr {
				t.Fatalf("got %q; want %q", gotstr, wantstr)
			}
		})
	}
}

func TestRepeatReaderNonEOFErr(t *testing.T) {
	want := errors.New("TestRepeatReaderNonEOFErr")
	repeater := iotest.NewRepeater(iotest.BrokenReader{Err: want}, 666)
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

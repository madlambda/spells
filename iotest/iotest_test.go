package iotest_test

import (
	"bytes"
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
	repeater := iotest.NewRepeater(iotest.BrokenReader{}, 666)
	data := make([]byte, 10)

	n, err := repeater.Read(data)
	if n != 0 {
		t.Errorf("got n=%d; want=0", n)
	}
	if err == nil {
		t.Error("got nil error; want non-nil error")
	}
	n, err2 := repeater.Read(data)
	if n != 0 {
		t.Errorf("got n=%d; want=0", n)
	}
	if err != err2 {
		t.Errorf("got error: %v; want error: %v", err2, err)
	}
}

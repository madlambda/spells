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

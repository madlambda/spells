package iotest_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/madlambda/spells/iotest"
)

func TestRepeatReader(t *testing.T) {

	type Test struct {
		data   []byte
		repeat uint
		want   []byte
	}

	tests := []Test{
		{
			data:   []byte("test"),
			repeat: 1,
			want:   []byte("test"),
		},
	}

	for _, test := range tests {
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
	}
}

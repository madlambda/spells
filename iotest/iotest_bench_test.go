package iotest_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/madlambda/spells/iotest"
)

// The core idea here is to repeat the same input
// different times and check that memory usage is constant (O(1))
// regarding N where N is amount of repetitions.

func BenchmarkRepeat1(b *testing.B) {
	benchmarkRepeatReader(b, newInput(), 1)
}

func BenchmarkRepeat10(b *testing.B) {
	benchmarkRepeatReader(b, newInput(), 10)
}

func BenchmarkRepeat100(b *testing.B) {
	benchmarkRepeatReader(b, newInput(), 100)
}

func BenchmarkRepeat1000(b *testing.B) {
	benchmarkRepeatReader(b, newInput(), 1000)
}

func benchmarkRepeatReader(b *testing.B, input io.Reader, repeatCount int) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repeater := iotest.NewRepeatReader(input, repeatCount)
		_, err := ioutil.ReadAll(repeater)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func newInput() *bytes.Buffer {
	return bytes.NewBuffer([]byte("benchmarking is cool"))
}

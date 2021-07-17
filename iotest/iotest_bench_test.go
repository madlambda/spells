package iotest_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/madlambda/spells/iotest"
)

// The core idea here is to repeat the same input
// different times and check that memory usage is constant (O(1))
// regarding N (for N = amount of repetitions).

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

func BenchmarkRepeat10000(b *testing.B) {
	benchmarkRepeatReader(b, newInput(), 10000)
}

func benchmarkRepeatReader(b *testing.B, input io.Reader, repeatCount int) {
	data := make([]byte, 10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		repeater := iotest.NewRepeatReader(input, repeatCount)
		var err error
		for err == nil {
			_, err = repeater.Read(data)
		}
		if err != io.EOF {
			b.Fatal(err)
		}
	}
}

func newInput() *bytes.Buffer {
	return bytes.NewBuffer([]byte("benchmarking is cool"))
}

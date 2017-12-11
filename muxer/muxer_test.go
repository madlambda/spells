package muxer_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/muxer"
)

func TestMux(t *testing.T) {
	for _, tcase := range []TestCase{
		TestCase{
			name:            "oneInputOneSource",
			expectedOutputs: []int{666},
			sourceChannels:  1,
		},
		TestCase{
			name:            "multipleInputsOneSource",
			expectedOutputs: []int{666, 777, 10, 0, 1},
			sourceChannels:  1,
		},
		TestCase{
			name:            "sameInputsAsSources",
			expectedOutputs: []int{666, 777, 10},
			sourceChannels:  3,
		},
		TestCase{
			name:            "lessInputsThanSources",
			expectedOutputs: []int{666, 777},
			sourceChannels:  3,
		},
		TestCase{
			name:            "moreInputsThanSources",
			expectedOutputs: []int{666, 777, 234},
			sourceChannels:  2,
		},
	} {
		testMux(t, tcase)
	}
}

func TestMuxClosedChannels(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)

	close(source1)
	close(source2)

	assert.NoError(t, muxer.Do(sink, source1, source2))

	v, ok := <-sink
	if ok {
		t.Fatalf("expected sink to be closed, instead got val[%d]", v)
	}
}

func TestErrorOnInvalidOutput(t *testing.T) {
}

func TestErrorOnInvalidInput(t *testing.T) {
}

func TestErrorOnIncompatibleInputsOutputs(t *testing.T) {
}

type TestCase struct {
	name            string
	expectedOutputs []int
	sourceChannels  int
}

func testMux(t *testing.T, tcase TestCase) {
	t.Run(tcase.name, func(t *testing.T) {
		sources := []chan int{}
		sourcesgen := []interface{}{}
		for i := 0; i < tcase.sourceChannels; i++ {
			source := make(chan int)
			sources = append(sources, source)
			sourcesgen = append(sourcesgen, source)
		}
		sink := make(chan int)
		assert.NoError(t, muxer.Do(sink, sourcesgen...))

		go func() {
			for i, v := range tcase.expectedOutputs {
				inindex := i % len(sources)
				sources[inindex] <- v
			}

			for _, source := range sources {
				close(source)
			}
		}()

		for _, want := range tcase.expectedOutputs {
			got := <-sink
			assert.EqualInts(t, want, got)
		}

		v, ok := <-sink
		if ok {
			t.Fatalf("expected sink to be closed, got val[%d]", v)
		}
	})
}

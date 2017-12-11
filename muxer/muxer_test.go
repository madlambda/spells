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
			expectedOutputs: []int{666, 777, 234, 1, 0},
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

func TestMuxDirectionedChannels(t *testing.T) {
	sink := make(chan string)
	source := make(chan string)

	var sinkd chan<- string = sink
	var sourced <-chan string = source

	const expectedVal string = "lambda"

	go func() {
		source <- expectedVal
		close(source)
	}()

	assert.NoError(t, muxer.Do(sinkd, sourced))

	v, _ := <-sink
	assert.EqualStrings(t, expectedVal, v)
}

func TestFailsOnWrongSourceDirection(t *testing.T) {
	// the reflect package validates like this =(
	//tt := (*chanType)(unsafe.Pointer(ch.typ))
	//if ChanDir(tt.dir)&RecvDir == 0 {
	//panic("reflect.Select: RecvDir case using send-only channel")
	//}
	t.Skip("no good way to verify direction =(")
	sink := make(chan string)
	source := make(chan string)

	var sourced chan<- string = source
	assert.Error(t, muxer.Do(sink, sourced))
}

func TestFailsOnWrongSinkDirection(t *testing.T) {
	// the reflect package validates like this =(
	//tt := (*chanType)(unsafe.Pointer(ch.typ))
	//if ChanDir(tt.dir)&RecvDir == 0 {
	//panic("reflect.Select: RecvDir case using send-only channel")
	//}
	t.Skip("no good way to verify direction =(")
	sink := make(chan string)
	source := make(chan string)

	var sinkd <-chan string = sink
	assert.Error(t, muxer.Do(sinkd, source))
}

func TestErrorOnInvalidSink(t *testing.T) {
	for name, sink := range invalidCases() {
		t.Run(name, func(t *testing.T) {
			source := make(chan int)
			assert.Error(t, muxer.Do(sink, source))
		})
	}
}

func TestErrorOnInvalidSource(t *testing.T) {
	for name, source := range invalidCases() {
		t.Run(name, func(t *testing.T) {
			sink := make(chan int)
			validSource := make(chan int)
			assert.Error(t, muxer.Do(sink, validSource, source))
			assert.Error(t, muxer.Do(sink, source, validSource))
		})
	}
}

type TestCase struct {
	name            string
	expectedOutputs []int
	sourceChannels  int
}

func invalidCases() map[string]interface{} {
	valid := make(chan int)
	var nilChannel chan int

	return map[string]interface{}{
		"nil":              nil,
		"nilChannel":       nilChannel,
		"notChannel":       1,
		"pointerToChannel": &valid,
		"wrongType":        make(chan uint),
	}
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

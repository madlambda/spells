package muxer_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/muxer"
)

func TestMux(t *testing.T) {
	for _, tcase := range []TestCase{
		TestCase{
			name:            "OneInputOneSource",
			expectedOutputs: []int{666},
			sourceChannels:  1,
		},
		TestCase{
			name:            "MultipleInputsOneSource",
			expectedOutputs: []int{666, 777, 10, 0, 1},
			sourceChannels:  1,
		},
		TestCase{
			name:            "SameInputsAsSources",
			expectedOutputs: []int{666, 777, 10},
			sourceChannels:  3,
		},
		TestCase{
			name:            "LessInputsThanSources",
			expectedOutputs: []int{666, 777},
			sourceChannels:  3,
		},
		TestCase{
			name:            "MoreInputsThanSources",
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

func TestMuxCloseFirstSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, muxer.Do(sink, source1, source2, source3))

	close(source1)

	want1 := 666
	want2 := 777
	go func() {
		source3 <- want1
		source2 <- want2
	}()

	assert.EqualInts(t, want1, <-sink)
	assert.EqualInts(t, want2, <-sink)
}

func TestMuxCloseMiddleSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, muxer.Do(sink, source1, source2, source3))

	close(source2)

	want1 := 100
	want2 := 7
	go func() {
		source3 <- want1
		source1 <- want2
	}()

	assert.EqualInts(t, want1, <-sink)
	assert.EqualInts(t, want2, <-sink)
}

func TestMuxCloseLastSource(t *testing.T) {
	sink := make(chan int)
	source1 := make(chan int)
	source2 := make(chan int)
	source3 := make(chan int)

	assert.NoError(t, muxer.Do(sink, source1, source2, source3))

	close(source3)

	want1 := 100
	want2 := 7
	go func() {
		source2 <- want1
		source1 <- want2
	}()

	assert.EqualInts(t, want1, <-sink)
	assert.EqualInts(t, want2, <-sink)
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

	v := <-sink
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
				srcindex := i % len(sources)
				sources[srcindex] <- v
			}

			for _, source := range sources {
				close(source)
			}
		}()

		gotOutputs := []int{}
		for got := range sink {
			gotOutputs = append(gotOutputs, got)
		}

		if len(gotOutputs) != len(tcase.expectedOutputs) {
			t.Fatalf("got %v != want %v", gotOutputs, tcase.expectedOutputs)
		}

		for i, want := range tcase.expectedOutputs {
			got := gotOutputs[i]
			assert.EqualInts(t, want, got)
		}
	})
}

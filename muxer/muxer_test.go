package muxer_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/muxer"
)

func TestMux(t *testing.T) {
	for _, tcase := range []TestCase{
		TestCase{
			name:            "oneInput",
			expectedOutputs: []int{666},
			inputChannels:   1,
		},
	} {
		testMux(t, tcase)
	}
}

func TestMultipleIntChannels(t *testing.T) {
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
	inputChannels   int
}

func testMux(t *testing.T, tcase TestCase) {
	t.Run(tcase.name, func(t *testing.T) {
		inputs := []chan int{}
		inputsgen := []interface{}{}
		for i := 0; i < tcase.inputChannels; i++ {
			input := make(chan int)
			inputs = append(inputs, input)
			inputsgen = append(inputsgen, input)
		}
		output := make(chan int)
		assert.NoError(t, muxer.Do(output, inputsgen...))

		go func() {
			for i, v := range tcase.expectedOutputs {
				inindex := i % len(inputs)
				inputs[inindex] <- v
			}
		}()

		for _, want := range tcase.expectedOutputs {
			got := <-output
			assert.EqualInts(t, want, got)
		}
	})
}

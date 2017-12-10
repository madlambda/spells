package muxer_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/muxer"
)

func TestSingleIntChannel(t *testing.T) {
	output := make(chan int)
	input := make(chan int)
	expectedOutput := 666

	assert.NoError(t, muxer.Do(output, input))
	go func() {
		input <- expectedOutput
	}()
	outVal := <-output
	assert.EqualInts(t, expectedOutput, outVal)
}

func TestMultipleIntChannels(t *testing.T) {
}

func TestErrorOnInvalidOutput(t *testing.T) {
}

func TestErrorOnInvalidInput(t *testing.T) {
}

func TestErrorOnIncompatibleInputsOutputs(t *testing.T) {
}

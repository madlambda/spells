package errutil_test

import (
	"errors"
	"fmt"

	"github.com/madlambda/spells/errutil"
)

func ExampleError() {
	// Declare your error sentinels using errutil.Error
	const (
		someError errutil.Error = "someError"
	)

	// You add some context to the sentinel error
	wrappedErr := fmt.Errorf("wrapping up: %w", someError)

	// Checking programmatically for the underlying error
	// Users of your API handle the sentinel opaquely
	fmt.Println(errors.Is(wrappedErr, someError))

	// Output: true
}

func ExampleChain() {
	// Declare your error sentinels using errutil.Error
	const (
		layer1Err errutil.Error = "layer1Err"
		layer2Err errutil.Error = "layer2Err"
		layer3Err errutil.Error = "layer3Err"
	)

	// Chain the errors
	err := errutil.Chain(layer1Err, layer2Err, layer3Err)

	// Checking programmatically for the underlying error
	// Users of your API handle the sentinels opaquely
	fmt.Println(errors.Is(err, layer1Err))
	fmt.Println(errors.Is(err, layer2Err))
	fmt.Println(errors.Is(err, layer3Err))

	// Output:
	// true
	// true
	// true
}

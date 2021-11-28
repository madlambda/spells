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
	fmt.Println(err)

	// Output:
	// true
	// true
	// true
	// layer1Err: layer2Err: layer3Err
}

func ExampleMerge() {
	// call multiple functions that may return an error
	// but none of them should interrupt overall computation
	var i int
	someFunc := func() error {
		i++
		return fmt.Errorf("error %d", i)
	}

	var errs []error

	errs = append(errs, someFunc())
	errs = append(errs, someFunc())
	errs = append(errs, someFunc())

	// Chain the errors
	err := errutil.Merge(errs...)

	// Checking programmatically for the underlying error
	// Users of your API handle the sentinels opaquely
	fmt.Println(err)

	// Output:
	// error 1: error 2: error 3
}

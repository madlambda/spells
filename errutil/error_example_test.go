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

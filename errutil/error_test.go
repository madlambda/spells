package errutil_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/madlambda/spells/errutil"
)

func TestErrorSentinelWrapping(t *testing.T) {
	const (
		someError        errutil.Error = "someError"
		someAnotherError errutil.Error = "someAnotherError"
	)

	wrappedErr := fmt.Errorf("wrapping up: %w", someError)
	assertErrorIsWrapped(t, wrappedErr, someError)

	wrappedErr2 := fmt.Errorf("wrapping up: %w", someAnotherError)
	assertErrorIsWrapped(t, wrappedErr2, someAnotherError)
}

func TestErrorRepresentation(t *testing.T) {
	const (
		someError errutil.Error = "someError"
	)

	got := someError.Error()
	want := string(someError)

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestErrorChain(t *testing.T) {
	errs := []error{
		errors.New("top error"),
		errors.New("wrapped error 1"),
		errors.New("wrapped error 2"),
		errors.New("wrapped error 3"),
	}

	err := errutil.Chain(errs...)
	assert.Error(t, err)

	got := err
	for i, want := range errs {
		if got == nil {
			t.Fatal("expected error to exist, got nil")
		}

		if !errors.Is(got, want) {
			t.Fatalf("error[%d] got: [%v] want: [%v]", i, got, want)
		}

		// We could only test chain through errors.Is
		// But wanted to check the unwrapping order too.
		got = errors.Unwrap(got)
	}

	if got != nil {
		t.Fatalf("wanted error chain to reach end (nil), got chain [%v] instead", got)
	}
}

func TestErrorChainStringRepresentation(t *testing.T) {
	type TestCase struct {
		name string
		errs []error
		want string
	}

	tcases := []TestCase{
		{
			name: "Single Error",
			errs: []error{
				errors.New("error 1"),
			},
			want: "error 1",
		},
		{
			name: "Two Chained Errors",
			errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
			},
			want: "error 1: error 2",
		},
		{
			name: "Three Chained Errors",
			errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			},
			want: "error 1: error 2: error 3",
		},
	}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {

			err := errutil.Chain(tc.errs...)
			assert.Error(t, err)

			got := err.Error()

			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}

}

func TestErrorChainForEmptyErrList(t *testing.T) {
	assert.NoError(t, errutil.Chain())
	errs := []error{}
	assert.NoError(t, errutil.Chain(errs...))
}

func TestErrorChainRespectIsMethodOfChainedErrors(t *testing.T) {
	var neverIs errorThatNeverIs

	err := errutil.Chain(neverIs)
	if errors.Is(err, neverIs) {
		t.Fatalf("errors.Is(%q, %q) = true, wanted false", err, neverIs)
	}
}

// To test the Is method the error must not be comparable.
// If it is comparable, Go always just compares it, the Is method
// is just a fallback, not an override of actual behavior.
type errorThatNeverIs []string

func (e errorThatNeverIs) Is(err error) bool {
	return false
}

func (e errorThatNeverIs) Error() string {
	return "never is"
}

func assertErrorIsWrapped(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("error [%v] is not wrapping [%v]", err, target)
	}
}

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
	testcases := [][]error{
		[]error{errors.New("single error")},
		[]error{
			errors.New("top error"),
			errors.New("wrapped error 1"),
		},
		[]error{
			errors.New("top error"),
			errors.New("wrapped error 1"),
			errors.New("wrapped error 2"),
		},
	}

	for _, errs := range testcases {

		name := fmt.Sprintf("%dErrors", len(errs))
		t.Run(name, func(t *testing.T) {

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
		})

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

func TestErrorChainTypeSelection(t *testing.T) {
	want1 := error1{
		data: "yay",
	}
	want2 := error2{
		data: 666,
	}

	err := errutil.Chain(want1, want2)
	assert.Error(t, err)

	var got1 error1

	if !errors.As(err, &got1) {
		t.Fatalf("errors.As(%v, %v) == false, want true", err, &got1)
	}
	assert.EqualStrings(t, want1.data, got1.data)

	var got2 error2
	if !errors.As(err, &got2) {
		t.Fatalf("errors.As(%v, %v) == false, want true", err, &got2)
	}
	assert.EqualInts(t, want2.data, got2.data)

	var unrelatedErr errutil.Error
	if errors.As(err, &unrelatedErr) {
		t.Fatalf("errors.As(%v, %v) == true, want false", err, &unrelatedErr)
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

func TestErrorReducing(t *testing.T) {
	type testcase struct {
		name   string
		input  []error
		reduce errutil.Reducer
		want   error
	}

	tests := []testcase{
		{
			name:  "merging two errors",
			input: []error{errors.New("one"), errors.New("two")},
			reduce: func(err1, err2 error) error {
				return fmt.Errorf("%v:%v", err1, err2)
			},
			want: errors.New("one:two"),
		},
		{
			name: "merging three errors",
			input: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return fmt.Errorf("%v/%v", err1, err2)
			},
			want: errors.New("one/two/three"),
		},
		{
			name: "filtering just first err",
			input: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return err1
			},
			want: errors.New("one"),
		},
		{
			name: "filtering just second err",
			input: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return err2
			},
			want: errors.New("three"),
		},
		{
			name: "reduces to nil",
			input: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return nil
			},
			want: nil,
		},
		{
			name:  "reduces empty err list to nil",
			input: []error{},
			reduce: func(err1, err2 error) error {
				panic("unreachable")
				return nil
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := errutil.Reduce(test.reduce, test.input...)

			if test.want == nil {
				if g == nil {
					return
				}
				t.Fatalf(
					"errutil.Reduce(%v)=%q; want nil",
					test.input,
					g,
				)
			}

			if g == nil {
				t.Fatalf(
					"errutil.Reduce(%v)=nil; want %q",
					test.input,
					test.want,
				)
			}

			got := g.Error()
			want := test.want.Error()

			if got != want {
				t.Fatalf(
					"errutil.Reduce(%v)=%q; want=%q",
					test.input,
					got,
					want,
				)
			}
		})
	}
}

func TestErrorMerging(t *testing.T) {
	type TestCase struct {
		name string
		errs []error
		want string
	}

	tcases := []TestCase{
		{
			name: "Two Merged Errors",
			errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
			},
			want: "error 1: error 2",
		},
		{
			name: "Three Merged Errors",
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

			err := errutil.Merge(tc.errs...)
			assert.Error(t, err)

			got := err.Error()

			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}

			for _, inputErr := range tc.errs {
				if errors.Is(err, inputErr) {
					t.Fatalf("errors.Is(%q, %q)=true; want false", err, inputErr)
				}
			}
		})
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

type error1 struct {
	data string
}

type error2 struct {
	data int
}

func (e error1) Error() string {
	return e.data
}

func (e error2) Error() string {
	return fmt.Sprint(e.data)
}

func assertErrorIsWrapped(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("error [%v] is not wrapping [%v]", err, target)
	}
}

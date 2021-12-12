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
	type testcase struct {
		name string
		errs []error
		want []error
	}

	const (
		sentinelErr  errutil.Error = "a sentinel error"
		sentinel2Err errutil.Error = "another sentinel error"
		sentinel3Err errutil.Error = "YASE"
	)

	testcases := []testcase{
		{
			name: "single error",
			errs: []error{errors.New("single error")},
		},
		{
			name: "two errors",
			errs: []error{
				errors.New("top error"),
				errors.New("wrapped error 1"),
			},
		},
		{
			name: "three errors",
			errs: []error{
				errors.New("top error"),
				errors.New("wrapped error 1"),
				errors.New("wrapped error 2"),
			},
		},
		{
			name: "errors is nil and err",
			errs: []error{
				nil,
				sentinelErr,
			},
			want: []error{
				sentinelErr,
			},
		},
		{
			name: "errors is err and nil",
			errs: []error{
				sentinelErr,
				nil,
			},
			want: []error{
				sentinelErr,
			},
		},
		{
			name: "errors is nil,err,nil",
			errs: []error{
				nil,
				sentinelErr,
				nil,
			},
			want: []error{
				sentinelErr,
			},
		},
		{
			name: "errors interleaved with nils",
			errs: []error{
				sentinelErr,
				nil,
				sentinel2Err,
				nil,
				sentinel3Err,
				nil,
			},
			want: []error{
				sentinelErr,
				sentinel2Err,
				sentinel3Err,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			err := errutil.Chain(tc.errs...)
			assert.Error(t, err)

			got := err

			if tc.want == nil {
				tc.want = tc.errs
			}

			for i, want := range tc.want {
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

func TestErrorChainForEmptyErrListIsNil(t *testing.T) {
	assert.NoError(t, errutil.Chain())
	errs := []error{}
	assert.NoError(t, errutil.Chain(errs...))
}

func TestErrorChainWithOnlyNilErrorsIsNil(t *testing.T) {
	assert.NoError(t, errutil.Chain(nil))
	assert.NoError(t, errutil.Chain(nil, nil))
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
		name    string
		errs    []error
		reduce  errutil.Reducer
		want    string
		wantNil bool
	}

	mergeWithComma := func(err1, err2 error) error {
		return fmt.Errorf("%v,%v", err1, err2)
	}

	tests := []testcase{
		{
			name: "reducing empty err list wont call reducer and returns nil",
			errs: []error{},
			reduce: func(err1, err2 error) error {
				panic("unreachable")
			},
			wantNil: true,
		},
		{
			name: "reducing one error wont call reducer and returns error",
			errs: []error{errors.New("one")},
			reduce: func(err1, err2 error) error {
				panic("should not be called")
			},
			want: "one",
		},
		{
			name:   "reducing two errors",
			errs:   []error{errors.New("one"), errors.New("two")},
			reduce: mergeWithComma,
			want:   "one,two",
		},
		{
			name: "reducing three errors",
			errs: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: mergeWithComma,
			want:   "one,two,three",
		},
		{
			name: "filtering just first err of 3",
			errs: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return err1
			},
			want: "one",
		},
		{
			name: "filtering just first err of single err",
			errs: []error{errors.New("one")},
			reduce: func(err1, err2 error) error {
				return err1
			},
			want: "one",
		},
		{
			name: "filtering just second err",
			errs: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return err2
			},
			want: "three",
		},
		{
			name: "reduces 3 errs to nil",
			errs: []error{
				errors.New("one"),
				errors.New("two"),
				errors.New("three"),
			},
			reduce: func(err1, err2 error) error {
				return nil
			},
			wantNil: true,
		},
		{
			name: "reduces 2 errs to nil",
			errs: []error{
				errors.New("one"),
				errors.New("two"),
			},
			reduce: func(err1, err2 error) error {
				return nil
			},
			wantNil: true,
		},
		{
			name: "first is nil",
			errs: []error{
				nil,
				errors.New("error 2"),
				errors.New("error 3"),
			},
			reduce: mergeWithComma,
			want:   "error 2,error 3",
		},
		{
			name: "second is nil",
			errs: []error{
				errors.New("error 1"),
				nil,
				errors.New("error 3"),
			},
			reduce: mergeWithComma,
			want:   "error 1,error 3",
		},
		{
			name: "third is nil",
			errs: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				nil,
			},
			reduce: mergeWithComma,
			want:   "error 1,error 2",
		},
		{
			name: "multiple nils interleaved",
			errs: []error{
				nil,
				nil,
				nil,
				errors.New("error 1"),
				nil,
				nil,
				errors.New("error 2"),
				nil,
				nil,
			},
			reduce: mergeWithComma,
			want:   "error 1,error 2",
		},
		{
			name: "first err among nils",
			errs: []error{
				errors.New("error 1"),
				nil,
				nil,
				nil,
			},
			reduce: mergeWithComma,
			want:   "error 1",
		},
		{
			name: "last err among nils",
			errs: []error{
				nil,
				nil,
				nil,
				errors.New("error 1"),
			},
			reduce: mergeWithComma,
			want:   "error 1",
		},
		{
			name: "reduces list with nils to nil",
			errs: []error{
				nil,
				nil,
				nil,
			},
			reduce:  mergeWithComma,
			wantNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := errutil.Reduce(test.reduce, test.errs...)

			if test.wantNil {
				if g == nil {
					return
				}
				t.Fatalf(
					"errutil.Reduce(%v)=%q; want nil",
					test.errs,
					g,
				)
			}

			if g == nil {
				t.Fatalf(
					"errutil.Reduce(%v)=nil; want %q",
					test.errs,
					test.want,
				)
			}

			got := g.Error()
			want := test.want

			if got != want {
				t.Fatalf(
					"errutil.Reduce(%v)=%q; want=%q",
					test.errs,
					got,
					want,
				)
			}
		})
	}
}

// To test the Is method the error base type must not be comparable.
// If it is comparable, Go always just compares it, the Is method
// is just a fallback, not an override of actual comparison behavior.
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

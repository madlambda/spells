package assert_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/madlambda/spells/assert"
)

func TestPartial(t *testing.T) {
	type testcase struct {
		name   string
		obj    interface{}
		target interface{}
		fail   bool
	}

	for _, tc := range []testcase{
		{
			name:   "same numbers",
			obj:    1,
			target: 1,
		},
		{
			name:   "different int64 numbers - bigger than int32",
			obj:    int64(math.MaxInt32 + 1),
			target: int64(math.MaxInt32 + 2),
			fail:   true,
		},
		{
			name:   "same uint64 numbers - bigger than int64",
			obj:    uint64(9223372036854775807),
			target: uint64(9223372036854775807 + 1),
			fail:   true,
		},
		{
			name:   "numbers mismatch",
			obj:    1,
			target: 0,
			fail:   true,
		},
		{
			name:   "same bool",
			obj:    true,
			target: true,
		},
		{
			name:   "bool mismatch",
			obj:    true,
			target: false,
			fail:   true,
		},
		{
			name:   "same floats",
			obj:    1.2,
			target: 1.2,
		},
		{
			name:   "different floats",
			obj:    1.3,
			target: 1.2,
			fail:   true,
		},
		{
			name:   "same empty struct",
			obj:    struct{}{},
			target: struct{}{},
		},
		{
			name: "different number of fields",
			obj: struct {
				A int
			}{},
			target: struct{}{},
			fail:   true,
		},
		{
			name: "different field types",
			obj: struct {
				A int
			}{},
			target: struct {
				B string
			}{},
			fail: true,
		},
		{
			name: "same struct types different field names",
			obj: struct {
				A int
			}{},
			target: struct {
				B int
			}{},
			fail: true,
		},
		{
			name: "same struct value",
			obj: struct {
				A int
			}{1},
			target: struct {
				A int
			}{1},
		},
		{
			name: "different struct field value",
			obj: struct {
				A string
			}{"test"},
			target: struct {
				A string
			}{"test2"},
			fail: true,
		},
		{
			obj: struct {
				A string
			}{"test"},
			target: struct {
				A string
			}{"test"},
			name: "same struct field value",
		},
		{
			name: "field contains",
			obj: struct {
				A string
			}{"testing"},
			target: struct {
				A string
			}{"test"},
		},
		{
			name: "field not contains",
			obj: struct {
				A string
			}{"testing"},
			target: struct {
				A string
			}{"ABC"},
			fail: true,
		},
		{
			name: "different nested struct",
			obj: struct {
				A string
				B struct {
					C int
				}
			}{
				A: "testing",
				B: struct {
					C int
				}{1},
			},
			target: struct {
				A string
				B struct {
					C int
				}
			}{
				A: "testing",
				B: struct {
					C int
				}{2},
			},
			fail: true,
		},
		{
			name: "different nested string contains",
			obj: struct {
				A string
				B struct {
					C string
				}
			}{
				A: "testing",
				B: struct {
					C string
				}{"ABCDEFG"},
			},
			target: struct {
				A string
				B struct {
					C string
				}
			}{
				A: "testing",
				B: struct {
					C string
				}{"ABCDEF"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t, func(assert *assert.Assert, details ...interface{}) {
				if !tc.fail {
					t.Fatalf("unexpected fail: %s.%s", tc.name, errordetails(details...))
				}
			}, tc.name)
			assert.Partial(tc.obj, tc.target)
			if assert.Success() != !tc.fail {
				t.Fatalf("assert.Success() is %t but should be %t",
					assert.Success(), !tc.fail)
			}
		})
	}
}

func errordetails(details ...interface{}) string {
	if len(details) == 1 {
		return details[0].(string)
	}

	if len(details) > 1 {
		return fmt.Sprintf(details[0].(string), details[1:]...)
	}
	return ""
}

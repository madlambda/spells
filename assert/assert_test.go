package assert_test

import (
	"fmt"
	"testing"

	"github.com/madlambda/spells/assert"
)

func TestPartial(t *testing.T) {
	type testcase struct {
		obj    interface{}
		target interface{}
		msg    string
		fail   bool
	}

	t.Parallel()

	for _, tc := range []testcase{
		{
			obj:    1,
			target: 1,
			msg:    "same numbers",
		},
		{
			obj:    1,
			target: 0,
			msg:    "numbers mismatch",
			fail:   true,
		},
		{
			obj:    true,
			target: true,
			msg:    "bool mismatch",
		},
		{
			obj:    true,
			target: false,
			msg:    "bool mismatch",
			fail:   true,
		},
		{
			obj:    struct{}{},
			target: struct{}{},
			msg:    "empty struct mismatch",
		},
		{
			obj: struct {
				A int
			}{},
			target: struct{}{},
			msg:    "different number of field",
			fail:   true,
		},
		{
			obj: struct {
				A int
			}{},
			target: struct {
				B string
			}{},
			msg:  "different field types",
			fail: true,
		},
		{
			obj: struct {
				A int
			}{},
			target: struct {
				B int
			}{},
			msg:  "same struct types different field names",
			fail: true,
		},
		{
			obj: struct {
				A int
			}{1},
			target: struct {
				A int
			}{1},
			msg: "same struct value",
		},
		{
			obj: struct {
				A string
			}{"test"},
			target: struct {
				A string
			}{"test2"},
			msg:  "different struct field value",
			fail: true,
		},
		{
			obj: struct {
				A string
			}{"test"},
			target: struct {
				A string
			}{"test"},
			msg: "same struct field value",
		},
		{
			obj: struct {
				A string
			}{"testing"},
			target: struct {
				A string
			}{"test"},
			msg: "field contains",
		},
	} {
		t.Run(tc.msg, func(t *testing.T) {
			assert := assert.New(t, func(assert *assert.Assert, details ...interface{}) {
				if !tc.fail {
					t.Fatalf("unexpected fail: %s.%s", tc.msg, errordetails(details...))
				}
			})
			assert.Partial(tc.obj, tc.target, tc.msg)
			if assert.Success() != !tc.fail {
				t.Fatalf("unexpected assert result: %t", assert.Success())
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

package assert_test

import (
	"fmt"
	"testing"

	"github.com/madlambda/spells/assert"
)

func TestPartial(t *testing.T) {
	type testcase struct {
		a, b interface{}
		msg  string
		fail bool
	}

	t.Parallel()

	for _, tc := range []testcase{
		{
			a:   1,
			b:   1,
			msg: "same numbers",
		},
		{
			a:    1,
			b:    0,
			msg:  "numbers mismatch",
			fail: true,
		},
		{
			a:   true,
			b:   true,
			msg: "bool mismatch",
		},
		{
			a:    true,
			b:    false,
			msg:  "bool mismatch",
			fail: true,
		},
		{
			a:   struct{}{},
			b:   struct{}{},
			msg: "empty struct mismatch",
		},
	} {
		t.Run(tc.msg, func(t *testing.T) {
			assert := assert.New(t, func(assert *assert.Assert, details ...interface{}) {
				if !tc.fail {
					t.Fatalf("unexpected fail: %s.%s", tc.msg, errordetails(details...))
				}
			})
			assert.Partial(tc.a, tc.b, tc.msg)
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

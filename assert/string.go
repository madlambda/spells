package assert

import (
	"regexp"
	"strings"
	"testing"
)

// StringContains asserts that string s contains the subst string and calls
// the failure function with details otherwise.
func (assert *Assert) StringContains(s string, substr string, details ...interface{}) {
	assert.IsTrue(strings.Contains(s, substr),
		errctx(details, "strings.Contains(%q, %q)",
			s, substr))
}

// StringMatch asserts that string matches the regex pattern and calls
// the failure function with details otherwise.
func (assert *Assert) StringMatch(pattern string, str string, details ...interface{}) {
	found, err := regexp.MatchString(pattern, str)
	assert.NoError(err, errctx(details, "failed to build regexp pattern %q", pattern))
	assert.IsTrue(found, errctx(details, "pattern[%s] not found in [%s]", pattern, str))
}

// StringContains asserts that string s contains the subst string and calls
// the Fatal() function with details otherwise.
func StringContains(t testing.TB, s, substr string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringContains(s, substr, details...)
}

// StringMatch asserts that string matches the regex pattern and calls
// the Fatal() function with details otherwise.
func StringMatch(t testing.TB, pattern string, str string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringMatch(pattern, str, details...)
}

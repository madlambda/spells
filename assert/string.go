package assert

import (
	"regexp"
	"strings"
	"testing"
)

// StringContains asserts that string s contains the subst string and calls
// the failure function with details otherwise.
func (assert *Assert) StringContains(s string, substr string, details ...interface{}) {
	assert.IsTrue(strings.Contains(s, substr), "strings.Contains(%q, %q).%s",
		s, substr, errordetails(details...))
}

// StringMatch asserts that string matches the regex pattern and calls
// the failure function with details otherwise.
func (assert *Assert) StringMatch(pattern string, str string, details ...interface{}) {
	matched, err := regexp.MatchString(pattern, str)
	assert.IsFalse(err != nil, "err != nil. %v", pattern, err)
	assert.IsTrue(matched, "pattern[%s] not found in [%s].%s",
		pattern, str, errordetails(details...))
}

// StringContains asserts that string s contains the subst string and calls
// the Fatal() function with details otherwise.
func StringContains(t *testing.T, s, substr string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringContains(s, substr, details...)
}

// StringMatch asserts that string matches the regex pattern and calls
// the Fatal() function with details otherwise.
func StringMatch(t *testing.T, pattern string, str string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringMatch(pattern, str, details...)
}

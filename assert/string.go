package assert

import (
	"regexp"
	"strings"
	"testing"
)

// StringContains asserts that string s contains the substring and calls
// t.Fatal using the details parameter otherwise.
// The details parameter can be a single string of a format string + parameters.
func (assert *Assert) StringContains(s string, substr string, details ...interface{}) {
	assert.True(strings.Contains(s, substr), "strings.Contains(%q, %q).%s",
		s, substr, errordetails(details...))
}

// StringMatch asserts that string matches the regex pattern and calls
// t.Fatal using the details parameter otherwise.
// The details parameter can be a single string of a format string + parameters.
func (assert *Assert) StringMatch(pattern string, str string, details ...interface{}) {
	matched, err := regexp.MatchString(pattern, str)
	assert.False(err != nil, "err != nil. %v", pattern, err)
	assert.True(matched, "pattern[%s] not found in [%s].%s",
		pattern, str, errordetails(details...))
}

// StringContains asserts that string s contains the substring and calls
// t.Fatal using the details parameter otherwise.
// The details parameter can be a single string of a format string + parameters.
func StringContains(t *testing.T, s, substr string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringContains(s, substr, details...)
}

// StringMatch asserts that string matches the regex pattern and calls
// t.Fatal using the details parameter otherwise.
// The details parameter can be a single string of a format string + parameters.
func StringMatch(t *testing.T, pattern string, str string, details ...interface{}) {
	assert := New(t, Fatal)
	assert.StringMatch(pattern, str, details...)
}

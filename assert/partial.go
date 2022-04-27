package assert

import (
	"reflect"
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

func (assert *Assert) Partial(obj interface{}, target interface{}, details ...interface{}) {
	elem := reflect.ValueOf(obj)
	targ := reflect.ValueOf(target)

	assert.True(elem.Kind() == targ.Kind(), "wanted object type[%s] but got[%s]",
		targ.Kind(), elem.Kind())

	if targ.Kind() == reflect.Ptr {
		elem = elem.Elem()
		targ = targ.Elem()

		assert.True(elem.Kind() == targ.Kind(), "wanted object type[%s] but got[%s]",
			targ.Kind(), elem.Kind())
	}

	switch targ.Kind() {
	case reflect.Bool:
		assert.Bool(targ.Bool(), elem.Bool(), "boolean mismatch")
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		// TODO(i4k): properly compare without conversion.
		assert.EqualInts(int(targ.Int()), int(elem.Int()), details...)
	case reflect.String:
		assert.StringContains(elem.String(), targ.String(), details...)
	case reflect.Struct:
		assert.True(targ.Type().Name() == elem.Type().Name(), "struct type mismatch")
		assert.partialStruct(elem, targ, details...)
	default:
		assert.t.Fatalf("Partial does not support comparing %s", targ.Kind())
	}
}

func (assert *Assert) partialStruct(obj reflect.Value, target reflect.Value, details ...interface{}) {
	objtype := obj.Type()
	targtype := target.Type()

	assert.EqualInts(obj.NumField(), target.NumField(),
		"number of struct fields mismatch.%s", errordetails(details...))

	for i := 0; i < target.NumField(); i++ {
		ofield := objtype.Field(i)
		tfield := targtype.Field(i)

		assert.Bool(ofield.Anonymous, tfield.Anonymous,
			"embedded field and non-embedded field.%s", errordetails(details...))

		assert.True(ofield.Type == tfield.Type,
			"field type mismatch: index %d (%s.%s (%s) == %s.%s (%s).%s", i,
			objtype.Name(), ofield.Name, ofield.Type,
			targtype.Name(), tfield.Name, tfield.Type,
			errordetails(details...),
		)

		assert.True(ofield.Name == tfield.Name,
			"field name mismatch: index %d (%s.%s (%s) == %s.%s (%s).%s",
			i,
			objtype.Name(), ofield.Name, ofield.Type,
			targtype.Name(), tfield.Name, tfield.Type,
			errordetails(details...),
		)

		assert.Partial(obj.Field(i).Interface(), target.Field(i).Interface(), details...)
	}
}

func Partial(t *testing.T, obj interface{}, target interface{}, details ...interface{}) {
	assert := New(t, Fatal)
	assert.Partial(obj, target, details...)
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

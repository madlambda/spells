package assert

import (
	"reflect"
	"testing"
)

func (assert *Assert) Partial(obj interface{}, target interface{}) {
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
		assert.EqualInts(int(targ.Int()), int(elem.Int()), "int mismatch")
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		assert.EqualUints(uint64(targ.Uint()), uint64(elem.Uint()), "uint mismatch")
	case reflect.String:
		assert.StringContains(elem.String(), targ.String(), "string mismatch")
	case reflect.Struct:
		assert.partialStruct(elem, targ, "struct mismatch")
	default:
		assert.t.Fatalf("Partial does not support comparing %s", targ.Kind())
	}
}

func (assert *Assert) partialStruct(obj reflect.Value, target reflect.Value, details ...interface{}) {
	objtype := obj.Type()
	targtype := target.Type()

	assert.True(targtype.Name() == objtype.Name(), "struct type mismatch.%s",
		errordetails(details...))
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

		assert.Partial(obj.Field(i).Interface(), target.Field(i).Interface())
	}
}

func Partial(t *testing.T, obj interface{}, target interface{}, details ...interface{}) {
	assert := New(t, Fatal)
	assert.Partial(obj, target)
}

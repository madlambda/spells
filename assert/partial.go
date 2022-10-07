package assert

import (
	"reflect"
	"testing"
)

// Partial recursively asserts that obj partially matches target.
// Below are the assertion rules:
// - booleans, integers, floats and complex numbers must be equal.
// - strings, slices, array and map: the obj must contains the target.
// - structs: the obj fields must recursively 8match the target fields.
func (assert *Assert) Partial(obj, target interface{}, details ...interface{}) {
	assert.t.Helper()
	elem := reflect.ValueOf(obj)
	targ := reflect.ValueOf(target)

	assert.failif(targ.IsValid() != elem.IsValid(),
		details, "internal reflection property mismatch")

	if !targ.IsValid() {
		return
	}

	assert.failif(elem.Kind() != targ.Kind(), details,
		"wanted object kind[%s] but got[%s]", targ.Kind(), elem.Kind(),
	)

	if targ.Kind() == reflect.Ptr {
		elem = elem.Elem()
		targ = targ.Elem()

		assert.failif(elem.Kind() != targ.Kind(), details,
			"wanted object type[%s] but got[%s]", targ.Kind(), elem.Kind())

		assert.failif(targ.IsValid() != elem.IsValid(), details,
			"internal reflection property mismatch")

		if !targ.IsValid() {
			return
		}
	}

	switch targ.Kind() {
	case reflect.Bool:
		assert.EqualBools(targ.Bool(), elem.Bool(),
			errctx(details, "boolean mismatch"))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		assert.EqualInts(int(targ.Int()), int(elem.Int()),
			errctx(details, "int mismatch"))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		assert.EqualUints(uint64(targ.Uint()), uint64(elem.Uint()),
			errctx(details, "uint mismatch"))
	case reflect.Float32, reflect.Float64:
		assert.EqualFloats(targ.Float(), elem.Float(),
			errctx(details, "float mismatch"))
	case reflect.Complex64, reflect.Complex128:
		assert.EqualComplexes(targ.Complex(), elem.Complex(),
			errctx(details, "complex numbers mismatch"))
	case reflect.String:
		assert.StringContains(elem.String(), targ.String(),
			errctx(details, "string mismatch"))
	case reflect.Struct:
		assert.partialStruct(elem.Interface(), targ.Interface(),
			errctx(details, "struct mismatch"))
	case reflect.Slice:
		assert.failif(targ.Len() > elem.Len(), details,
			"target length is bigger than object: %d > %d", targ.Len(), elem.Len())
		for i := 0; i < targ.Len(); i++ {
			assert.Partial(elem.Index(i).Interface(), targ.Index(i).Interface(),
				errctx(details, "slice index %d mismatch", i))
		}
	case reflect.Array:
		assert.failif(targ.Len() > elem.Len(), details,
			"target length is bigger than object")
		for i := 0; i < targ.Len(); i++ {
			assert.Partial(elem.Index(i).Interface(), targ.Index(i).Interface(),
				errctx(details, "array index %d mismatch", i))
		}
	case reflect.Map:
		assert.failif(targ.Len() > elem.Len(), details, "target map has more keys")
		tkeys := targ.MapKeys()
		for _, tkey := range tkeys {
			tval := targ.MapIndex(tkey)
			eval := elem.MapIndex(tkey)
			if !eval.IsValid() {
				assert.fail(details, "target key %v not found in object", tkey)
				continue
			}
			assert.Partial(tval.Interface(), eval.Interface(),
				errctx(details, "comparing map values"))
		}
	default:
		assert.t.Fatalf("Partial does not support comparing %s", targ.Kind())
	}
}

func (assert *Assert) partialStruct(obj interface{}, target interface{}, details ...interface{}) {
	assert.t.Helper()
	elem := reflect.ValueOf(obj)
	targ := reflect.ValueOf(target)
	objtype := elem.Type()
	targtype := targ.Type()

	assert.failif(targ.NumField() > elem.NumField(),
		details, "target.NumField() > obj.NumField()")

	for i := 0; i < targtype.NumField(); i++ {
		tfield := targtype.Field(i)
		ofield, found := objtype.FieldByName(tfield.Name)

		assert.failif(!found, details, "field %s not found in the object", tfield.Name)

		if !found {
			continue
		}

		assert.failif(ofield.Anonymous != tfield.Anonymous,
			details, "embedded field and non-embedded field")

		fieldElem := elem.FieldByName(ofield.Name)
		assert.failif(!fieldElem.IsValid(), details, "object field %q not found",
			ofield.Name)

		if !fieldElem.IsValid() {
			continue
		}

		fieldTarg := targ.FieldByName(tfield.Name)
		assert.failif(!fieldTarg.IsValid(), details, "target field %q not found",
			tfield.Name)

		// tfield.IsExported() was introduced only in go1.17
		if tfield.PkgPath == "" {
			assert.Partial(
				fieldElem.Interface(), fieldTarg.Interface(),
				errctx(details, "comparing struct field %s and %s",
					tfield.Name, ofield.Name))

		}
	}
}

func Partial(t testing.TB, obj interface{}, target interface{}, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal, details...)
	assert.Partial(obj, target)
}

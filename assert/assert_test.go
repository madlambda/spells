package assert_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/madlambda/spells/assert"
)

type testIface interface {
	A() bool
}

type testStruct1 struct {
	Val        int
	IfaceField testIface
}

func (testStruct1) A() bool { return true }

type testStruct2 struct {
	Val        int
	IfaceField testIface
}

func (testStruct2) A() bool { return true }

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
			name: "struct with no fields match any struct",
			obj: struct {
				A int
			}{},
			target: struct{}{},
		},
		{
			name:   "different struct names",
			obj:    testStruct1{},
			target: testStruct2{},
		},
		{
			name: "comparing different struct that fields match",
			obj: testStruct1{
				Val: 10,
			},
			target: struct {
				Val int
			}{10},
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
			name: "different unexported field types",
			obj: struct {
				a int
			}{},
			target: struct {
				a string
			}{},
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
			name: "same struct int value",
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
			name: "same struct field value",
			obj: struct {
				A string
			}{"test"},
			target: struct {
				A string
			}{"test"},
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
		{
			name:   "same interfaces - sanity check",
			obj:    testIface(testStruct1{}),
			target: testIface(testStruct1{}),
		},
		{
			name:   "same empty slices",
			obj:    []int{},
			target: []int{},
		},
		{
			name:   "same slices",
			obj:    []int{1},
			target: []int{1},
		},
		{
			name:   "same slices - contains target elements in order",
			obj:    []int{1, 2, 3},
			target: []int{1, 2},
		},
		{
			name:   "same slices",
			obj:    []string{"test"},
			target: []string{"test"},
		},
		{
			name:   "same slices with string contains",
			obj:    []string{"testing"},
			target: []string{"test"},
		},
		{
			name:   "same arrays",
			obj:    [2]string{"1", "2"},
			target: [2]string{"1", "2"},
		},
		{
			name:   "empty target slice matches any value",
			obj:    []string{"test"},
			target: []string{},
		},
		{
			name: "struct with different slices",
			obj: struct {
				A []int
			}{[]int{1, 2}},
			target: struct {
				A []int
			}{[]int{1, 2}},
		},
		{
			name:   "same maps",
			obj:    map[int]int{},
			target: map[int]int{},
		},
		{
			name: "same maps with values",
			obj: map[int]int{
				1: 1,
				3: 1,
			},
			target: map[int]int{
				1: 1,
				3: 1,
			},
		},
		{
			name: "same maps - obj contains target keys",
			obj: map[int]int{
				1:   1,
				3:   1,
				666: 666,
			},
			target: map[int]int{
				1: 1,
				3: 1,
			},
		},
		{
			name: "different maps - obj doesn't contains all target keys",
			obj: map[int]int{
				1:   1,
				3:   1,
				666: 666,
			},
			target: map[int]int{
				1:   1,
				3:   1,
				667: 666,
			},
			fail: true,
		},
		{
			name: "different maps - same keys, different value",
			obj: map[int]int{
				1:   1,
				3:   1,
				666: 666,
			},
			target: map[int]int{
				1:   1,
				3:   1,
				666: 667,
			},
			fail: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failures := 0
			assert := assert.New(t, func(assert *assert.Assert, msg string) {
				failures++
				if !tc.fail {
					t.Fatalf("unexpected fail: %s: %s", tc.name, msg)
				}
			}, tc.name)
			assert.Partial(tc.obj, tc.target)
			if failures > 0 != tc.fail {
				t.Fatalf("there was %d errors but tc.fail is %t",
					failures, tc.fail)
			}
		})
	}
}

func TestAssertErrMessages(t *testing.T) {
	t.Run("error: no details", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil"
			assert.EqualStrings(t, want, got)
		})
		a.Error(nil)
	})

	t.Run("error: constructor msg and details msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: func fmt 777: constructor fmt 666"
			assert.EqualStrings(t, want, got)
		}, "constructor fmt %d", 666)
		a.Error(nil, "func fmt %d", 777)
	})

	t.Run("error: constructor msg and details msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: func msg: constructor msg"
			assert.EqualStrings(t, want, got)
		}, "constructor msg")
		a.Error(nil, "func msg")
	})

	t.Run("error: constructor msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: constructor fmt 666"
			assert.EqualStrings(t, want, got)
		}, "constructor fmt %d", 666)
		a.Error(nil)
	})

	t.Run("error: constructor msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: constructor msg"
			assert.EqualStrings(t, want, got)
		}, "constructor msg")
		a.Error(nil)
	})

	t.Run("error: detail msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: func fmt 666"
			assert.EqualStrings(t, want, got)
		})
		a.Error(nil, "func fmt %d", 666)
	})

	t.Run("error: detail msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := "expected error, got nil: func msg"
			assert.EqualStrings(t, want, got)
		})
		a.Error(nil, "func msg")
	})
}

func TestAssertNoErrMessages(t *testing.T) {
	err := errors.New("test")
	errmsg := fmt.Sprintf("unexpected error[%v]", err)

	t.Run("no error: no details", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := errmsg
			assert.EqualStrings(t, want, got)
		})
		a.NoError(err)
	})

	t.Run("no error: constructor msg and details msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: func fmt 777: constructor fmt 666", errmsg)
			assert.EqualStrings(t, want, got)
		}, "constructor fmt %d", 666)
		a.NoError(err, "func fmt %d", 777)
	})

	t.Run("no error: constructor msg and details msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: func msg: constructor msg", errmsg)
			assert.EqualStrings(t, want, got)
		}, "constructor msg")
		a.NoError(err, "func msg")
	})

	t.Run("no error: constructor msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: constructor fmt 666", errmsg)
			assert.EqualStrings(t, want, got)
		}, "constructor fmt %d", 666)
		a.NoError(err)
	})

	t.Run("no error: constructor msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: constructor msg", errmsg)
			assert.EqualStrings(t, want, got)
		}, "constructor msg")
		a.NoError(err)
	})

	t.Run("no error: detail msg with fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: func fmt 666", errmsg)
			assert.EqualStrings(t, want, got)
		})
		a.NoError(err, "func fmt %d", 666)
	})

	t.Run("error: detail msg no fmt", func(t *testing.T) {
		a := assert.New(t, func(a *assert.Assert, got string) {
			want := fmt.Sprintf("%s: func msg", errmsg)
			assert.EqualStrings(t, want, got)
		})
		a.NoError(err, "func msg")
	})
}

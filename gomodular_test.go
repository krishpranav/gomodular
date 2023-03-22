package gomodular_test

import (
	"errors"
	"testing"

	"github.com/krishpranav/gomodular"
	"github.com/stretchr/testify/assert"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c Circle) GetArea() int {
	return c.a
}

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
	return true
}

var instance = gomodular.New()

func TestGomodular_Singleton(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(666)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestGomodular_SingletonLazy(t *testing.T) {
	err := instance.SingletonLazy(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(666)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestGomodular_Singleton_With_Missing_Dependency_Resolve(t *testing.T) {
	err := instance.Singleton(func(db Database) Shape {
		return &Circle{a: 13}
	})
	assert.EqualError(t, err, "gomodular: no concrete found for: gomodular_test.Database")
}

func TestGomodular_Singleton_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.Singleton(func() {})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_SingletonLazy_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.SingletonLazy(func() {})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_Singleton_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.Singleton(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.Error(t, err, "app: error")
}

func TestGomodular_SingletonLazy_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.SingletonLazy(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.NoError(t, err)

	var s Shape
	err = instance.Resolve(&s)
	assert.Error(t, err, "app: error")
}

func TestGomodular_Singleton_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.Singleton("STRING!")
	assert.EqualError(t, err, "gomodular: the resolver must be a function")
}

func TestGomodular_SingletonLazy_With_NonFunction_Resolver_It_Should_Fail(t *testing.T) {
	err := instance.SingletonLazy("STRING!")
	assert.EqualError(t, err, "gomodular: the resolver must be a function")
}

func TestGomodular_Singleton_With_Resolvable_Arguments(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), 666)
		return &MySQL{}
	})
	assert.NoError(t, err)
}

func TestGomodular_SingletonLazy_With_Resolvable_Arguments(t *testing.T) {
	err := instance.SingletonLazy(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.SingletonLazy(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), 666)
		return &MySQL{}
	})
	assert.NoError(t, err)

	var s Shape
	err = instance.Resolve(&s)
	assert.NoError(t, err)
}

func TestGomodular_Singleton_With_Non_Resolvable_Arguments(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func(s Shape) Shape {
		return &Circle{a: s.GetArea()}
	})
	assert.EqualError(t, err, "gomodular: resolver function signature is invalid - depends on abstract it returns")
}

func TestGomodular_SingletonLazy_With_Non_Resolvable_Arguments(t *testing.T) {
	instance.Reset()

	err := instance.SingletonLazy(func(s Shape) Shape {
		return &Circle{a: s.GetArea()}
	})
	assert.EqualError(t, err, "gomodular: resolver function signature is invalid - depends on abstract it returns")
}

func TestGomodular_NamedSingleton(t *testing.T) {
	err := instance.NamedSingleton("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestGomodular_NamedSingletonLazy(t *testing.T) {
	err := instance.NamedSingletonLazy("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestGomodular_Transient(t *testing.T) {
	err := instance.Transient(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(13)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestGomodular_TransientLazy(t *testing.T) {
	err := instance.TransientLazy(func() Shape {
		return &Circle{a: 666}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s1 Shape) {
		s1.SetArea(13)
	})
	assert.NoError(t, err)

	err = instance.Call(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, 666)
	})
	assert.NoError(t, err)
}

func TestGomodular_Transient_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.Transient(func() {})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_TransientLazy_With_Resolve_That_Returns_Nothing(t *testing.T) {
	err := instance.TransientLazy(func() {})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_Transient_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.Transient(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.Error(t, err, "app: error")

	firstCall := true
	err = instance.Transient(func() (Database, error) {
		if firstCall {
			firstCall = false
			return &MySQL{}, nil
		}
		return nil, errors.New("app: second call error")
	})
	assert.NoError(t, err)

	var db Database
	err = instance.Resolve(&db)
	assert.Error(t, err, "app: second call error")
}

func TestGomodular_TransientLazy_With_Resolve_That_Returns_Error(t *testing.T) {
	err := instance.TransientLazy(func() (Shape, error) {
		return nil, errors.New("app: error")
	})
	assert.NoError(t, err)

	var s Shape
	err = instance.Resolve(&s)
	assert.Error(t, err, "app: error")

	firstCall := true
	err = instance.TransientLazy(func() (Database, error) {
		if firstCall {
			firstCall = false
			return &MySQL{}, nil
		}
		return nil, errors.New("app: second call error")
	})
	assert.NoError(t, err)

	var db Database
	err = instance.Resolve(&db)
	assert.NoError(t, err)

	err = instance.Resolve(&db)
	assert.Error(t, err, "app: second call error")
}

func TestGomodular_Transient_With_Resolve_With_Invalid_Signature_It_Should_Fail(t *testing.T) {
	err := instance.Transient(func() (Shape, Database, error) {
		return nil, nil, nil
	})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_TransientLazy_With_Resolve_With_Invalid_Signature_It_Should_Fail(t *testing.T) {
	err := instance.TransientLazy(func() (Shape, Database, error) {
		return nil, nil, nil
	})
	assert.Error(t, err, "gomodular: resolver function signature is invalid")
}

func TestGomodular_NamedTransient(t *testing.T) {
	err := instance.NamedTransient("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestGomodular_NamedTransientLazy(t *testing.T) {
	err := instance.NamedTransientLazy("theCircle", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	var sh Shape
	err = instance.NamedResolve(&sh, "theCircle")
	assert.NoError(t, err)
	assert.Equal(t, sh.GetArea(), 13)
}

func TestGomodular_Call_With_Multiple_Resolving(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.NoError(t, err)
}

func TestGomodular_Call_With_Dependency_Missing_In_Chain(t *testing.T) {
	var instance = gomodular.New()
	err := instance.SingletonLazy(func() (Database, error) {
		var s Shape
		if err := instance.Resolve(&s); err != nil {
			return nil, err
		}
		return &MySQL{}, nil
	})
	assert.NoError(t, err)

	err = instance.Call(func(m Database) {
		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
	assert.EqualError(t, err, "gomodular: no concrete found for: gomodular_test.Shape")
}

func TestGomodular_Call_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Call("STRING!")
	assert.EqualError(t, err, "gomodular: invalid function")
}

func TestGomodular_Call_With_Second_UnBounded_Argument(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape, d Database) {})
	assert.EqualError(t, err, "gomodular: no concrete found for: gomodular_test.Database")
}

func TestGomodular_Call_With_A_Returning_Error(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape) error {
		return errors.New("app: some context error")
	})
	assert.EqualError(t, err, "app: some context error")
}

func TestGomodular_Call_With_A_Returning_Nil_Error(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape) error {
		return nil
	})
	assert.Nil(t, err)
}

func TestGomodular_Call_With_Invalid_Signature(t *testing.T) {
	instance.Reset()

	err := instance.Singleton(func() Shape {
		return &Circle{}
	})
	assert.NoError(t, err)

	err = instance.Call(func(s Shape) (int, error) {
		return 13, errors.New("app: some context error")
	})
	assert.EqualError(t, err, "gomodular: receiver function signature is invalid")
}

func TestGomodular_Resolve_With_Reference_As_Resolver(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	var (
		s Shape
		d Database
	)

	err = instance.Resolve(&s)
	assert.NoError(t, err)
	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	err = instance.Resolve(&d)
	assert.NoError(t, err)
	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestGomodular_Resolve_With_Unsupported_Receiver_It_Should_Fail(t *testing.T) {
	err := instance.Resolve("STRING!")
	assert.EqualError(t, err, "gomodular: invalid abstraction")
}

func TestGomodular_Resolve_With_NonReference_Receiver_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Resolve(s)
	assert.EqualError(t, err, "gomodular: invalid abstraction")
}

func TestGomodular_Resolve_With_UnBounded_Reference_It_Should_Fail(t *testing.T) {
	instance.Reset()

	var s Shape
	err := instance.Resolve(&s)
	assert.EqualError(t, err, "gomodular: no concrete found for: gomodular_test.Shape")
}

func TestGomodular_Fill_With_Struct_Pointer(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.NamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `gomodular:"type"`
		D Database `gomodular:"type"`
		C Shape    `gomodular:"name"`
		X string
	}{}

	err = instance.Fill(&myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.S)
	assert.IsType(t, &MySQL{}, myApp.D)
}

func TestGomodular_Fill_Unexported_With_Struct_Pointer(t *testing.T) {
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape    `gomodular:"type"`
		d Database `gomodular:"type"`
		y int
	}{}

	err = instance.Fill(&myApp)
	assert.NoError(t, err)

	assert.IsType(t, &Circle{}, myApp.s)
	assert.IsType(t, &MySQL{}, myApp.d)
}

func TestGomodular_Fill_With_Invalid_Field_It_Should_Fail(t *testing.T) {
	err := instance.NamedSingleton("C", func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	type App struct {
		S string `gomodular:"name"`
	}

	myApp := App{}

	err = instance.Fill(&myApp)
	assert.EqualError(t, err, "gomodular: cannot make S field")
}

func TestGomodular_Fill_With_Invalid_Tag_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `gomodular:"invalid"`
	}

	myApp := App{}

	err := instance.Fill(&myApp)
	assert.EqualError(t, err, "gomodular: S has an invalid struct tag")
}

func TestGomodular_Fill_With_Invalid_Field_Name_It_Should_Fail(t *testing.T) {
	type App struct {
		S string `gomodular:"name"`
	}

	myApp := App{}

	err := instance.Fill(&myApp)
	assert.EqualError(t, err, "gomodular: cannot make S field")
}

func TestGomodular_Fill_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	invalidStruct := 0
	err := instance.Fill(&invalidStruct)
	assert.EqualError(t, err, "gomodular: invalid structure")
}

func TestGomodular_Fill_With_Invalid_Pointer_It_Should_Fail(t *testing.T) {
	var s Shape
	err := instance.Fill(s)
	assert.EqualError(t, err, "gomodular: invalid structure")
}

func TestGomodular_Fill_With_Dependency_Missing_In_Chain(t *testing.T) {
	var instance = gomodular.New()
	err := instance.Singleton(func() Shape {
		return &Circle{a: 5}
	})
	assert.NoError(t, err)

	err = instance.NamedSingletonLazy("C", func() (Shape, error) {
		var s Shape
		if err := instance.NamedResolve(&s, "foo"); err != nil {
			return nil, err
		}
		return &Circle{a: 5}, nil
	})
	assert.NoError(t, err)

	err = instance.Singleton(func() Database {
		return &MySQL{}
	})
	assert.NoError(t, err)

	myApp := struct {
		S Shape    `gomodular:"type"`
		D Database `gomodular:"type"`
		C Shape    `gomodular:"name"`
		X string
	}{}

	err = instance.Fill(&myApp)
	assert.EqualError(t, err, "gomodular: no concrete found for: gomodular_test.Shape")
}

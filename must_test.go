package gomodular_test

import (
	"errors"
	"testing"

	"github.com/krishpranav/gomodular"
)

func TestMustSingleton_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustSingleton(c, func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustSingletonLazy_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustSingletonLazy(c, func() {})
	t.Errorf("panic expcted.")
}

func TestMustNamedSingleton_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustNamedSingleton(c, "name", func() (Shape, error) {
		return nil, errors.New("error")
	})
	t.Errorf("panic expcted.")
}

func TestMustNamedSingletonLazy_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustNamedSingletonLazy(c, "name", func() {})
	t.Errorf("panic expcted.")
}

func TestMustTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustTransient(c, func() (Shape, error) {
		return nil, errors.New("error")
	})

	var resVal Shape
	gomodular.MustResolve(c, &resVal)

	t.Errorf("panic expcted.")
}

func TestMustTransientLazy_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustTransientLazy(c, func() {
	})

	var resVal Shape
	gomodular.MustResolve(c, &resVal)

	t.Errorf("panic expcted.")
}

func TestMustNamedTransient_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustNamedTransient(c, "name", func() (Shape, error) {
		return nil, errors.New("error")
	})

	var resVal Shape
	gomodular.MustNamedResolve(c, &resVal, "name")

	t.Errorf("panic expcted.")
}

func TestMustNamedTransientLazy_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustNamedTransientLazy(c, "name", func() {
	})

	var resVal Shape
	gomodular.MustNamedResolve(c, &resVal, "name")

	t.Errorf("panic expcted.")
}

func TestMustCall_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	defer func() { recover() }()
	gomodular.MustCall(c, func(s Shape) {
		s.GetArea()
	})
	t.Errorf("panic expcted.")
}

func TestMustResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	var s Shape

	defer func() { recover() }()
	gomodular.MustResolve(c, &s)
	t.Errorf("panic expcted.")
}

func TestMustNamedResolve_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	var s Shape

	defer func() { recover() }()
	gomodular.MustNamedResolve(c, &s, "name")
	t.Errorf("panic expcted.")
}

func TestMustFill_It_Should_Panic_On_Error(t *testing.T) {
	c := gomodular.New()

	myApp := struct {
		S Shape `gomodular:"type"`
	}{}

	defer func() { recover() }()
	gomodular.MustFill(c, &myApp)
	t.Errorf("panic expcted.")
}

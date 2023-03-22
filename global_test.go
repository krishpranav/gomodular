package gomodular_test

import (
	"testing"

	"github.com/krishpranav/gomodular"
	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	gomodular.Reset()

	err := gomodular.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestSingletonLazy(t *testing.T) {
	gomodular.Reset()

	err := gomodular.SingletonLazy(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedSingleton(t *testing.T) {
	gomodular.Reset()

	err := gomodular.NamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedSingletonLazy(t *testing.T) {
	gomodular.Reset()

	err := gomodular.NamedSingletonLazy("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestTransient(t *testing.T) {
	gomodular.Reset()

	err := gomodular.Transient(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestTransientLazy(t *testing.T) {
	gomodular.Reset()

	err := gomodular.TransientLazy(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedTransient(t *testing.T) {
	gomodular.Reset()

	err := gomodular.NamedTransient("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestNamedTransientLazy(t *testing.T) {
	gomodular.Reset()

	err := gomodular.NamedTransientLazy("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)
}

func TestCall(t *testing.T) {
	gomodular.Reset()

	err := gomodular.Call(func() {})
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {
	gomodular.Reset()

	var s Shape

	err := gomodular.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = gomodular.Resolve(&s)
	assert.NoError(t, err)
}

func TestNamedResolve(t *testing.T) {
	gomodular.Reset()

	var s Shape

	err := gomodular.NamedSingleton("rounded", func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	err = gomodular.NamedResolve(&s, "rounded")
	assert.NoError(t, err)
}

func TestFill(t *testing.T) {
	gomodular.Reset()

	err := gomodular.Singleton(func() Shape {
		return &Circle{a: 13}
	})
	assert.NoError(t, err)

	myApp := struct {
		s Shape `Global:"type"`
	}{}

	err = gomodular.Fill(&myApp)
	assert.NoError(t, err)
}

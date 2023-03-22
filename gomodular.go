package gomodular

import (
	"errors"
	"reflect"
)

type binding struct {
	resolver    interface{}
	concrete    interface{}
	isSingleton bool
}

func (b *binding) make(c Gomodular) (interface{}, error) {
	if b.concrete != nil {
		return b.concrete, nil
	}

	retVal, err := c.invoke(b.resolver)
	if b.isSingleton {
		b.concrete = retVal
	}

	return retVal, err
}

type Gomodular map[reflect.Type]map[string]*binding

func New() Gomodular {
	return make(Gomodular)
}

func (c Gomodular) invoke(function interface{}) (interface{}, error) {
	arguments, err := c.arguments(function)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(function).Call(arguments)
	if len(values) == 2 && values[1].CanInterface() {
		if err, ok := values[1].Interface().(error); ok {
			return values[0].Interface(), err
		}
	}
	return values[0].Interface(), nil
}

func (c Gomodular) arguments(function interface{}) ([]reflect.Value, error) {
	reflectedFunction := reflect.TypeOf(function)
	argumentsCount := reflectedFunction.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := reflectedFunction.In(i)
		if concrete, exist := c[abstraction][""]; exist {
			instance, err := concrete.make(c)
			if err != nil {
				return nil, err
			}
			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, errors.New("container: no concrete found for: " + abstraction.String())
		}
	}

	return arguments, nil
}

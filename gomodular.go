package gomodular

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
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

func (c Gomodular) bind(resolver interface{}, name string, isSingleton bool, isLazy bool) error {
	reflectedResolver := reflect.TypeOf(resolver)
	if reflectedResolver.Kind() != reflect.Func {
		return errors.New("gomodular: the resolver must be a function")
	}

	if reflectedResolver.NumOut() > 0 {
		if _, exist := c[reflectedResolver.Out(0)]; !exist {
			c[reflectedResolver.Out(0)] = make(map[string]*binding)
		}
	}

	if err := c.validateResolverFunction(reflectedResolver); err != nil {
		return err
	}

	var concrete interface{}
	if !isLazy {
		var err error
		concrete, err = c.invoke(resolver)
		if err != nil {
			return err
		}
	}

	if isSingleton {
		c[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, concrete: concrete, isSingleton: isSingleton}
	} else {
		c[reflectedResolver.Out(0)][name] = &binding{resolver: resolver, isSingleton: isSingleton}
	}

	return nil
}

func (c Gomodular) validateResolverFunction(funcType reflect.Type) error {
	retCount := funcType.NumOut()

	if retCount == 0 || retCount > 2 {
		return errors.New("gomodular: resolver function signature is invalid - it must return abstract, or abstract and error")
	}

	resolveType := funcType.Out(0)
	for i := 0; i < funcType.NumIn(); i++ {
		if funcType.In(i) == resolveType {
			return fmt.Errorf("gomodular: resolver function signature is invalid - depends on abstract it returns")
		}
	}

	return nil
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
			return nil, errors.New("gomodular: no concrete found for: " + abstraction.String())
		}
	}

	return arguments, nil
}

func (c Gomodular) Reset() {
	for k := range c {
		delete(c, k)
	}
}

func (c Gomodular) Singleton(resolver interface{}) error {
	return c.bind(resolver, "", true, false)
}

func (c Gomodular) SingletonLazy(resolver interface{}) error {
	return c.bind(resolver, "", true, true)
}

func (c Gomodular) NamedSingleton(name string, resolver interface{}) error {
	return c.bind(resolver, name, true, false)
}

func (c Gomodular) NamedSingletonLazy(name string, resolver interface{}) error {
	return c.bind(resolver, name, true, true)
}

func (c Gomodular) Transient(resolver interface{}) error {
	return c.bind(resolver, "", false, false)
}

func (c Gomodular) TransientLazy(resolver interface{}) error {
	return c.bind(resolver, "", false, true)
}

func (c Gomodular) NamedTransient(name string, resolver interface{}) error {
	return c.bind(resolver, name, false, false)
}

func (c Gomodular) NamedTransientLazy(name string, resolver interface{}) error {
	return c.bind(resolver, name, false, true)
}

func (c Gomodular) Call(function interface{}) error {
	receiverType := reflect.TypeOf(function)
	if receiverType == nil || receiverType.Kind() != reflect.Func {
		return errors.New("gomodular: invalid function")
	}

	arguments, err := c.arguments(function)
	if err != nil {
		return err
	}

	result := reflect.ValueOf(function).Call(arguments)

	if len(result) == 0 {
		return nil
	} else if len(result) == 1 && result[0].CanInterface() {
		if result[0].IsNil() {
			return nil
		}
		if err, ok := result[0].Interface().(error); ok {
			return err
		}
	}

	return errors.New("gomodular: receiver function signature is invalid")
}

func (c Gomodular) Resolve(abstraction interface{}) error {
	return c.NamedResolve(abstraction, "")
}

func (c Gomodular) NamedResolve(abstraction interface{}, name string) error {
	receiverType := reflect.TypeOf(abstraction)
	if receiverType == nil {
		return errors.New("gomodular: invalid abstraction")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()

		if concrete, exist := c[elem][name]; exist {
			if instance, err := concrete.make(c); err == nil {
				reflect.ValueOf(abstraction).Elem().Set(reflect.ValueOf(instance))
				return nil
			} else {
				return err
			}
		}

		return errors.New("gomodular: no concrete found for: " + elem.String())
	}

	return errors.New("gomodular: invalid abstraction")
}

func (c Gomodular) Fill(structure interface{}) error {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		return errors.New("gomodular: invalid structure")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()
		if elem.Kind() == reflect.Struct {
			s := reflect.ValueOf(structure).Elem()

			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)

				if t, exist := s.Type().Field(i).Tag.Lookup("gomodular"); exist {
					var name string

					if t == "type" {
						name = ""
					} else if t == "name" {
						name = s.Type().Field(i).Name
					} else {
						return fmt.Errorf("gomodular: %v has an invalid struct tag", s.Type().Field(i).Name)
					}

					if concrete, exist := c[f.Type()][name]; exist {
						instance, err := concrete.make(c)
						if err != nil {
							return err
						}

						ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
						ptr.Set(reflect.ValueOf(instance))

						continue
					}

					return fmt.Errorf("gomodular: cannot make %v field", s.Type().Field(i).Name)
				}
			}

			return nil
		}
	}

	return errors.New("gomodular: invalid structure")
}

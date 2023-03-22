package gomodular

var Global = New()

func Singleton(resolver interface{}) error {
	return Global.Singleton(resolver)
}

func SingletonLazy(resolver interface{}) error {
	return Global.SingletonLazy(resolver)
}

func NamedSingleton(name string, resolver interface{}) error {
	return Global.NamedSingleton(name, resolver)
}

func NamedSingletonLazy(name string, resolver interface{}) error {
	return Global.NamedSingletonLazy(name, resolver)
}

func Transient(resolver interface{}) error {
	return Global.Transient(resolver)
}

func TransientLazy(resolver interface{}) error {
	return Global.TransientLazy(resolver)
}

func NamedTransient(name string, resolver interface{}) error {
	return Global.NamedTransient(name, resolver)
}

func NamedTransientLazy(name string, resolver interface{}) error {
	return Global.NamedTransientLazy(name, resolver)
}

func Reset() {
	Global.Reset()
}

func Call(receiver interface{}) error {
	return Global.Call(receiver)
}

func Resolve(abstraction interface{}) error {
	return Global.Resolve(abstraction)
}

func NamedResolve(abstraction interface{}, name string) error {
	return Global.NamedResolve(abstraction, name)
}

func Fill(receiver interface{}) error {
	return Global.Fill(receiver)
}

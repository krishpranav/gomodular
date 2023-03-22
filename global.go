package gomodular

var Global = New()

func Singleton(resolver interface{}) error {
	return Global.Singleton(resolver)
}

func SingletonLazy(resolver interface{}) error {
	return Global.SingletonLazy(resolver)
}

func Reset() {
	Global.Reset()
}

func Call(receiver interface{}) error {
	return Global.Call(receiver)
}

func NamedResolve(abstraction interface{}, name string) error {
	return Global.NamedResolve(abstraction, name)
}

func Fill(receiver interface{}) error {
	return Global.Fill(receiver)
}

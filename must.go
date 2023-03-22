package gomodular

func MustSingleton(c Gomodular, resolver interface{}) {
	if err := c.Singleton(resolver); err != nil {
		panic(err)
	}
}

func MustSingletonLazy(c Gomodular, resolver interface{}) {
	if err := c.SingletonLazy(resolver); err != nil {
		panic(err)
	}
}

func MustNamedSingleton(c Gomodular, name string, resolver interface{}) {
	if err := c.NamedSingleton(name, resolver); err != nil {
		panic(err)
	}
}

func MustNamedSingletonLazy(c Gomodular, name string, resolver interface{}) {
	if err := c.NamedSingletonLazy(name, resolver); err != nil {
		panic(err)
	}
}

func MustTransient(c Gomodular, resolver interface{}) {
	if err := c.Transient(resolver); err != nil {
		panic(err)
	}
}

func MustTransientLazy(c Gomodular, resolver interface{}) {
	if err := c.TransientLazy(resolver); err != nil {
		panic(err)
	}
}

func MustNamedTransient(c Gomodular, name string, resolver interface{}) {
	if err := c.NamedTransient(name, resolver); err != nil {
		panic(err)
	}
}

func MustNamedTransientLazy(c Gomodular, name string, resolver interface{}) {
	if err := c.NamedTransientLazy(name, resolver); err != nil {
		panic(err)
	}
}

func MustCall(c Gomodular, receiver interface{}) {
	if err := c.Call(receiver); err != nil {
		panic(err)
	}
}

func MustResolve(c Gomodular, abstraction interface{}) {
	if err := c.Resolve(abstraction); err != nil {
		panic(err)
	}
}

func MustNamedResolve(c Gomodular, abstraction interface{}, name string) {
	if err := c.NamedResolve(abstraction, name); err != nil {
		panic(err)
	}
}

func MustFill(c Gomodular, receiver interface{}) {
	if err := c.Fill(receiver); err != nil {
		panic(err)
	}
}

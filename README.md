## gomodular
- Golang Dependency Injection Framework. 

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

## Installing:
```
go get -u github.com/krishpranav/gomodular
```

# Tutorial:
## Singleton Typed Binding:
```golang
err := gomodular.Singleton(func() Abstraction {
  return Implementation
})

err := gomodular.Singleton(func() (Abstraction, error) {
  return Implementation, nil
})
```

- error singleton
```golang
err := gomodular.Singleton(func() Database {
  return &MySQL{}
})
```

## Transient:
```golang
err := gomodular.Transient(func() Shape {
  return &Rectangle{}
})
```

## Bindings:
```golang
err := gomodular.NamedSingleton("square", func() Shape {
    return &Rectangle{}
})
err := gomodular.NamedSingleton("rounded", func() Shape {
    return &Circle{}
})

err := gomodular.NamedTransient("sql", func() Database {
    return &MySQL{}
})
err := gomodular.NamedTransient("noSql", func() Database {
    return &MongoDB{}
})
```

## Resolver Errors:
```golang
err := gomodular.Transient(func() (Shape, error) {
  return nil, errors.New("my-app: cannot create a Shape implementation")
})
```

## Resolving:
- Resolves Dependencies such as ```Resolve()```, ```Call()``` & ```Fill()``` methods
```golang
var a Abstraction
err := gomodular.Resolve(&a)
```

- resolving using references:
```golang
var m Mailer
err := gomodular.Resolve(&m)
m.Send("example@gmail.com", "Hello World!")
```

- named resolving using references:
```golang
var s Shape
err := gomodular.NamedResolve(&s, "rounded")
```

- closures
```golang
err := gomodular.Call(func(a Abstraction) {
})
```

- resolving using closures
```golang
err := gomodular.Call(func(db Database) {
  db.Query("...")
})
```

- resolve using multiple abstractions
```golang
err := gomodular.Call(func(db Database, s Shape) {
  db.Query("...")
  s.Area()
})
```

- raising error in receiver function
```golang
err := gomodular.Call(func(db Database) error {
  return db.Ping()
})
```

## Structs
- using ```Fill()``` method in structs
```golang
type App struct {
    database Database
    other int
}

myApp := App{}

err := gomodular.Fill(&myApp)

```

- binding time
```golang
err := gomodular.Singleton(func() Config {
    return &JsonConfig{...}
})

err := gomodular.Singleton(func(c Config) Database {
    return &MySQL{
        Username: c.Get("DB_USERNAME"),
        Password: c.Get("DB_PASSWORD"),
    }
})
```

## Standalone Instance:
```golang
c := gomodular.New()

err := c.Singleton(func() Database {
    return &MySQL{}
})

err := c.Call(func(db Database) {
    db.Query("...")
})
```

## Helpers:
```golang
g := gomodular.New()

gomodular.MustSingleton(g, func() Shape {
    return &Circle{a: 13}
})

gomodular.MustCall(g, func(s Shape) {
})
```

## Lazy binding:
there are many lazy bindings some of them are
- ```gomodular.SingletonLazy()```
- ```gomodular.NamedSingletonLazy()```
- ```gomodular.TransientLazy()```
- ```gomodular.NamedTransientLazy()```

## Contributing:
- gomodular is an open-source project, this is still in development adding more dependency injection stuffs and many more features is always welcomed.

## License:
- gomodular is licensed under [MIT-License](https://github.com/krishpranav/gomodular/blob/main/LICENSE)
```
MIT License

Copyright (c) 2023 Krisna Pranav

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
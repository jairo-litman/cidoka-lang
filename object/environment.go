package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, *Environment, bool) {
	env := e
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, env, ok = e.outer.Get(name)
	}
	return obj, env, ok
}

func (e *Environment) GetNoRecursion(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func SetInEnv(name string, val Object, env *Environment) Object {
	env.store[name] = val
	return val
}

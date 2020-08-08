package colobj

// Env : container for variables' and functions' bindings
type Env struct {
	bindings    map[string]Object
	ContainedIn *Env
}

// NewEnv : constructs and returns a new Env
func NewEnv() *Env {
	bind := make(map[string]Object)
	return &Env{
		bindings:    bind,
		ContainedIn: nil,
	}
}

// NewInnerEnv : constructs a new enviroment within an existing enviroment
func NewInnerEnv(extern *Env) *Env {
	newEnv := NewEnv()
	newEnv.ContainedIn = extern
	return newEnv
}

// Get : to retrieve the value bound to a name from an env. If no binding
// is found, it checks if the current environment is contained within
// another environment. If it is, it looks for a binding to that same
// name as earlier.
func (e *Env) Get(name string) (Object, bool) {
	val, ok := e.bindings[name]
	if !ok && e.ContainedIn != nil {
		val, ok = e.ContainedIn.Get(name)
	}
	return val, ok
}

// Set : to bind a new value to a given name. If the name was
// previously bound to some value, it will be bound to the
// new value after this function is called
func (e *Env) Set(name string, val Object) Object {
	e.bindings[name] = val
	return val
}

// IsInside : To find out if the current environment is embedded in another
// environment, which happens when evaluation goes inside of functions or
// loops
func (e *Env) IsInside() bool {
	if e.ContainedIn == nil {
		return false
	}
	return true
}

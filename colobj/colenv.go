package colobj

// Env : container for variables' and functions' bindings
type Env struct {
	bindings map[string]Object
}

// NewEnv : constructs and returns a new Env
func NewEnv() *Env {
	bind := make(map[string]Object)
	return &Env{bindings: bind}
}

// Get : to retrieve the value bound to a name from the env,
// (also, if the name hasn't been bound to any value yet,
// Get() returns an err)
func (e *Env) Get(name string) (Object, bool) {
	val, ok := e.bindings[name]
	return val, ok
}

// Set : to bind a new value to a given name. If the name was
// previously bound to some value, it will be bound to the
// new value after this function is called
func (e *Env) Set(name string, val Object) Object {
	e.bindings[name] = val
	return val
}

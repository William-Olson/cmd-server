package cmddeps

import "fmt"

// T is just a blank interface type
type T interface{}

// Deps is the manager that holds the dependencies
type Deps struct {
	deps map[string]T
}

// NewDeps creates a new Deps
func NewDeps() Deps {

	return Deps{map[string]T{}}

}

// Set will create or orvewrite the dependency with given name
func (dep *Deps) Set(name string, t T) {

	fmt.Printf("Setting dependency[%s]\n", name)
	dep.deps[name] = t

}

// Get retrieves the dependency with given name
func (dep *Deps) Get(name string) T {

	return dep.deps[name]

}

// Print outputs all currently mapped dependencies to the console
func (dep *Deps) Print() {

	for k, d := range dep.deps {
		fmt.Printf("logging deps[%s] = %v\n", k, d)
	}

}

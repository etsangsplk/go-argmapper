package argmapper

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/go-argmapper/internal/graph"
)

const (
	// weightNormal is the typcal edge weight.
	weightNormal = 1

	// weightTyped is the weight to use for edges that connected to any
	// type-only vertex. We weigh these heavier since we prefer valued vertices.
	weightTyped = 5

	// weightMatchingName is the weight to use for the edges to any value
	// vertex with a matching name. This has the effect of preferring edges
	// from "A string" to "A int" for example (over "B string" to "A int"),
	// since we'd prefer to convert our original type.
	weightMatchingName = -1
)

// valueVertex represents any named and typed value.
type valueVertex struct {
	Name string
	Type reflect.Type

	Value reflect.Value
}

func (v *valueVertex) Hashcode() interface{} {
	return fmt.Sprintf("%s/%s", v.Name, v.Type.String())
}

// funcVertex is our target function. There is only ever one of these
// in the graph.
type funcVertex struct {
	Func *Func
}

func (v *funcVertex) Hashcode() interface{} { return v.Func }
func (v *funcVertex) String() string        { return "func: " + v.Func.fn.String() }

// typedArgVertex represents a typed argument to a function. These have no
// name and match any matching types.
type typedArgVertex struct {
	Name string
	Type reflect.Type

	Value reflect.Value
}

func (v *typedArgVertex) Hashcode() interface{} {
	return fmt.Sprintf("arg: %s", v.Type.String())
}

func (v *typedArgVertex) String() string { return v.Hashcode().(string) }

// typedOutputVertex represents an output from a function that is typed
// only and has no name. This can be inherited by any value with a matching
// type.
type typedOutputVertex struct {
	Name string
	Type reflect.Type

	Value reflect.Value
}

func (v *typedOutputVertex) Hashcode() interface{} {
	return fmt.Sprintf("out: %s", v.Type.String())
}

func (v *typedOutputVertex) String() string { return v.Hashcode().(string) }

// rootVertex tracks the root of a function call. This should have
// in-edges only from the inputs. We use this to get a single root.
type rootVertex struct{}

func (v *rootVertex) String() string { return "root" }

var (
	_ graph.VertexHashable = (*funcVertex)(nil)
	_ graph.VertexHashable = (*valueVertex)(nil)
	_ graph.VertexHashable = (*typedArgVertex)(nil)
	_ graph.VertexHashable = (*typedOutputVertex)(nil)
)

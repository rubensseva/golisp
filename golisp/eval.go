package golisp

import (
	"fmt"
)

func Eval(node Node) Node {
	switch t := node.Data.(type) {
	case []Node:
		// Get a reference to the function symbol
		// In lisp this can be dynamically found so
		// we need to eval to get it
		f := Eval(t[0])
		if f.Name == "" {
			panic(fmt.Sprintf("expected symbol, got: %+v", f))
		}

		fn, ok := fnMap()[f.Name]
		if !ok {
			return dynamicBuiltin(f.Name, t[1:]...)
		}

		finalParams := t[1:]
		if !fn.IsSpecial {
			evaldParams := []Node{}
			for _, p := range finalParams {
				evaldParams = append(evaldParams, Eval(p))
			}
			finalParams = evaldParams
		}

		return fn.FN(finalParams...)

	default:
		// If its a int, string or other primitive, we
		// just return it
		return node
	}
}

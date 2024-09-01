package golisp

import (
	"fmt"
)

func Eval(node Node) Node {
	switch node.Type {
	case TypeList:
		// Get a refernce to the function symbol
		// In lisp this can be dynamically found so
		// we need to eval to get it
		f := Eval(node.Nested[0])
		if f.Type != TypeSymbol {
			panic(fmt.Sprintf("expected symbol, got: %+v", f))
		}

		fn, ok := fnMap()[f.Name]
		if !ok {
			panic(fmt.Sprintf("unknown symbol: %v", fn))
		}

		finalParams := node.Nested[1:]
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

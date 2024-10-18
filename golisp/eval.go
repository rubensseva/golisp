package golisp

import (
	"fmt"
	"log"
)

func Eval(node any) any {
	switch t := node.(type) {
	case LiteralMap:
			if len(t)%2 != 0 {
				log.Fatalf("map did not have even key-val pairs: %v, %+v", len(t), t)
			}

			newMap := map[any]any{}
			for i := 0; i < len(t); i += 2 {
				key := Eval(t[0])
				val := Eval(t[1])
				newMap[key] = val
			}

		    return newMap
	case LiteralList:
		// Get a reference to the function symbol
		// In lisp this can be dynamically found so
		// we need to eval to get it
		f := Eval(t[0])
		sym, ok := f.(Symbol)
		if !ok {
			panic(fmt.Sprintf("expected symbol, got: %+v", f))
		}

		fn, ok := fnMap()[string(sym)]
		if !ok {
			return dynamicBuiltin(string(sym), t[1:]...)
		}

		finalParams := t[1:]
		if !fn.IsSpecial {
			evaldParams := []any{}
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

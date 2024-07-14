package parser

import (
	"fmt"
)

type Op func(n ...Node) Node

func BuiltinPlus(nodes ...Node) Node {
	sum := int64(0)
	for _, n := range nodes {
		num, ok := n.Data.(int64)
		if !ok {
			panic(fmt.Sprintf("plus function got something that was not an int: %+v, %T", n, n.Data))
		}
		sum += num
	}
	return Node{
		Type:   TypeInt,
		Data:   sum,
	}
}

var fnMap = map[string]Op{
	"+": BuiltinPlus}

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

		fn, ok := fnMap[f.Name]
		if !ok {
			panic(fmt.Sprintf("unknown symbol: %v", fn))
		}

		evaldParams := []Node{}
		for _, p := range node.Nested[1:] {
			evaldParams = append(evaldParams, Eval(p))
		}

		return fn(evaldParams...)

	default:
		// If its a int, string or other primitive, we
		// just return it
		return node

	}


}

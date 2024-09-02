package golisp

import (
	"fmt"
)

type OpFn func(n ...Node) Node
type Builtin struct {
	FN        OpFn
	IsSpecial bool
}

func fnMap() map[string]Builtin {
	return map[string]Builtin{
		"+":       {BuiltinPlus, false},
		"-":       {BuiltinMinus, false},
		"cond":    {BuiltinCond, true},
		"=":       {BuiltinEq, false},
		"println": {BuiltinPrintln, false},
	}
}

func NodeField[T any](n Node) T {
	val, ok := n.Data.(T)
	if !ok {
		panic(fmt.Sprintf("tried to convert node data %+v, but it was %T", n, n.Data))
	}
	return val
}

func IsTrue(n Node) bool {
	switch t := n.Data.(type) {
	case []Node:
		return len(t) != 0
	case int, int32, int64:
		return NodeField[int64](n) != 0
	case string:
		return NodeField[string](n) != ""
	case bool:
		return NodeField[bool](n)
	default:
		panic(fmt.Sprintf("unrecognized type: %s", n))
	}
}

func Equal(nodes ...Node) bool {
	prevVal := nodes[0].Data
	for _, node := range nodes[1:] {
		if prevVal != node.Data {
			return false
		}
		prevVal = node.Data
	}
	return true
}

func BuiltinPlus(nodes ...Node) Node {
	if len(nodes) == 0 {
		return Node{
			Data: int64(0),
		}
	}
	sum := int64(0)
	for _, n := range nodes {
		num := NodeField[int64](n)
		sum += num
	}
	return Node{
		Data: sum,
	}
}

func BuiltinMinus(nodes ...Node) Node {
	if len(nodes) == 0 {
		return Node{
			Data: int64(0),
		}
	}
	res := NodeField[int64](nodes[0])
	for _, n := range nodes[1:] {
		num := NodeField[int64](n)
		res -= num
	}
	return Node{
		Data: res,
	}
}

func BuiltinEq(nodes ...Node) Node {
	return Node{
		Data: Equal(nodes...),
	}
}

func BuiltinCond(nodes ...Node) Node {
	if len(nodes)%2 != 0 {
		panic("cond expects even number of arguments")
	}

	for i := 0; i < len(nodes); i += 2 {
		evald := Eval(nodes[i])
		if IsTrue(evald) {
			return Eval(nodes[i+1])
		}
	}

	panic(fmt.Sprintf("no conditions of cond was hit: %+v", nodes))
}

func BuiltinPrintln(nodes ...Node) Node {
	for _, n := range nodes {
		n.NodePprint()
	}
	fmt.Println()
	return Node{}
}

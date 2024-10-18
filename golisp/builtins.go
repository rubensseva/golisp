package golisp

import (
	"fmt"
	"reflect"
)

type OpFn func(n ...any) any
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

func NodeField[T any](n any) T {
	val, ok := n.(T)
	if !ok {
		panic(fmt.Sprintf("tried to convert node data %+v, but it was %T", n, n))
	}
	return val
}

func IsTrue(n any) bool {
	switch n.(type) {
	case int, int32, int64:
		return NodeField[int64](n) != 0
	case string:
		return NodeField[string](n) != ""
	case bool:
		return NodeField[bool](n)
	default:
		v := reflect.ValueOf(n)
		if v.Kind() == reflect.Slice {
			return v.Len() != 0 // Get the length of the slice using reflection
		}
		panic(fmt.Sprintf("unrecognized type: %s", n))
	}
}

func Equal(nodes ...any) any {
	prevVal := nodes[0]
	for _, node := range nodes[1:] {
		if prevVal != node {
			return false
		}
		prevVal = node
	}
	return true
}

func BuiltinPlus(nodes ...any) any {
	if len(nodes) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nodes {
		num := NodeField[int](n)
		sum += num
	}
	return sum
}

func BuiltinMinus(nodes ...any) any {
	if len(nodes) == 0 {
		return 0
	}
	res := NodeField[int](nodes[0])
	for _, n := range nodes[1:] {
		num := NodeField[int](n)
		res -= num
	}
	return res
}

func BuiltinEq(nodes ...any) any {
	return Equal(nodes...)
}

func BuiltinCond(nodes ...any) any {
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

func BuiltinPrintln(nodes ...any) any {
	for _, n := range nodes {
		NodePprint(n)
	}
	fmt.Println()
	return nil
}

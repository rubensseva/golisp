package golisp

import (
	"log"
	"reflect"
	"strings"
)

var dynamicFnMap = map[string]any{
	"strings.split":    strings.Split,
	"strings.split-n":  strings.SplitN,
	"strings.count":    strings.Count,
	"strings.replace":  strings.Replace,
	"strings.to-lower": strings.ToLower,
	"strings.to-upper": strings.ToUpper,
	"strings.trim":     strings.Trim,
}

func dynamicBuiltin(name string, n ...Node) Node {
	fn, ok := dynamicFnMap[name]
	if !ok {
		log.Fatalf("no dynamic builtin with name: %v", name)
	}

	// typ := reflect.TypeOf(fn)
	val := reflect.ValueOf(fn)

	args := []reflect.Value{}

	for _, nn := range n {
		args = append(args, reflect.ValueOf(nn.Data))
	}

	log.Printf("calling: %+v, with args: %+v", fn, n)
	rres := val.Call(args)

	res := []Node{}
	for _, r := range rres {
		res = append(res, Node{
			Type: TypeAny,
			// Name:   name,
			Data: r.Interface(),
		})
	}

	return Node{
		Type:   TypeList,
		Nested: res,
	}
}

func strings_split(n ...Node) Node {
	if len(n) != 2 {
		panic("split expects two strings")
	}
	for _, nn := range n {
		if nn.Type != TypeString {
			panic("split expects all arguments to be strings")
		}
	}
	s := n[0].Data.(string)
	sep := n[1].Data.(string)

	strs := strings.Split(s, sep)

	res := []Node{}
	for _, s := range strs {
		res = append(res, Node{
			Type: TypeString,
			// Name:   "",
			Data: s,
			// Nested: n,
		})
	}
	return Node{
		Type:   TypeList,
		Nested: res,
	}
}

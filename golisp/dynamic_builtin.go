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

	"strings.testfoo": TestFoo,
}

var dynamicStructMap = map[string]any{
	"strings.foo": Foo{},
}

type Foo struct {
	Bar string
}

func TestFoo(f Foo) string {
	return f.Bar
}

func mapToStruct(m map[any]any, structInstance any) {
	// Get the reflect.Type of the expected struct
	structType := reflect.TypeOf(structInstance)

	// Create a new instance of the struct
	newStruct := reflect.New(structType).Elem()

	// Iterate over the fields of the struct
	for i := 0; i < newStruct.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name

		// Check if the map has a value for this field
		if value, ok := m[fieldName]; ok {
			structField := newStruct.FieldByName(fieldName)

			// Ensure the value from the map is assignable to the struct field
			val := reflect.ValueOf(value)
			if val.Type().AssignableTo(structField.Type()) {
				structField.Set(val)
			}
		}
	}
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
			// Name:   name,
			Data: r.Interface(),
		})
	}

	return Node{
		Data: res,
	}
}

func strings_split(n ...Node) Node {
	if len(n) != 2 {
		panic("split expects two strings")
	}
	for _, nn := range n {
		_, ok := nn.Data.(string)
		if !ok {
			panic("split expects all arguments to be strings")
		}
	}
	s := n[0].Data.(string)
	sep := n[1].Data.(string)

	strs := strings.Split(s, sep)

	res := []Node{}
	for _, s := range strs {
		res = append(res, Node{
			// Name:   "",
			Data: s,
			// Nested: n,
		})
	}
	return Node{
		Nested: res,
	}
}

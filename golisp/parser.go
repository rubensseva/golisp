package golisp

import (
	"fmt"
	// "golisp/parser"
	"strconv"
	"strings"
)

type Type string

const (
	TypeList   = Type("list")
	TypeString = Type("string")
	TypeInt    = Type("int")
	TypeSymbol = Type("symbol")
	TypeBool   = Type("bool")
	TypeAny    = Type("any")
)

type Node struct {
	Type Type
	// Only present on functions and symbols
	Name string

	Data any

	// Nested represents list data. Lists can represent function invocation, in
	// that case the first element is the function name, and the rest are
	// function params
	Nested []Node
}

func (n Node) NodePprint() {
	switch n.Type {
	case TypeList:
		fmt.Print("(")
		for i, n := range n.Nested {
			if i != 0 {
				fmt.Printf(" ")
			}
			n.NodePprint()
		}
		fmt.Print(")")
	case TypeSymbol:
		fmt.Print(n.Name)
	default:
		fmt.Print(n.Data)
	}
}

func IsStr(str string) bool {
	return strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"")
}

func IsInt(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

func expectFound(found bool) {
	if !found {
		panic("expected more tokens")
	}
}

func Parse(tokenizer *Tokenizer) Node {
	for {
		token := tokenizer.Token()
		switch {
		case token == "(":
			elements := []Node{}
			for {
				val := tokenizer.Peek()
				if val == ")" {
					tokenizer.Token() // purge the )
					break
				}

				newEl := Parse(tokenizer)

				elements = append(elements, newEl)
			}

			return Node{
				Type:   TypeList,
				Nested: elements,
			}

		case token == ")":
			panic("should never happen")

		case IsStr(token):
			return Node{
				Type: TypeString,
				Data: strings.Trim(token, "\""),
			}

		case IsInt(token):
			n, _ := strconv.ParseInt(token, 10, 64)
			return Node{
				Type: TypeInt,
				Data: n,
			}

		case token == "true" || token == "false":
			if token == "true" {
				return Node{
					Type: TypeBool,
					Data: true,
				}
			}
			return Node{
				Type: TypeBool,
				Data: false,
			}

		default:
			return Node{
				Type: TypeSymbol,
				Name: token,
			}
		}
	}

	panic("should never get here")
}

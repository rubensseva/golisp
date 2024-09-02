package golisp

import (
	"fmt"
	// "golisp/parser"
	"strconv"
	"strings"
)

type Type string

type Node struct {
	// Only present on functions and symbols
	Name string

	// Data contains the actual data, and type information in the form of
	// the underlying Golang type.
	// If its a []Node it represents a golisp lisp. Anything else is just a
	// literal value.
	// Lists ([]Node) can represent function invocation, in
	// that case the first element is the function name, and the rest are
	// function params.
	Data any
}

func (n Node) NodePprint() {
	if n.Name != "" {
		fmt.Print(n.Name)
	} else {
		switch t := n.Data.(type) {
		case []Node:
			fmt.Print("(")
			for i, n := range t {
				if i != 0 {
					fmt.Printf(" ")
				}
				n.NodePprint()
			}
			fmt.Print(")")
		default:
			fmt.Print(n.Data)
		}
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
				Data: elements,
			}

		case token == ")":
			panic("should never happen")

		case IsStr(token):
			return Node{
				Data: strings.Trim(token, "\""),
			}

		case IsInt(token):
			n, _ := strconv.ParseInt(token, 10, 64)
			return Node{
				Data: n,
			}

		case token == "true" || token == "false":
			if token == "true" {
				return Node{
					Data: true,
				}
			}
			return Node{
				Data: false,
			}

		default:
			return Node{
				Name: token,
			}
		}
	}

	panic("should never get here")
}

package parser

import (
	"fmt"
	// "golisp/parser"
	"log"
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
)

type Node struct {
	Type Type
	Name string
	// Data contains the data AND the type
	Data any

	// Nested can be func params, or func invocation arguments
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
		token, err := tokenizer.Token()
		if err != nil {
			log.Fatalf("reading token: %v", err)
		}
		switch {
		case token == "(":
			elements := []Node{}
			for {
				val, err := tokenizer.Peek()
				if err != nil {
					log.Fatalf("peeking: %v", err)
				}
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

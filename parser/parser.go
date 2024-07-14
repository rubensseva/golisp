package parser

import (
	"fmt"
	"golisp/channel"
	"strconv"
	"strings"
)

type Type string

const (
	TypeFuncInvocation = Type("funcinvocation")
	TypeString = Type("string")
	TypeInt = Type("int")
	TypeSymbol = Type("symbol")
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
	case TypeFuncInvocation:
		fmt.Printf("(%s ", n.Name)
		for i, n := range n.Nested {
			if i != 0 {
				fmt.Printf(" ")
			}
			n.NodePprint()
		}
		fmt.Printf(")")
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

func Parse(tokens *channel.PeekableChannel) (Node) {
	for {
		token, found := tokens.Receive()
		if !found {
			break
		}
		fmt.Println("handling token ", token)
		switch  {
		case token == "(":
			funcName, found := tokens.Receive()
			expectFound(found)
			fmt.Println("detected func", funcName)

			params := []Node{}
			for {
				val, found := tokens.Peek()
				expectFound(found)
				fmt.Println("peeked val", val)
				if val == ")" {
					break
				}

				{
				val, found := tokens.Peek()
				fmt.Println("peeked val", val)
				expectFound(found)
				}

				newParam := Parse(tokens)

				params = append(params, newParam)
			}

			return Node{
				Type:     TypeFuncInvocation,
				Name: funcName,
				Data:     nil,
				Nested:   params,
			}

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

		default:
			return Node{
				Type:     TypeSymbol,
				Name: token,
			}
		}
	}

	panic("should never get here")
}

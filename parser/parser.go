package parser

import (
	"fmt"
	"golisp/channel"
	"strconv"
	"strings"
)

type Type string

const (
	TypeList = Type("list")
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

func Parse(tokens *channel.PeekableChannel) (Node) {
	for {
		token, found := tokens.Receive()
		if !found {
			break
		}
		fmt.Println("handling token ", token)
		switch  {
		case token == "(":
			fmt.Println("STARTING LIST HANDLING")
			elements := []Node{}
			for {
				val, found := tokens.Peek()
				expectFound(found)
				fmt.Println("got val", val)
				if val == ")" {
					fmt.Println("detected ), breaking")
					tokens.Receive() // purge the )
					break
				}

				newEl := Parse(tokens)

				elements = append(elements, newEl)
			}

			fmt.Println("ENDING LIST HANDLING")
			return Node{
				Type:     TypeList,
				Nested:   elements,
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

		default:
			return Node{
				Type:     TypeSymbol,
				Name: token,
			}
		}
	}

	panic("should never get here")
}

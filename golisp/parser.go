package golisp

import (
	"fmt"
	"log"
	// "golisp/parser"
	"strconv"
	"strings"
)

type Type string

type LiteralMap []any
type LiteralList []any
type Symbol string

// type Primitive any

func NodePprint(n any) {
		switch t := n.(type) {
		case LiteralList:
			fmt.Print("(")
			for i, n := range t {
				if i != 0 {
					fmt.Printf(" ")
				}
				NodePprint(n)
			}
			fmt.Print(")")
		case LiteralMap:
			fmt.Print("(")
			for i, n := range t {
				if i != 0 {
					fmt.Printf(" ")
				}
				NodePprint(n)
			}
			fmt.Print(")")
		default:
			fmt.Print(t)
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

func Parse(tokenizer *Tokenizer) any {
	for {
		token := tokenizer.Token()
		switch {
		case token == "(" || token == "{":
			isMap := false
			if token == "{" {
				isMap = true
			}

			elements := LiteralList{}
			for {
				val := tokenizer.Peek()
				if val == ")" || val == "}" {
					if isMap && val == ")" {
						log.Fatalln("got ending of list wihle expecting ending of map")
					}
					tokenizer.Token() // purge the )
					break
				}

				newEl := Parse(tokenizer)

				elements = append(elements, newEl)
			}

			if isMap {
				return LiteralMap(elements)
			}
			return elements

		case token == ")" || token == "}":
			panic("should never happen")

		case IsStr(token):
		    return strings.Trim(token, "\"")

		case IsInt(token):
			n, err := strconv.ParseInt(token, 10, 64)
			if err != nil {
				panic(err)
			}
			return n

		case token == "true" || token == "false":
			return token == "true"

		default:
			return Symbol(token)
		}
	}

	panic("should never get here")
}

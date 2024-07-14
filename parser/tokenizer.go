package parser

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"golisp/channel"
)

var delims = []rune{'(', ')'}

func tokenize(str string, cout *channel.PeekableChannel) {
	str = strings.Trim(str, " ")

	prevI := 0
	skipSpace := false

	for i, r := range str {
		// If we previously encountered whitespace, skip
		// until we dont encounter more whitespace
		if skipSpace {
			if r == ' ' {
				continue
			}
			prevI = i
			skipSpace = false
		}

		if slices.Contains([]rune{')', '('}, r) {
			if r == ')' {
				skipSpace = true
			}
			if i > prevI {
				cout.Send(str[prevI:i])
			}
			cout.Send(string(r))
			prevI = i + 1
			continue
		}
		if r == ' ' {
			cout.Send(str[prevI:i])
			skipSpace = true
			continue
		}
	}
}



func Tokenize(in io.Reader, cout *channel.PeekableChannel) (error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("reading all: %w", err)
	}

	tokenize(string(data), cout)
	return nil
}


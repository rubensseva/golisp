package golisp

import (
	"io"
)



type Tokenizer struct {
	buffer []rune
	r      io.RuneReader
	tokenBuffer []string
}

func NewTokenizer(r io.RuneReader) *Tokenizer {
	return &Tokenizer{
		buffer:      []rune{},
		r:           r,
		tokenBuffer: []string{},
	}
}

func (t *Tokenizer) peekRune() rune {
	if len(t.buffer) > 0 {
		return t.buffer[0]
	}
	r, _, err := t.r.ReadRune()
	if err != nil {
		panic(err)
	}
	t.buffer = []rune{r}
	return r
}

func (t *Tokenizer) readRune() rune {
	r := t.peekRune()
	t.buffer = t.buffer[1:]
	return r
}

func (t *Tokenizer) Token() string {

	if len(t.tokenBuffer) > 0 {
		token := t.tokenBuffer[0]
		t.tokenBuffer = t.tokenBuffer[1:]
		return token
	}

	atStart := true
	isString := false
	var token []rune
	for {
		r := t.peekRune()
		// Trim spaces and newlines
		if atStart && r == ' ' || r == '\n' {
			t.readRune()
			continue
		}



		// If the first thing we encounter (except for spaces) is a parenthesis,
		// we consume and return it
		if atStart && (r == '(' || r == ')') {
			t.readRune()
			return string(r)
		}

		// At this point we are no longer at the beginning
		atStart = false

		// Now we need to check if we are consuming a string token, in which
		// case we should ignore delimiters inside the string token
		if r == '"' && !isString {
			isString = true
			token = append(token, r)
			t.readRune()
			continue
		}
		if isString {
			token = append(token, r)
			// If we are inside the string, we need to consume a rune to get to
			// the next one on the next iteration. If we are at the end of a
			// string, meaning the peeked rune is == '"', then we also should
			// consume a rune so that we dont encounter the " rune on the next
			// run of the tokenizer
			t.readRune()
			if r == '"' {
				return string(token)
			}
			continue
		}

		// Check for delimiters
		if r == '(' || r == ')' || r == ' ' || r == '\n' {
			return string(token)
		}

		token = append(token, r)
		t.readRune()
	}
}

func (t *Tokenizer) Peek() string {
	if len(t.tokenBuffer) > 0 {
		return t.tokenBuffer[0]
	}

	token := t.Token()
	t.tokenBuffer = []string{token}
	return token
}

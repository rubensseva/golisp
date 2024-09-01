package golisp

import (
	"fmt"
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

func (t *Tokenizer) peekRune() (rune, error) {
	var z rune

	if len(t.buffer) > 0 {
		return t.buffer[0], nil
	}
	r, _, err := t.r.ReadRune()
	if err != nil {
		return z, fmt.Errorf("reading rune: %w", err)
	}
	t.buffer = []rune{r}
	return r, nil
}

func (t *Tokenizer) readRune() (rune, error) {
	var z rune

	r, err := t.peekRune()
	if err != nil {
		return z, fmt.Errorf("peeking one: %w", err)
	}
	t.buffer = t.buffer[1:]
	return r, nil
}

func (t *Tokenizer) Token() (string, error) {
	var z string

	if len(t.tokenBuffer) > 0 {
		token := t.tokenBuffer[0]
		t.tokenBuffer = t.tokenBuffer[1:]
		return token, nil
	}

	atStart := true
	var token []rune
	for {
		r, err := t.peekRune()
		if err != nil {
			return z, fmt.Errorf("reading one: %w", err)
		}
		// Trim spaces and newlines
		if atStart && r == ' ' || r == '\n' {
			if _, err := t.readRune(); err != nil {
				return z, fmt.Errorf("trimming: %w", err)
			}
			continue
		}

		// If the first thing we encounter (except for spaces) is a parenthesis,
		// we consume and return it
		if atStart && (r == '(' || r == ')') {
			if _, err := t.readRune(); err != nil {
				return z, fmt.Errorf("discarding a parenthesis: %w", err)
			}
			return string(r), nil
		}

		// At this point we are no longer at the beginning
		atStart = false

		// This will never trigger at the start, because we already check for
		// spaces, newlines and parenthesis. So if this condition passes, token
		// is not empty
		if r == '(' || r == ')' || r == ' ' || r == '\n' {
			return string(token), nil
		}

		token = append(token, r)
		if _, err := t.readRune(); err != nil {
			return z, fmt.Errorf("discarding a rune: %w", err)
		}
	}
}

func (t *Tokenizer) Peek() (string, error) {
	var z string

	if len(t.tokenBuffer) > 0 {
		return t.tokenBuffer[0], nil
	}

	token, err := t.Token()
	if err != nil {
		return z, fmt.Errorf("reading a token: %w", err)
	}
	t.tokenBuffer = []string{token}
	return token, nil
}

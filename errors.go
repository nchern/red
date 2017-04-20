package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errNotARequestString = errors.New("Not a request string")
	errIndentFailed      = errors.New("Failed to format json")
	errBadRequestString  = errors.New("Bad formatted request string. Must be in a form: <METHOD> /<index>/<action>")
)

func findPrevLineIndex(s string, offset int) int {
	prefixIdx := strings.LastIndex(s[:offset], "\n")
	if prefixIdx < 0 {
		prefixIdx = 0
	}
	return prefixIdx
}

type ParseError struct {
	n        int
	msg      string
	filename string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:#%d: %s", e.filename, e.n, e.msg)
}

type JsonifyError struct {
	Inner  error
	Source string
}

func (e *JsonifyError) Error() string { return e.Inner.Error() }

func (e *JsonifyError) Highlighted(offset int64) string {
	prefixIdx := findPrevLineIndex(e.Source, int(offset))
	prefixIdx = findPrevLineIndex(e.Source[:prefixIdx], prefixIdx)
	return e.Source[prefixIdx:]
}

package app

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrFormatFailed indicates that JSON chunk can not be formatted
	ErrFormatFailed = errors.New("Failed to format json")

	errNotARequestString = errors.New("Not a request string")
	errBadRequestString  = errors.New("Bad formatted request string. Must be in a form: <METHOD> /<index>/<action>")
)

func findPrevLineIndex(s string, offset int) int {
	prefixIdx := strings.LastIndex(s[:offset], "\n")
	if prefixIdx < 0 {
		prefixIdx = 0
	}
	return prefixIdx
}

type parseError struct {
	n        int
	msg      string
	filename string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("%s:#%d: %s", e.filename, e.n, e.msg)
}

// JsonifyError describes the invalid request body that must be a valid json
type JsonifyError struct {
	Inner  error
	Source string
}

func (e *JsonifyError) Error() string { return e.Inner.Error() }

// Highlighted returns the problematic json part that caused the error
func (e *JsonifyError) Highlighted(offset int64) string {
	prefixIdx := findPrevLineIndex(e.Source, int(offset))
	prefixIdx = findPrevLineIndex(e.Source[:prefixIdx], prefixIdx)
	return e.Source[prefixIdx:]
}

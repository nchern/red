package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	terminateParseToken = ">>>EOF<<<"
)

var errTimeout = errors.New("timeout")

func tryParseStdinAsync() (*ParsedRequest, error) {
	var err error
	var selection *ParsedRequest
	var buf bytes.Buffer

	tee := io.TeeReader(os.Stdin, &buf)

	finished := make(chan bool)

	go func() {
		selection, err = parseScript(tee)
		finished <- true
	}()

	select {
	case <-finished:
		// Mirror stdin to stout - this allows processing selections in vim correctly
		os.Stdout.Write(buf.Bytes())
	case <-time.After(50 * time.Millisecond):
		return nil, errTimeout
	}

	return selection, err
}

func tryParseRequestString(line string, req *ParsedRequest) error {
	for _, method := range methods {
		if !strings.HasPrefix(line, method) {
			continue
		}
		req.Method = method
		req.URI = strings.TrimSpace(strings.TrimPrefix(line, method))
		if req.URI == "" {
			return errBadRequestString
		}
		return nil
	}
	return errNotARequestString
}

func parseScript(reader io.Reader) (*ParsedRequest, error) {
	result := &ParsedRequest{}

	i := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		i++
		line := scanner.Text()
		normalizedLine := strings.TrimSpace(line)
		if normalizedLine == "" {
			continue
		}
		if strings.HasPrefix(normalizedLine, "#") {
			continue
		}
		if strings.HasPrefix(normalizedLine, "@") {
			result.Host = normalizedLine[1:]
			continue
		}
		if normalizedLine == terminateParseToken {
			break
		}
		if err := tryParseRequestString(normalizedLine, result); err != nil {
			if err != errNotARequestString {
				return nil, &ParseError{n: i, msg: err.Error(), filename: queryFilePath}
			}
		} else {
			continue
		}
		result.bodyLines = append(result.bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type ParsedRequest struct {
	Host   string
	Method string
	URI    string

	bodyLines []string
}

func NewParsedRequest() *ParsedRequest {
	return &ParsedRequest{Host: defaultHost}
}

func (r *ParsedRequest) RawBody() string {
	return strings.Join(r.bodyLines, "\n")
}

func (r *ParsedRequest) Url() string {
	url := r.Host + r.URI
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	return url
}

/*
func (r *ParsedRequest) Merge(src *ParsedRequest) *ParsedRequest {
	r.Host = notEmpty(src.Host, r.Host)
	r.Method = notEmpty(src.Method, r.Method)
	r.URI = notEmpty(src.URI, r.URI)

	if len(src.bodyLines) > 0 {
		r.bodyLines = src.bodyLines
	}
	return r
}
*/

func (r *ParsedRequest) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("Host is empty")
	}
	if r.Method == "" {
		return fmt.Errorf("Method is empty")
	}
	if r.URI == "" {
		return fmt.Errorf("Uri is empty")
	}
	return nil
}

func (r *ParsedRequest) JSON() (string, error) {
	src := r.RawBody()
	if src == "" {
		return "", nil
	}
	var obj struct{}
	if err := json.Unmarshal([]byte(src), &obj); err != nil {
		return "", &JsonifyError{Source: src, Inner: err}
	}

	return src, nil
}

package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	terminateParseToken = ">>>EOF<<<"

	// TemplateAsset sets path to template text bin asset
	TemplateAsset = "assets/template.txt"

	parseTimeout = 100 * time.Millisecond
)

var (
	errTimeout = errors.New("timeout")
	methods    = []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "HEAD"}
)

// HTTPRequest represents the request to be made
type HTTPRequest struct {
	Host   string
	Method string
	URI    string

	Headers http.Header

	bodyLines []string
}

func newHTTPRequest() *HTTPRequest {
	return &HTTPRequest{Headers: http.Header{}}
}

// CopyBodyFrom copies body from src
func (r *HTTPRequest) CopyBodyFrom(src *HTTPRequest) {
	r.bodyLines = src.bodyLines
}

func (r *HTTPRequest) rawBody() string {
	return strings.Join(r.bodyLines, "\n")
}

// URL returns url
func (r *HTTPRequest) URL() string {
	url := r.Host + r.URI
	if toks := strings.Split(url, ":"); len(toks) < 2 || !strings.HasPrefix(toks[0], "http") {
		url = "http://" + url
	}
	return url
}

// Validate validates the request
func (r *HTTPRequest) Validate() error {
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

// JSON returns body as JSON
func (r *HTTPRequest) JSON() (string, error) {
	src := r.rawBody()
	if src == "" {
		return "", nil
	}
	var obj struct{}
	if err := json.Unmarshal([]byte(src), &obj); err != nil {
		// try if it is a json array
		var obj []struct{}
		if err := json.Unmarshal([]byte(src), &obj); err == nil {
			return src, nil
		}
		return "", &JsonifyError{Source: src, Inner: err}
	}

	return src, nil
}

func (r *HTTPRequest) tryParseRequestString(line string) error {
	for _, method := range methods {
		if !strings.HasPrefix(line, method) {
			continue
		}
		r.Method = method
		r.URI = strings.TrimSpace(strings.TrimPrefix(line, method))
		if r.URI == "" {
			return errBadRequestString
		}
		return nil
	}
	return errNotARequestString
}

// TryParseAsync tries to parse given reader and times out if it can't read on time
func TryParseAsync(src io.Reader) (*HTTPRequest, error) {
	var err error
	var selection *HTTPRequest

	finished := make(chan bool)

	go func() {
		selection, err = ParseRequest(src)
		finished <- true
	}()

	select {
	case <-finished:
	case <-time.After(parseTimeout):
		return nil, errTimeout
	}

	return selection, err
}

// ParseRequest parses request info from the given reader
func ParseRequest(reader io.Reader) (*HTTPRequest, error) {
	result := newHTTPRequest()

	i := 0
	isBodyStarted := false
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		i++
		line := scanner.Text()
		normalizedLine := strings.TrimSpace(line)
		if normalizedLine == terminateParseToken {
			break
		}
		if !isBodyStarted {
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
			if strings.Index(normalizedLine, ":") > -1 {
				// parse header
				tokens := strings.SplitN(normalizedLine, ":", 2)
				tokens[0], tokens[1] = strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
				result.Headers.Add(tokens[0], tokens[1])
				continue
			}
			if err := result.tryParseRequestString(normalizedLine); err != nil {
				if err != errNotARequestString {
					//TODO: filename: queryFilePath
					return nil, &parseError{n: i, msg: err.Error()}
				}
			} else {
				continue
			}
		}
		result.bodyLines = append(result.bodyLines, line)
		isBodyStarted = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

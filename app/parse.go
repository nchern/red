package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	errTimeout  = errors.New("timeout")
	defaultHost = "http://localhost:9200"
	methods     = []string{"GET", "POST", "DELETE", "PUT"}
)

// HTTPRequest represents the request to be made
type HTTPRequest struct {
	Host   string
	Method string
	URI    string

	bodyLines []string
}

func newHTTPRequest() *HTTPRequest {
	return &HTTPRequest{Host: defaultHost}
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

// TryParseAsync parses stdin and times out if stdin is open and empty
func TryParseAsync(src io.Reader, out io.Writer) (*HTTPRequest, error) {
	var err error
	var selection *HTTPRequest
	var buf bytes.Buffer

	tee := io.TeeReader(src, &buf)

	finished := make(chan bool)

	go func() {
		selection, err = ParseRequest(tee)
		finished <- true
	}()

	select {
	case <-finished:
		// Mirror stdin to stout - this allows processing selections in vim correctly
		out.Write(buf.Bytes())
	case <-time.After(parseTimeout):
		return nil, errTimeout
	}

	return selection, err
}

func tryParseRequestString(line string, req *HTTPRequest) error {
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

// ParseRequest parses request info from the given reader
func ParseRequest(reader io.Reader) (*HTTPRequest, error) {
	result := newHTTPRequest()

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
				//TODO: filename: queryFilePath
				return nil, &parseError{n: i, msg: err.Error()}
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

package app

import (
	"bytes"
	"reflect"
	"regexp"
	"testing"
)

const (
	expectedBody = `{"size":10,"query":{"term":{"field":"value"}}}`

	srcUnderTest = `
@localhost:9200

# HELP:
# '@<hostname>:9200' at the begining of a line sets server host:port to make request to
# <METHOD> /foo/bar at the begining of a line sets request METHOD/URI a for the query
# >>>EOF<<< at the begining of a line stops parsing the source


# Example of Elasticsearch query

GET /foo/_search
{
    "size": 10,
    "query": {
        "term": {
            "field": "value"
        }
    }
}`
)

func TestParse(t *testing.T) {
	result, err := ParseRequest(bytes.NewBufferString(srcUnderTest))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}

	if err := result.Validate(); err != nil {
		t.Errorf("result.Validate returned %s", err)
	}

	if result.Host != "localhost:9200" {
		t.Errorf("Incorrect host: '%s'", result.Host)
	}

	if result.Method != "GET" {
		t.Errorf("Incorrect method: '%s'", result.Method)
	}

	if result.URI != "/foo/_search" {
		t.Errorf("Incorrect URI: '%s'", result.URI)
	}
	if result.URL() != "http://localhost:9200/foo/_search" {
		t.Errorf("Incorrect Url(): '%s'", result.URL())
	}

	json, err := result.JSON()
	if err != nil {
		t.Errorf("Malformed json: %s", err)
	}
	rx := regexp.MustCompile(`\s`)
	json = rx.ReplaceAllString(json, "")
	if json != expectedBody {
		t.Errorf("Unexpected json: [%s]", json)
	}
}

func TestParseWithHeaders(t *testing.T) {
	srcWithHeaders := `
POST /bar/_search
X-Request-ID: foobar
Accept-Encoding: en-US

{
    "size": 10,
}`
	result, err := ParseRequest(bytes.NewBufferString(srcWithHeaders))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}

	if result.Host != "" {
		t.Errorf("Incorrect host: '%s'", result.Host)
	}

	if result.Method != "POST" {
		t.Errorf("Incorrect method: '%s'", result.Method)
	}
	if result.URI != "/bar/_search" {
		t.Errorf("Incorrect URI: '%s'", result.URI)
	}

	if result.Headers.Get("X-Request-ID") != "foobar" {
		t.Errorf("Bad or unexpected headers: '%+v'", result.Headers)
	}
	if result.Headers.Get("Accept-Encoding") != "en-US" {
		t.Errorf("Bad or unexpected headers: '%+v'", result.Headers)
	}
	if len(result.Headers) != 2 {
		t.Errorf("Bad or unexpected headers: '%+v'", result.Headers)
	}
}

func TestUrlScheme(t *testing.T) {
	var req HTTPRequest

	var tests = []struct {
		expected string
		given    string
	}{
		{"http://localhost", "localhost"},
		{"http://localhost:8080", "localhost:8080"},
		{"http://localhost:8080", "http://localhost:8080"},
		{"http://localhost", "http://localhost"},
		{"https://localhost", "https://localhost"},
	}
	for _, tt := range tests {
		req.Host = tt.given
		actual := req.URL()
		if actual != tt.expected {
			t.Errorf("given: %s; expected %s, actual %s", tt.given, tt.expected, actual)
		}
	}
}

func TestParseEmptyRequestBody(t *testing.T) {
	selected := "POST /foo/bar"
	result, err := ParseRequest(bytes.NewBufferString(selected))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}
	body, err := result.JSON()
	if err != nil {
		t.Errorf("result.JSON returned %s", err)
	}
	if body != "" {
		t.Errorf("expected '' got %s", body)
	}
}

func TestParsePartial(t *testing.T) {
	selected := `
	POST /foo/bar
	{"size":10,"query":{"term":{"field":"value"}}}
	`
	result, err := ParseRequest(bytes.NewBufferString(selected))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}

	if result.Method != "POST" {
		t.Errorf("Incorrect method: '%s'", result.Method)
	}
	if result.URI != "/foo/bar" {
		t.Errorf("Incorrect URI: '%s'", result.URI)
	}
	if _, err := result.JSON(); err != nil {
		t.Errorf("Must be valid json but got: %s", err)
	}

}

func TestParseErrors(t *testing.T) {
	if _, err := ParseRequest(bytes.NewBufferString("GET")); err == nil {
		t.Errorf("Must return an error")
	}

	result, err := ParseRequest(bytes.NewBufferString(""))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}
	if err := result.Validate(); err == nil {
		t.Errorf("Must be invalid")
	}
}

func TestParseJsonTopLevelArray(t *testing.T) {
	selected := `
	POST /foo/bar
	[
		{"size":10}
	]
	`
	result, err := ParseRequest(bytes.NewBuffer([]byte(selected)))
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}

	if result.Method != "POST" {
		t.Errorf("Incorrect method: '%s'", result.Method)
	}
	if result.URI != "/foo/bar" {
		t.Errorf("Incorrect URI: '%s'", result.URI)
	}
	if _, err := result.JSON(); err != nil {
		t.Errorf("Must be valid json but got: %s %s", err, reflect.TypeOf(err))
	}
}

type blockingReader struct{}

func (r *blockingReader) Read(p []byte) (n int, err error) {
	select {}
}

func TestTryParseAsync(t *testing.T) {
	// should read input
	selected := "POST /foo/bar\n{}"
	var src = bytes.NewBufferString(selected)

	_, err := TryParseAsync(src)
	if err != nil {
		t.Errorf("ParseScript returned %s", err)
	}

	// should handle blocked stream
	if _, err := TryParseAsync(&blockingReader{}); err != errTimeout {
		t.Errorf("expected: %v; actual: %v", errTimeout, err)
	}
}

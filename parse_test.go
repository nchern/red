package main

import (
	"bytes"
	"regexp"
	"testing"
)

const (
	expectedBody = `{"size":10,"query":{"term":{"field":"value"}}}`
)

func TestParse(t *testing.T) {
	src := MustAsset(templateAsset)

	result, err := parseScript(bytes.NewBuffer(src))
	if err != nil {
		t.Errorf("parseScript returned %s", err)
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
	if result.Url() != "http://localhost:9200/foo/_search" {
		t.Errorf("Incorrect Url(): '%s'", result.Url())
	}

	json, err := result.JSON()
	if err != nil {
		t.Errorf("Malformed json: %s", err)
	}
	rx := regexp.MustCompile("\\s")
	json = rx.ReplaceAllString(json, "")
	if json != expectedBody {
		t.Errorf("Unexpected json: [%s]", json)
	}
}

func TestParseEmptyRequestBody(t *testing.T) {
	selected := "POST /foo/bar"
	result, err := parseScript(bytes.NewBuffer([]byte(selected)))
	if err != nil {
		t.Errorf("parseScript returned %s", err)
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
	result, err := parseScript(bytes.NewBuffer([]byte(selected)))
	if err != nil {
		t.Errorf("parseScript returned %s", err)
	}

	if result.Method != "POST" {
		t.Errorf("Incorrect method: '%s'", result.Method)
	}
	if result.URI != "/foo/bar" {
		t.Errorf("Incorrect URI: '%s'", result.URI)
	}

}

func TestParseErrors(t *testing.T) {
	if _, err := parseScript(bytes.NewBuffer([]byte("GET"))); err == nil {
		t.Errorf("Must return an error")
	}

	result, err := parseScript(bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Errorf("parseScript returned %s", err)
	}
	if err := result.Validate(); err == nil {
		t.Errorf("Must be invalid")
	}
}
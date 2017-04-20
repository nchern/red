package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/nchern/red/app"

	color "gopkg.in/fatih/color.v1"
)

const (
	filenameBase = "query"
)

var (
	editor      = env("EDITOR", "vim")
	editorFlags = env("EDITOR_FLAGS", "-O")

	appHomePath   = path.Join(os.Getenv("HOME"), ".red")
	queryFilename = filenameBase + ".txt"
	outFilename   = filenameBase + ".out"

	queryFilePath = path.Join(appHomePath, queryFilename)
	outFilePath   = path.Join(appHomePath, outFilename)
)

func notEmpty(val, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func env(key, defaultVal string) string {
	return notEmpty(os.Getenv(key), defaultVal)
}

func openEditor() error {

	if _, err := os.Stat(queryFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(appHomePath, 0700); err != nil {
			return err
		}
		// path/to/whatever does not exist
		if err := ioutil.WriteFile(queryFilePath, app.MustAsset(app.TemplateAsset), 0644); err != nil {
			return err
		}
	}

	cmdArgs := strings.Split(editorFlags, " ")
	cmdArgs = append(cmdArgs, queryFilePath, outFilePath)

	cmd := exec.Command(editor, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func doRequest(req *app.ParsedRequest) (int, []byte, error) {
	src, err := req.JSON()
	if err != nil {
		return 0, nil, err
	}
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	httpReq, err := http.NewRequest(req.Method, req.Url(), bytes.NewBufferString(src))
	if err != nil {
		return 0, nil, err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, tryFormatJSON(body), nil
}

func tryFormatJSON(body []byte) []byte {
	var out bytes.Buffer
	if err := json.Indent(&out, body, "", app.JsIndent); err != nil {
		return body
	}
	return out.Bytes()
}

func runQuery() error {
	file, err := os.Open(queryFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	request, err := app.ParseScript(file)
	if err != nil {
		return err
	}

	if sel, err := app.TryParseStdinAsync(); err == nil {
		// got the whole query file or it is enough input to use parsed data from stdin
		if sel.Validate() == nil {
			request = sel
		} else {
			request.URI = sel.URI
			request.Method = sel.Method
			request.CopyBodyFrom(sel)
		}
	}

	if err := request.Validate(); err != nil {
		return err
	}
	code, body, err := doRequest(request)
	if err != nil {
		return err
	}
	w, err := os.Create(outFilePath)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "#> %d %s %s\n\n", code, request.Method, request.URI); err != nil {
		return err
	}

	_, err = w.Write(body)
	return err
}

func example() error {
	data := app.MustAsset(app.TemplateAsset)
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}

func doCmd() error {
	if len(os.Args) < 2 {
		return openEditor()
	}

	action := os.Args[1]
	if action == "run" {
		return runQuery()
	}
	if action == "example" {
		return example()
	}
	if action == "fmt" {
		return app.JsonFormat(os.Stdin, os.Stdout, false)
	}

	return fmt.Errorf("Unknown action: %s", action)
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s ", color.RedString("ERROR"))
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func main() {

	if err := doCmd(); err != nil {
		if err == app.ErrFormatFailed {
			os.Exit(1)
		}

		switch err := err.(type) {
		case *exec.ExitError:
		case *app.JsonifyError:
			if syntaxErr, ok := err.Inner.(*json.SyntaxError); ok {
				errorf("Bad JSON query: %s", err)
				fmt.Fprintf(os.Stderr, "%s\n", color.RedString(err.Highlighted(syntaxErr.Offset)))
				break
			}
			errorf("%s", err)
		default:
			errorf("%s", err)
		}
		os.Exit(1)
	}
}

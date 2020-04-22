package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	jsonIndent = "   "

	filenameBase = "query"
)

var (
	editor      = env("EDITOR", "vim")
	editorFlags = env("EDITOR_FLAGS", "-O")

	appHomePath = path.Join(os.Getenv("HOME"), ".red")

	queryFilename = filenameBase + ".red"
	outFilename   = filenameBase + ".redout"

	queryFilePath = path.Join(appHomePath, queryFilename)
	outFilePath   = path.Join(appHomePath, outFilename)

	client = &http.Client{
		Timeout: 3 * time.Second,
	}

	flagCmd = flag.String("c", "edit", "Command to exectue. One of: edit, run, example")

	// opens the editor of preference to edit requests
	cmdEdit = "edit"

	// runs a given query, either from stdin or a query file(TODO: make it accept "-")
	cmdRun = "run"

	// prints out example of request file
	cmdExample = "example"
)

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

func doRequest(req *app.HTTPRequest) (int, []byte, error) {
	src, err := req.JSON()
	if err != nil {
		return 0, nil, err
	}
	httpReq, err := http.NewRequest(req.Method, req.URL(), bytes.NewBufferString(src))
	if err != nil {
		return 0, nil, err
	}
	httpReq.Header = req.Headers
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

func runQuery(primaryReader io.Reader, secondryReader io.Reader, out io.Writer) error {
	request, err := app.ParseRequest(primaryReader)
	if err != nil {
		return err
	}

	if sel, err := app.TryParseAsync(secondryReader); err == nil {
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

	if _, err := fmt.Fprintf(out, "#> %d %s %s\n\n", code, request.Method, request.URL()); err != nil {
		return err
	}

	_, err = out.Write(body)
	return err
}

func example() error {
	data := app.MustAsset(app.TemplateAsset)
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}

func run() error {
	srcReader, err := os.Open(queryFilePath)
	if err != nil {
		return err
	}
	defer srcReader.Close()

	w, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// TODO: get rid of the logic especially if we can accept "-" from cmd line(see corresponding todo)?
	// Mirror stdin to stout - this allows processing selections in vim correctly
	secondary := io.TeeReader(os.Stdin, os.Stdout)

	return runQuery(srcReader, secondary, w)
}

func doCmd() error {
	switch *flagCmd {
	case cmdEdit:
		return openEditor()
	case cmdExample:
		return example()
	case cmdRun:
		return run()
	}

	return fmt.Errorf("Unknown action: %s", *flagCmd)
}

func main() {
	flag.Parse()

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

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s ", color.RedString("ERROR"))
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func notEmpty(val, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func env(key, defaultVal string) string {
	return notEmpty(os.Getenv(key), defaultVal)
}

func tryFormatJSON(body []byte) []byte {
	var out bytes.Buffer
	if err := json.Indent(&out, body, "", jsonIndent); err != nil {
		return body
	}
	return out.Bytes()
}

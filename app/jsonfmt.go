package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	// JsIndent is the identation used to format JSON
	JsIndent = "   "
)

// JSONFormat formats json input in the reader and writes the formatted output to writer
func JSONFormat(reader io.Reader, writer io.Writer, verbose bool) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err := json.Indent(&out, body, "", JsIndent); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			if verbose {
				fmt.Fprintf(os.Stdout, "<<< FORMAT ERROR: %s\n", err)
				os.Stdout.Write(body[:syntaxErr.Offset])
				fmt.Fprintln(os.Stdout, "\n>>>>")
				os.Stdout.Write(body[syntaxErr.Offset:])
				return ErrFormatFailed
			}
			// not verbose: leave input unmodified
			os.Stdout.Write(body)
			return ErrFormatFailed
		}
		return err
	}
	_, err = out.WriteTo(writer)
	return err
}

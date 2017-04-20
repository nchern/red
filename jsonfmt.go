package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func jsonFormat(reader io.Reader, writer io.Writer, verbose bool) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err := json.Indent(&out, body, "", jsIndent); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			if verbose {
				fmt.Fprintf(os.Stdout, "<<< FORMAT ERROR: %s\n", err)
				os.Stdout.Write(body[:syntaxErr.Offset])
				fmt.Fprintln(os.Stdout, "\n>>>>")
				os.Stdout.Write(body[syntaxErr.Offset:])
				return errIndentFailed
			}
			// not verbose: leave input unmodified
			os.Stdout.Write(body)
			return errIndentFailed
		}
		return err
	}
	_, err = out.WriteTo(writer)
	return err
}

// jsonfmt reads json content given on stdin or file, formats it and writes
// result to stdout.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	var rewrite bool
	flag.BoolVar(&rewrite, "w", false, "write result to (source) file instead of stdout")
	flag.Parse()
	src := ""
	if args := flag.Args(); len(args) > 0 {
		src = args[0]
	}
	if err := run(src, rewrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(src string, rewrite bool) error {
	var rd io.ReadCloser = ioutil.NopCloser(os.Stdin)
	if src != "" {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		rd = f
	}
	var data json.RawMessage
	if err := json.NewDecoder(rd).Decode(&data); err != nil {
		return err
	}
	rd.Close()
	if src == "" || !rewrite {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		return enc.Encode(data)
	}
	of, err := os.Create(src)
	if err != nil {
		return err
	}
	defer of.Close()
	enc := json.NewEncoder(of)
	enc.SetIndent("", "\t")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return of.Close()
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: jsonfmt [-w] [path]")
		flag.PrintDefaults()
	}
}

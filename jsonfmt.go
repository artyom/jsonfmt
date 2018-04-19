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
	"path/filepath"
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
	var rd io.Reader = os.Stdin
	if src != "" {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		rd = f
	}
	var enc *json.Encoder
	var wr io.WriteCloser
	var err error
	switch {
	case src == "" || !rewrite:
		wr = os.Stdout
	default:
		var tf *os.File
		if tf, err = ioutil.TempFile(filepath.Dir(src), ".jsonfmt-"); err != nil {
			return err
		}
		defer func() {
			switch err {
			case nil:
				os.Rename(tf.Name(), src)
			default:
				os.Remove(tf.Name())
			}
		}()
		wr = tf
	}
	enc = json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	dec := json.NewDecoder(rd)
	var data json.RawMessage
	for {
		switch err = dec.Decode(&data); err {
		case io.EOF:
			return nil
		case nil:
		default:
			return err
		}
		if err = enc.Encode(data); err != nil {
			return err
		}
	}
	err = wr.Close() // explicit assignment so it's visible in defered closure
	return err
}

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: jsonfmt [-w] [path]")
		flag.PrintDefaults()
	}
}

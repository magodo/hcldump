package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, `Usage:
	hcldump [filename]`)
		os.Exit(1)
	}
	fname := os.Args[1]

	src, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read file: %v", err)
		os.Exit(1)
	}

	file, diags := hclsyntax.ParseConfig(src, fname, hcl.InitialPos)
	if len(diags) != 0 {
		fmt.Fprintln(os.Stderr, "failed to parse: %v", diags.Error())
		os.Exit(1)
	}

	walker := &Walker{0, os.Stdout}

	// The body in the returned file has dynamic type *hclsyntax.Body, so callers may freely
	// type-assert this to get access to the full hclsyntax API in situations where detailed
	// access is required.
	// https://pkg.go.dev/github.com/hashicorp/hcl2/hcl/hclsyntax#ParseConfig
	diags = hclsyntax.Walk(file.Body.(*hclsyntax.Body), walker)
	if len(diags) != 0 {
		fmt.Fprintln(os.Stderr, "failed to walk the ast: %v", diags.Error())
		os.Exit(1)
	}
}

type Walker struct {
	indent int
	writer io.Writer
}

func (w *Walker) Enter(node hclsyntax.Node) hcl.Diagnostics {
	w.indent++
	switch node := node.(type) {
	case *hclsyntax.Attribute:
		// TODO
		w.writer.Write([]byte(strings.Repeat(" ", w.indent) + node.Name + "\n"))
	}
	// TODO
	return nil
}

func (w *Walker) Exit(node hclsyntax.Node) hcl.Diagnostics {
	w.indent--
	return nil
}

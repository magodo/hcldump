package main

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"

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
		fmt.Fprintf(os.Stderr, "cannot read file: %v\n", err)
		os.Exit(1)
	}

	f, diags := hclsyntax.ParseConfig(src, fname, hcl.InitialPos)
	if len(diags) != 0 {
		fmt.Fprintf(os.Stderr, "failed to parse: %v\n", diags.Error())
		os.Exit(1)
	}
	ast.Print(nil, f)
}

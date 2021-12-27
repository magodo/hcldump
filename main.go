package main

import (
	"flag"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

var (
	verbose    bool
	expression bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Show hcl.Range node together")
	flag.BoolVar(&expression, "e", false, "Parse the HCL as an expression")
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  hcldump [filename]

`)
	flag.PrintDefaults()
}

func main() {
	flag.CommandLine.Usage = usage
	flag.Usage = usage

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	fname := args[0]

	src, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read file: %v\n", err)
		os.Exit(1)
	}

	var (
		root  interface{}
		diags hcl.Diagnostics
	)
	if expression {
		root, diags = hclsyntax.ParseExpression(src, fname, hcl.InitialPos)
	} else {
		root, diags = hclsyntax.ParseConfig(src, fname, hcl.InitialPos)
	}
	if len(diags) != 0 {
		fmt.Fprintf(os.Stderr, "failed to parse: %v\n", diags.Error())
		//os.Exit(1)
	}

	var filter ast.FieldFilter
	if !verbose {
		filter = func(name string, value reflect.Value) bool {
			if _, ok := value.Interface().(hcl.Range); ok {
				return false
			}
			if _, ok := value.Interface().([]hcl.Range); ok {
				return false
			}
			return true
		}
	}
	if err := ast.Fprint(os.Stdout, nil, root, filter); err != nil {
		fmt.Fprintf(os.Stderr, "print error: %v", err)
		os.Exit(1)
	}
}

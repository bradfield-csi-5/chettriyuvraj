package main

import (
	"bytes"
	"log"
	"os"
	"sort"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const src string = `package foo

import (
	"fmt"
	"time"
)

func baz() {
	fmt.Println("Hello, world!")
}

type A int

const b = "testing"

func bar() {
	fmt.Println(time.Now())
}

func aar() {
	fmt.Println(time.Now())
}
`

// Moves all top-level functions to the end, sorted in alphabetical order.
// The "source file" is given as a string (rather than e.g. a filename).
func SortFunctions(src string) (string, error) {
	// TODO
	f, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}

	/* Separate funcs and non funcs in two lists */
	rootList := f.Decls
	funcList := make([]*dst.FuncDecl, 0)
	nonFuncList := make([]dst.Decl, 0)
	for _, decl := range rootList {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if ok == true {
			funcList = append(funcList, funcDecl)
		} else {
			nonFuncList = append(nonFuncList, decl)
		}
	}

	/* Sort funcList */
	sort.Slice(funcList, func(i, j int) bool { return funcList[i].Name.Name < funcList[j].Name.Name })

	/* Combine the two lists in original list */
	rootIdx := 0
	for _, decl := range nonFuncList {
		rootList[rootIdx] = decl
		rootIdx += 1
	}
	for _, decl := range funcList {
		rootList[rootIdx] = decl
		rootIdx += 1
	}

	/* Print new source */
	var b bytes.Buffer
	err = decorator.Fprint(&b, f)
	if err != nil {
		log.Fatal(err)
	}

	return b.String(), nil
}

func main() {
	f, err := decorator.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	// Print AST
	err = dst.Fprint(os.Stdout, f, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Convert AST back to source
	err = decorator.Print(f)
	if err != nil {
		log.Fatal(err)
	}

	SortFunctions(src)
}

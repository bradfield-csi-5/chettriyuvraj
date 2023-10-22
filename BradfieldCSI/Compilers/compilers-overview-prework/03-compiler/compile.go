package main

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
)

// Given an AST node corresponding to a function (guaranteed to be
// of the form `func f(x, y byte) byte`), compile it into assembly
// code.
//
// Recall from the README that the input parameters `x` and `y` should
// be read from memory addresses `1` and `2`, and the return value
// should be written to memory address `0`.
func compile(node *ast.FuncDecl) (string, error) {
	var b strings.Builder
	assemblyList := make([]strings.Builder, 0)

	/* Iterate over Body.List of FuncDecl */
	for _, n := range node.Body.List {
		assembly, err := evalNode(n)
		if err != nil {
			return "", err
		}
		assemblyList = append(assemblyList, assembly...)
	}

	/* Compile all string builders into one */
	for _, sb := range assemblyList {
		fmt.Fprintf(&b, sb.String())
	}

	return b.String(), nil
}

func evalNode(node ast.Stmt) ([]strings.Builder, error) {
	var res []strings.Builder
	var err error

	switch n := node.(type) {
	case *ast.ReturnStmt:
		res, err = evalReturn(n)
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

/* General Format of ReturnStmt to assembly:
pushi <evaluated value>
pop 0
halt
*/

func evalReturn(node *ast.ReturnStmt) ([]strings.Builder, error) {
	var b strings.Builder

	/* Since its the return statement, we are expecting a single elem in node.Results */
	val, err := evalExpr(node.Results[0])
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(&b, "pushi %d\npop 0\nhalt", val)
	// fmt.Printf("\n\n string %s\n\n", b.String())
	return []strings.Builder{b}, nil
}

func evalExpr(node ast.Expr) (byte, error) {
	var res byte
	var err error

	// fmt.Printf("\n\n %s", node)

	switch n := node.(type) {
	case *ast.BasicLit:
		// fmt.Printf("\n\n node val %s", n.Value)
		tempRes, tempErr := strconv.Atoi(n.Value)
		err = tempErr
		if err != nil {
			break
		}
		res = byte(tempRes)
	}

	if err != nil {
		return 0, err
	}
	return res, nil
}

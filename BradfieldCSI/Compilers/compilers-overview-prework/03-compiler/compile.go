package main

import (
	"fmt"
	"go/ast"
	"go/token"
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
	/* Since its the return statement, we are expecting a single elem in node.Results */
	var b strings.Builder
	sb, err := evalExpr(node.Results[0])
	if err != nil {
		return nil, err
	}

	b.WriteString("pop 0\nhalt\n")
	sb = append(sb, b)
	return sb, nil
}

func evalExpr(node ast.Expr) ([]strings.Builder, error) {
	var res []strings.Builder
	var err error

	switch n := node.(type) {
	case *ast.BasicLit:
		var b strings.Builder

		tempRes, tempErr := strconv.Atoi(n.Value)
		err = tempErr
		if err != nil {
			break
		}

		fmt.Fprintf(&b, "pushi %d\n", byte(tempRes))
		res = append(res, b)

	case *ast.BinaryExpr:
		leftRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}
		res = append(res, leftRes...)

		rightRes, tempErr := evalExpr(n.Y)
		err = tempErr
		if err != nil {
			break
		}
		res = append(res, rightRes...)
		res = append(res, getOpAsStringBuilder(n.Op))

	case *ast.ParenExpr:
		tempRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}

		res = tempRes
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

func evalOp(left byte, right byte, op token.Token) byte {
	var res byte
	switch x := op; x {
	case token.ADD:
		res = left + right
	case token.SUB:
		res = left - right
	case token.MUL:
		res = left * right
	case token.QUO:
		res = left / right
	}
	return res
}

func getOpAsStringBuilder(op token.Token) strings.Builder {
	var res strings.Builder
	switch x := op; x {
	case token.ADD:
		res.WriteString("add\n")
	case token.SUB:
		res.WriteString("sub\n")
	case token.MUL:
		res.WriteString("mul\n")
	case token.QUO:
		res.WriteString("div\n")
	}
	return res
}

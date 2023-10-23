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

/* Map variable name to its location in memory + next location in line to be mapped */
var varToMem map[string]int = make(map[string]int)
var nextMemoryLocation int = 3

func compile(node *ast.FuncDecl) (string, error) {
	var b strings.Builder

	/* Assumptions */
	varToMem["x"] = 1
	varToMem["y"] = 2

	/* node.Body always expected to be BlockStmt */
	sbList, err := evalNode(node.Body)
	if err != nil {
		return "", err
	}

	/* Compile all string builders into one */
	for _, sb := range sbList {
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
	case *ast.AssignStmt:
		res, err = evalAssign(n)
	case *ast.IfStmt:
		res, err = evalIfElse(n)
	case *ast.ForStmt:
		res, err = evalFor(n)
	case *ast.BlockStmt:
		var blockRes []strings.Builder
		for _, blockNode := range n.List {
			curBlockRes, err := evalNode(blockNode)
			if err != nil {
				break
			}
			blockRes = append(blockRes, curBlockRes...)
		}
		res = blockRes
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

/* General Format of AssignStmt to assembly:
{assembly for evaluated expression which puts it on top of the stack}
pop <next memory location>
*/

func evalAssign(node *ast.AssignStmt) ([]strings.Builder, error) {
	/* Since its the assign stmt, we are expecting an identifier on the LHS */
	var b strings.Builder
	sb, err := evalExpr(node.Rhs[0])
	if err != nil {
		return nil, err
	}

	/* Expecting Ident on LHS */
	ident, ok := node.Lhs[0].(*ast.Ident)
	if ok != true {
		return nil, fmt.Errorf("invalid RHS of assignment")
	}

	/* Take assignment var, (record in map for later reference depending - might be new/existing), also generate assembly to pop it to memory location*/
	varName := ident.Name
	existingMemoryLocation, exists := varToMem[varName]
	if exists != true {
		varToMem[varName] = nextMemoryLocation
		fmt.Fprintf(&b, "pop %d\n", nextMemoryLocation)
		nextMemoryLocation += 1
		sb = append(sb, b)
		return sb, nil
	}

	fmt.Fprintf(&b, "pop %d\n", existingMemoryLocation)
	sb = append(sb, b)
	return sb, nil
}

/* Currently handles only single if/else block */
func evalIfElse(node *ast.IfStmt) ([]strings.Builder, error) {
	sb := []strings.Builder{}
	var l1, l2 strings.Builder
	l1.WriteString("label l1\n")
	l2.WriteString("label l2\n")

	/* Determine which if label to be executed depending on result of if body */
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return nil, err
	}
	condStr := strings.Builder{}
	condStr.WriteString("jeqz l2\n")
	sbcond = append(sbcond, condStr)

	/* We will assign all conditions a label */
	sbif, err := evalNode(node.Body)
	if err != nil {
		return nil, err
	}

	sbelse, err := evalNode(node.Else)
	if err != nil {
		return nil, err
	}

	sb = append(sb, sbcond...)
	sb = append(sb, l1)
	sb = append(sb, sbif...)
	sb = append(sb, l2)
	sb = append(sb, sbelse...)

	return sb, nil
}

/* General Format of ForStmt to assembly:
{condition evaluate and jump}
label f1:
{assembly for inner expr}
{condition evaluate and jump}
pop <next memory location>
*/

/*
pushi 1
ne      -> 0 != 1 -> 1 So in case where

ne      -> 1 != 1 -> 0
*/

func evalFor(node *ast.ForStmt) ([]strings.Builder, error) {
	/* Jump to loop test */
	var s strings.Builder
	fmt.Fprint(&s, "jump l2\n")

	/* Define labels for loop content (l1) and loop test (l2) */
	sb := []strings.Builder{s}
	var l1, l2 strings.Builder
	l1.WriteString("label l1\n")
	l2.WriteString("label l2\n")

	/* Repeat loop content if condition satisfied, inverting result of condition so we can use jeqz to loop*/
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return nil, err
	}
	condStr := strings.Builder{}
	condStr.WriteString("pushi 1\nneq\njeqz l1\n")
	sbcond = append(sbcond, condStr)

	sbbody, err := evalNode(node.Body)
	if err != nil {
		return nil, err
	}

	sb = append(sb, l1)
	sb = append(sb, sbbody...)
	sb = append(sb, l2)
	sb = append(sb, sbcond...)

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

	case *ast.Ident:
		var b strings.Builder

		memoryLocation, exists := varToMem[n.Name]
		if exists != true {
			err = fmt.Errorf("variable does not exist")
			break
		}

		fmt.Fprintf(&b, "push %d\n", memoryLocation)
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
	case token.EQL:
		res.WriteString("eq\n")
	case token.LSS:
		res.WriteString("lt\n")
	case token.GTR:
		res.WriteString("gt\n")
	}
	return res
}

/*

ne

*/

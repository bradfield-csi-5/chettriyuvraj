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

	/* Assumptions */
	varToMem["x"] = 1
	varToMem["y"] = 2

	/* node.Body always expected to be BlockStmt */
	sb, err := evalNode(node.Body)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func evalNode(node ast.Stmt) (*strings.Builder, error) {
	res := &strings.Builder{}
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
		for _, blockNode := range n.List {
			blockRes, err := evalNode(blockNode)
			if err != nil {
				break
			}
			res.WriteString(blockRes.String())
		}
	}

	if err != nil {
		return &strings.Builder{}, err
	}
	return res, nil
}

/* General Format of ReturnStmt to assembly:
pushi <evaluated value>
pop 0
halt
*/

func evalReturn(node *ast.ReturnStmt) (*strings.Builder, error) {
	/* Since its the return statement, we are expecting a single elem in node.Results */
	sb, err := evalExpr(node.Results[0])
	if err != nil {
		return &strings.Builder{}, err
	}

	sb.WriteString("pop 0\nhalt\n")
	return sb, nil
}

/* General Format of AssignStmt to assembly:
{assembly for evaluated expression which puts it on top of the stack}
pop <next memory location>
*/

func evalAssign(node *ast.AssignStmt) (*strings.Builder, error) {
	/* Since its the assign stmt, we are expecting an identifier on the LHS */
	sb, err := evalExpr(node.Rhs[0])
	if err != nil {
		return &strings.Builder{}, err
	}

	/* Expecting Ident on LHS */
	ident, ok := node.Lhs[0].(*ast.Ident)
	if ok != true {
		return &strings.Builder{}, fmt.Errorf("invalid RHS of assignment")
	}

	/* Take assignment var, (record in map for later reference depending - might be new/existing), also generate assembly to pop it to memory location*/
	varName := ident.Name
	existingMemoryLocation, exists := varToMem[varName]
	if exists != true {
		varToMem[varName] = nextMemoryLocation
		fmt.Fprintf(sb, "pop %d\n", nextMemoryLocation)
		nextMemoryLocation += 1
		return sb, nil
	}

	fmt.Fprintf(sb, "pop %d\n", existingMemoryLocation)
	return sb, nil
}

/* Currently handles only single if/else block */
func evalIfElse(node *ast.IfStmt) (*strings.Builder, error) {

	/* Determine which if label to be executed depending on result of if body */
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return &strings.Builder{}, err
	}
	sbcond.WriteString("jeqz l2\n")

	/* We will assign all conditions a label */
	sbif, err := evalNode(node.Body)
	if err != nil {
		return &strings.Builder{}, err
	}

	sbelse, err := evalNode(node.Else)
	if err != nil {
		return &strings.Builder{}, err
	}

	/* Append in sequence to create assembly for if else statement */
	sbcond.WriteString("label l1\n")
	sbcond.WriteString(sbif.String())
	sbcond.WriteString("label l2\n")
	sbcond.WriteString(sbelse.String())

	return sbcond, nil
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

func evalFor(node *ast.ForStmt) (*strings.Builder, error) {
	/* Jump to loop test */
	var sb strings.Builder

	/* Label l2 - condition: Repeat loop content if condition satisfied, inverting result of condition so we can use jeqz to loop*/
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return &strings.Builder{}, err
	}
	sbcond.WriteString("pushi 1\nneq\njeqz l1\n")

	/* Label l1 - loop body*/
	sbbody, err := evalNode(node.Body)
	if err != nil {
		return &strings.Builder{}, err
	}

	/* Append in sequence to create assembly for if else statement */
	sb.WriteString("jump l2\n")
	sb.WriteString("label l1\n")
	sb.WriteString(sbbody.String())
	sb.WriteString("label l2\n")
	sb.WriteString(sbcond.String())

	return &sb, nil
}

func evalExpr(node ast.Expr) (*strings.Builder, error) {
	var res strings.Builder
	var err error

	switch n := node.(type) {
	case *ast.BasicLit:

		tempRes, tempErr := strconv.Atoi(n.Value)
		err = tempErr
		if err != nil {
			break
		}

		fmt.Fprintf(&res, "pushi %d\n", byte(tempRes))

	case *ast.Ident:
		memoryLocation, exists := varToMem[n.Name]
		if exists != true {
			err = fmt.Errorf("variable does not exist")
			break
		}

		fmt.Fprintf(&res, "push %d\n", memoryLocation)

	case *ast.BinaryExpr:
		leftRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}
		res.WriteString(leftRes.String())

		rightRes, tempErr := evalExpr(n.Y)
		err = tempErr
		if err != nil {
			break
		}
		opsb := getOpAsStringBuilder(n.Op)
		res.WriteString(rightRes.String())
		res.WriteString(opsb.String())

	case *ast.ParenExpr:
		tempRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}

		res.WriteString(tempRes.String())
	}

	if err != nil {
		return &strings.Builder{}, err
	}
	return &res, nil
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

func getOpAsStringBuilder(op token.Token) *strings.Builder {
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
	return &res
}

/*

ne

*/

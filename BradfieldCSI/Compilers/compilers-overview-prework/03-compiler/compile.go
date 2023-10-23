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

	/* Variables always expected in these locations */
	varToMem["x"] = 1
	varToMem["y"] = 2

	/* node.Body always expected to be BlockStmt */
	res, err := evalStmt(node.Body)
	if err != nil {
		return "", err
	}

	return res, nil
}

/****** Functions for evaluating interface types ******/

/**
 * Evaluate nodes implementing the Stmt interface
 *
 * @param ast.Stmt
 * @return string
 * @return error
 *
 * - Concrete type check
 * - Evaluate the concrete type to get its assembly
 * - Append result to string builder
 * - Return string builder as string
 **/

func evalStmt(node ast.Stmt) (string, error) {
	var sb strings.Builder
	var res string
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
			blockRes, err := evalStmt(blockNode)
			if err != nil {
				break
			}
			sb.WriteString(blockRes)
		}
	}

	if err != nil {
		return "", err
	}
	sb.WriteString(res)

	return sb.String(), nil
}

/**
 * Evaluate nodes implementing Expr Interface
 *
 * @param ast.Expr
 * @return string
 * @return error
 *
 * - Type check for concrete type of node
 * - Generate assembly for concrete type
 * - Return assembly as string
 *
 **/

func evalExpr(node ast.Expr) (string, error) {
	var sb strings.Builder
	var err error

	switch n := node.(type) {
	case *ast.BasicLit:

		tempRes, tempErr := strconv.Atoi(n.Value)
		err = tempErr
		if err != nil {
			break
		}

		fmt.Fprintf(&sb, "pushi %d\n", byte(tempRes))

	case *ast.Ident:
		memoryLocation, exists := varToMem[n.Name]
		if exists != true {
			err = fmt.Errorf("variable does not exist")
			break
		}

		fmt.Fprintf(&sb, "push %d\n", memoryLocation)

	case *ast.BinaryExpr:
		leftRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}
		sb.WriteString(leftRes)

		rightRes, tempErr := evalExpr(n.Y)
		err = tempErr
		if err != nil {
			break
		}
		opsb := getOpAsStringBuilder(n.Op)
		sb.WriteString(rightRes)
		sb.WriteString(opsb)

	case *ast.ParenExpr:
		tempRes, tempErr := evalExpr(n.X)
		err = tempErr
		if err != nil {
			break
		}

		sb.WriteString(tempRes)
	}

	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

/****** Functions for evaluating concrete types ******/

/**
 * Evaluate nodes of concrete type ReturnStmt
 *
 * @param *ast.ReturnStmt
 * @return string
 * @return error
 *
 * General Format of ReturnStmt to assembly:
 * pushi <evaluated value of expreesion>
 * pop 0
 * halt
 *
 * - Extract node of interface type ast.Expr from node.Result list
 * - Grab the assembly as result of evaluating the expression
 * - Generate assembly in correct format
 * - Return generated assembly as string
 **/

func evalReturn(node *ast.ReturnStmt) (string, error) {
	var sb strings.Builder

	/* Since its the return statement, we are expecting a single elem in node.Results */
	res, err := evalExpr(node.Results[0])
	if err != nil {
		return "", err
	}

	sb.WriteString(res)
	sb.WriteString("pop 0\nhalt\n")
	return sb.String(), nil
}

/**
 * Evaluate nodes of concrete type ast.AssignStmt
 *
 * @param *ast.AssignStmt
 * @return string
 * @return error
 *
 * General Format of AssignStmt to assembly:
 * {assembly for evaluated expression which puts result on top of the stack}
 * pop <correct memory location>
 *
 * - Grab result of evaluating expression on RHS of assignment
 * - Grab identifier on LHS of expression
 * - Generate assembly in correct format
 * - Return assembly as string
 **/

func evalAssign(node *ast.AssignStmt) (string, error) {
	var sb strings.Builder

	/* Since its the assign stmt, we are expecting an expression on the RHS */
	exprRes, err := evalExpr(node.Rhs[0])
	if err != nil {
		return "", err
	}

	/* Expecting Ident on LHS */
	ident, ok := node.Lhs[0].(*ast.Ident)
	if ok != true {
		return "", fmt.Errorf("invalid RHS of assignment")
	}

	/* Generate assembly to pop it to memory location, depending on new or existing variable*/
	sb.WriteString(exprRes)

	varName := ident.Name
	existingMemoryLocation, exists := varToMem[varName]
	if exists != true {
		varToMem[varName] = nextMemoryLocation
		fmt.Fprintf(&sb, "pop %d\n", nextMemoryLocation)
		nextMemoryLocation += 1
	} else {
		fmt.Fprintf(&sb, "pop %d\n", existingMemoryLocation)
	}

	return sb.String(), nil
}

/**
 * Evaluate nodes of the concrete type ast.IfStmt
 *
 * @param *ast.IfStmt
 * @return string
 * @return error
 *
 * General Format of IfStmt to assembly:
 * {assembly for condition evaluation}
 * jeqz l2
 * label l1
 * {assembly for if block}
 * label l2
 * {assemblyt for else block}
 *
 * - Generate assembly from evaluating condition
 * - Generate assembly for if block
 * - Generate assembly for else block
 * - Generate assembly for if-else by combining these along with jeqz and labels
 * - Return assembly as string
 *
 * NOTE:
 * 1. Label names hardcoded so can't evaluate multiple nested blocks or if-else inside for loop
 **/

func evalIfElse(node *ast.IfStmt) (string, error) {
	var sb strings.Builder

	/* Generate ssembly for condition statement */
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return "", err
	}

	/* We will assign all conditions a label */
	sbif, err := evalStmt(node.Body)
	if err != nil {
		return "", err
	}

	sbelse, err := evalStmt(node.Else)
	if err != nil {
		return "", err
	}

	/* Append in sequence to create assembly for if else statement */
	sb.WriteString(sbcond)
	sb.WriteString("jeqz l2\n")
	sb.WriteString("label l1\n")
	sb.WriteString(sbif)
	sb.WriteString("label l2\n")
	sb.WriteString(sbelse)

	return sb.String(), nil
}

/**
 * Evaluate nodes of the concrete type ForStmt
 *
 * @param *ast.ForStmt
 * @return string
 * @return error
 *
 * General Format of ForStmt to assembly:
 * jump l2
 * label l1:
 * {assembly for loop block}
 * label l2:
 * {assembly for loop condition}
 * pushi 1
 * neq
 * jeqz l1
 * pop <next memory location>
 *
 * - Generate assembly for evaluating loop condition
 * - Generate assembly for loop block
 * - Generate assembly for entire loop by combining these along with jeqz, neq and labels
 * - Return assembly as string
 *
 * NOTE:
 * 1. We invert the result of the loop condition using neq so that jeqz can be used to jump to loop block
 * 2. Label names hardcoded so can't evaluate multiple nested blocks or if-else inside for loop
 **/

func evalFor(node *ast.ForStmt) (string, error) {
	/* Jump to loop test */
	var sb strings.Builder

	/* Label l2 - condition: Repeat loop content if condition satisfied, inverting result of condition so we can use jeqz to loop*/
	sbcond, err := evalExpr(node.Cond)
	if err != nil {
		return "", err
	}

	/* Label l1 - loop body*/
	sbbody, err := evalStmt(node.Body)
	if err != nil {
		return "", err
	}

	/* Append in sequence to create assembly for if else statement */
	sb.WriteString("jump l2\n")
	sb.WriteString("label l1\n")
	sb.WriteString(sbbody)
	sb.WriteString("label l2\n")
	sb.WriteString(sbcond)
	sb.WriteString("pushi 1\nneq\njeqz l1\n")

	return sb.String(), nil
}

/****** Helpers ******/

func getOpAsStringBuilder(op token.Token) string {
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
	return res.String()
}

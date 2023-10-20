package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
)

// Given an expression containing only int types, evaluate
// the expression and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	/* Base case */
	basicLit, ok := expr.(*ast.BasicLit)
	if ok == true {
		return strconv.Atoi(basicLit.Value)
	}

	/* Recursive case */
	parenExpr, ok := expr.(*ast.ParenExpr)
	if ok == true {
		val, err := Evaluate((parenExpr.X))
		if err != nil {
			return -1, err
		}
		return val, nil
	}

	binaryExpr, ok := expr.(*ast.BinaryExpr)
	if ok == true {
		left, err := Evaluate((binaryExpr.X))
		if err != nil {
			return -1, err
		}

		right, err := Evaluate((binaryExpr.Y))
		if err != nil {
			return -1, err
		}

		res := 0
		switch x := binaryExpr.Op; x {
		case token.ADD:
			res = left + right
		case token.SUB:
			res = left - right
		case token.MUL:
			res = left * right
		case token.QUO:
			res = left / right
		}
		return res, nil
	}

	return -1, fmt.Errorf("Invalid AST assumptions")
}

/* Test cases
1 + 2 + 3
1 + (3 * 4)
(1+2)+(3+4)
*/

func main() {
	// expr, err := parser.ParseExpr("1 + 2 - 3 * 4")
	// expr, err := parser.ParseExpr("(1 + 2) - (3 * 4)")
	expr, err := parser.ParseExpr("1 + 2 * 3")
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	err = ast.Print(fset, expr)
	if err != nil {
		log.Fatal(err)
	}

	res, err := Evaluate(expr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

/*
		BEXP
  	X:BEXP	 Y:BEXP
X:LIT   Y:LIT

Inorder traversal:
# base case
	#basicLit
		#print(val);; return val


# rec case
	#binaryExp
	#leftVal = recurseLeft ie X
	#print op
	#rightVal = recurse right ie Y
	#return left op right

	#paranExp
	#print paranopen
	#val = recurseX
	#print paranClose
	#return val

1 + 2 - 3 + 5



*/

package coleval

import (
	ast "colon/colast"
	obj "colon/colobj"
	"fmt"
	"os"
)

var (
	booleanTrue  = &obj.Boolean{Value: true}
	booleanFalse = &obj.Boolean{Value: false}
	// EMPTY : the Null / Nil equivalent object in colon
	EMPTY = &obj.Empty{}
)

// PostEvalOutput : contains the results of running the program
// var PostEvalOutput []obj.Object

// Eval : evaluates the ast obtained after parsing
func Eval(node ast.Node) obj.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.PrefixExpression:
		rightExpression := Eval(node.RightExpression)
		return evalPrefixExpression(node.Operator, rightExpression)

	case *ast.InfixExpression:
		leftExpression := Eval(node.LeftExpression)
		rightExpression := Eval(node.RightExpression)
		return evalInfixExpression(node.Operator, leftExpression, rightExpression)

	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}

	case *ast.FloatingLiteral:
		return &obj.Floating{Value: node.Value}

	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}

	case *ast.BooleanLiteral:
		if node.Value == true {
			return booleanTrue
		}
		return booleanFalse

	case *ast.Block:
		return evalBlock(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		retVal := Eval(node.ReturnValue)
		return &obj.ReturnValue{Value: retVal}
	}
	return nil
}

func evalProgram(program *ast.Program) obj.Object {
	var res obj.Object
	for _, statement := range program.Statements {
		res = Eval(statement)
		// PostEvalOutput = append(PostEvalOutput, res)
		if retVal, ok := res.(*obj.ReturnValue); ok {
			return retVal.Value
		}
	}
	return res
}

func evalBlock(block *ast.Block) obj.Object {
	var res obj.Object
	for _, statement := range block.Statements {
		res = Eval(statement)
		// PostEvalOutput = append(PostEvalOutput, res)
		if res != nil && res.ObType() == obj.RETVAL {
			return res
		}
	}
	return res
}

func evalPrefixExpression(operator string, rightExpression obj.Object) obj.Object {
	switch operator {
	case "!":
		return evalLogicalNotOf(rightExpression)
	case "-":
		return evalNumericNegation(rightExpression)
	default:
		reportRuntimeError(fmt.Sprintf("unknown prefix operation.\nOperator used => [ %v ]", operator))
	}
	return nil
}

func evalLogicalNotOf(rightExpression obj.Object) obj.Object {
	switch rightExpression.ObType() {

	case obj.BOOLEAN:
		return makeBooleanObject(!(rightExpression.(*obj.Boolean).Value))

	default:
		reportRuntimeError(fmt.Sprintf("cannot perform LOGICAL_NOT operation on type %q\n", rightExpression.ObType()))
	}
	return nil
}

func makeBooleanObject(b bool) obj.Object {
	if b {
		return booleanTrue
	}
	return booleanFalse
}

func evalNumericNegation(rightExpression obj.Object) obj.Object {
	switch rightExpression.ObType() {
	case obj.INTEGER:
		return &obj.Integer{Value: -(rightExpression.(*obj.Integer).Value)}
	case obj.FLOATING:
		return &obj.Floating{Value: -(rightExpression.(*obj.Floating).Value)}
	default:
		reportRuntimeError(fmt.Sprintf("Runtime Error: cannot perform NUMERIC_NEGATION operation on type %q\n", rightExpression.ObType()))
	}
	return nil
}

func evalInfixExpression(operator string, leftExpression obj.Object, rightExpression obj.Object) obj.Object {
	leftExprType := leftExpression.ObType()
	rightExprType := rightExpression.ObType()

	if leftExprType == obj.INTEGER && rightExprType == obj.INTEGER {
		return evalIntIntInfix(operator, leftExpression, rightExpression)
	} else if leftExprType == obj.FLOATING && rightExprType == obj.FLOATING {
		return evalFltFltInfix(operator, leftExpression, rightExpression)
	} else if leftExprType == obj.INTEGER && rightExprType == obj.FLOATING {
		return evalIntFltInfix(operator, leftExpression, rightExpression)
	} else if leftExprType == obj.FLOATING && rightExprType == obj.INTEGER {
		return evalFltIntInfix(operator, leftExpression, rightExpression)
	} else if leftExprType == obj.STRING && rightExprType == obj.STRING {
		return evalStrStrInfix(operator, leftExpression, rightExpression)
	} else if leftExprType == obj.BOOLEAN && rightExprType == obj.BOOLEAN {
		return evalBolBolInfix(operator, leftExpression, rightExpression)
	} else {
		reportRuntimeError(fmt.Sprintf("illegal infix expression encountered. Operator that operates on %q and %q at not found", leftExprType, rightExprType))
	}

	return nil
}

func evalIntIntInfix(op string, l obj.Object, r obj.Object) obj.Object {
	lVal := l.(*obj.Integer).Value
	rVal := r.(*obj.Integer).Value
	switch op {
	case "+":
		return &obj.Integer{
			Value: lVal + rVal,
		}
	case "-":
		return &obj.Integer{
			Value: lVal - rVal,
		}
	case "*":
		return &obj.Integer{
			Value: lVal * rVal,
		}
	case "/":
		return &obj.Integer{
			Value: lVal / rVal,
		}
	case "%":
		return &obj.Integer{
			Value: lVal % rVal,
		}
	case "^":
		var result int64 = 1
		var i int64
		for i = 0; i < rVal; i++ {
			result = result * lVal
		}
		return &obj.Integer{
			Value: result,
		}
	case ">":
		return makeBooleanObject(lVal > rVal)
	case "<":
		return makeBooleanObject(lVal < rVal)
	case ">=":
		return makeBooleanObject(lVal >= rVal)
	case "<=":
		return makeBooleanObject(lVal <= rVal)
	case "==":
		return makeBooleanObject(lVal == rVal)
	case "!=":
		return makeBooleanObject(lVal != rVal)
	default:
		reportRuntimeError(fmt.Sprintf("unknown operation performed.\nOperator used => [ %v ]", op))
	}
	return nil
}

func evalFltFltInfix(op string, l obj.Object, r obj.Object) obj.Object {
	lVal := l.(*obj.Floating).Value
	rVal := r.(*obj.Floating).Value
	switch op {
	case "+":
		return &obj.Floating{
			Value: lVal + rVal,
		}
	case "-":
		return &obj.Floating{
			Value: lVal - rVal,
		}
	case "*":
		return &obj.Floating{
			Value: lVal * rVal,
		}
	case "/":
		return &obj.Floating{
			Value: lVal / rVal,
		}
	case ">":
		return makeBooleanObject(lVal > rVal)
	case "<":
		return makeBooleanObject(lVal < rVal)
	case ">=":
		return makeBooleanObject(lVal >= rVal)
	case "<=":
		return makeBooleanObject(lVal <= rVal)
	case "==":
		return makeBooleanObject(lVal == rVal)
	case "!=":
		return makeBooleanObject(lVal != rVal)
	default:
		reportRuntimeError(fmt.Sprintf("unknown operation performed.\nOperator used => [ %v ]", op))
	}
	return nil
}

func evalIntFltInfix(op string, l obj.Object, r obj.Object) obj.Object {
	lVal := float64(l.(*obj.Integer).Value)
	rVal := r.(*obj.Floating).Value
	switch op {
	case "+":
		return &obj.Floating{
			Value: lVal + rVal,
		}
	case "-":
		return &obj.Floating{
			Value: lVal - rVal,
		}
	case "*":
		return &obj.Floating{
			Value: lVal * rVal,
		}
	case "/":
		return &obj.Floating{
			Value: lVal / rVal,
		}
	case ">":
		return makeBooleanObject(lVal > rVal)
	case "<":
		return makeBooleanObject(lVal < rVal)
	case ">=":
		return makeBooleanObject(lVal >= rVal)
	case "<=":
		return makeBooleanObject(lVal <= rVal)
	case "==":
		return makeBooleanObject(lVal == rVal)
	case "!=":
		return makeBooleanObject(lVal != rVal)
	default:
		reportRuntimeError(fmt.Sprintf("unknown operation performed.\nOperator used => [ %v ]", op))
	}
	return nil
}

func evalFltIntInfix(op string, l obj.Object, r obj.Object) obj.Object {
	lVal := l.(*obj.Floating).Value
	rVal := float64(r.(*obj.Integer).Value)
	switch op {
	case "+":
		return &obj.Floating{
			Value: lVal + rVal,
		}
	case "-":
		return &obj.Floating{
			Value: lVal - rVal,
		}
	case "*":
		return &obj.Floating{
			Value: lVal * rVal,
		}
	case "/":
		return &obj.Floating{
			Value: lVal / rVal,
		}
	case ">":
		return makeBooleanObject(lVal > rVal)
	case "<":
		return makeBooleanObject(lVal < rVal)
	case ">=":
		return makeBooleanObject(lVal >= rVal)
	case "<=":
		return makeBooleanObject(lVal <= rVal)
	case "==":
		return makeBooleanObject(lVal == rVal)
	case "!=":
		return makeBooleanObject(lVal != rVal)
	default:
		reportRuntimeError(fmt.Sprintf("unknown operation performed.\nOperator used => [ %v ]", op))
	}
	return nil
}

func evalStrStrInfix(op string, l obj.Object, r obj.Object) obj.Object {
	if op != "+" {
		reportRuntimeError(fmt.Sprintf("cannot use this operator on string operands.\nOperator used => [ %v ]", op))
	}
	return &obj.String{
		Value: l.(*obj.String).Value + r.(*obj.String).Value,
	}
}

func contains(arr []obj.ObjectType, typ obj.ObjectType) bool {
	for _, v := range arr {
		if v == typ {
			return true
		}
	}
	return false
}

func getIntValueFromObj(object obj.Object) int64 {
	return object.(*obj.Integer).Value
}

func getFltValueFromObj(object obj.Object) float64 {
	return object.(*obj.Floating).Value
}

func getBolValueFromObj(object obj.Object) bool {
	return object.(*obj.Boolean).Value
}

func getStrValueFromObj(object obj.Object) string {
	return object.(*obj.String).Value
}

func evalIfExpression(ife *ast.IfExpression) obj.Object {
	condition := Eval(ife.Condition)
	if condition.ObType() != obj.BOOLEAN {
		reportRuntimeError(fmt.Sprintf("Condition does not evaluate to `true` or `false`"))
	}
	if getBolValueFromObj(condition) {
		return Eval(ife.IfBody)
	} else if ife.ElseBody != nil {
		return Eval(ife.ElseBody)
	} else {
		return EMPTY
	}
}

func reportRuntimeError(msg string) {
	fmt.Printf("Runtime Error: %s", msg)
	os.Exit(22)
}

func evalBolBolInfix(op string, l obj.Object, r obj.Object) obj.Object {
	switch op {
	case "==":
		return makeBooleanObject(l == r)
	case "!=":
		return makeBooleanObject(l != r)
	case "&":
		return makeBooleanObject(getBolValueFromObj(l) && getBolValueFromObj(r))
	case "|":
		return makeBooleanObject(getBolValueFromObj(l) || getBolValueFromObj(r))
	default:
		reportRuntimeError(fmt.Sprintf("operator %q cannot operator on BOOLEAN values", op))
	}
	return nil
}

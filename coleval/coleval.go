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
)

// PostEvalOutput : contains the results of running the program
var PostEvalOutput []obj.Object

// Eval : evaluates the ast obtained after parsing
func Eval(node ast.Node) obj.Object {
	switch node := node.(type) {

	case *ast.Program:
		var res obj.Object
		for _, statement := range node.Statements {
			res = Eval(statement)
			PostEvalOutput = append(PostEvalOutput, res)
		}
		return res

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
	}

	return nil
}

func evalPrefixExpression(operator string, rightExpression obj.Object) obj.Object {
	switch operator {
	case "!":
		return evalLogicalNotOf(rightExpression)
	case "-":
		return evalNumericNegation(rightExpression)
	default:
		fmt.Printf("Runtime Error : unknown prefix operation.\nOperator used => [ %v ]", operator)
		os.Exit(22)
	}
	return nil
}

func evalLogicalNotOf(rightExpression obj.Object) obj.Object {
	switch rightExpression.ObType() {

	case obj.BOOLEAN:
		return makeBooleanObject(!(rightExpression.(*obj.Boolean).Value))

	default:
		fmt.Printf("Runtime Error: cannot perform LOGICAL_NOT operation on type %q\n", rightExpression.ObType())
		os.Exit(22)
	}
	return nil
}

func makeBooleanObject(b bool) obj.Object {
	switch b {
	case true:
		return booleanTrue
	case false:
		return booleanFalse
	}
	return nil
}

func evalNumericNegation(rightExpression obj.Object) obj.Object {
	switch rightExpression.ObType() {
	case obj.INTEGER:
		return &obj.Integer{Value: -(rightExpression.(*obj.Integer).Value)}
	case obj.FLOATING:
		return &obj.Floating{Value: -(rightExpression.(*obj.Floating).Value)}
	default:
		fmt.Printf("Runtime Error: cannot perform NUMERIC_NEGATION operation on type %q\n", rightExpression.ObType())
		os.Exit(22)
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
	}

	return nil
}

func evalIntIntInfix(op string, l obj.Object, r obj.Object) obj.Object {
	switch op {
	case "+":
		return &obj.Integer{
			Value: l.(*obj.Integer).Value + r.(*obj.Integer).Value,
		}
	case "-":
		return &obj.Integer{
			Value: l.(*obj.Integer).Value - r.(*obj.Integer).Value,
		}
	case "*":
		return &obj.Integer{
			Value: l.(*obj.Integer).Value * r.(*obj.Integer).Value,
		}
	case "/":
		return &obj.Integer{
			Value: l.(*obj.Integer).Value / r.(*obj.Integer).Value,
		}
	case "%":
		return &obj.Integer{
			Value: l.(*obj.Integer).Value % r.(*obj.Integer).Value,
		}
	case "^":
		var result int64 = 1
		var i int64
		for i = 0; i < r.(*obj.Integer).Value; i++ {
			result = result * l.(*obj.Integer).Value
		}
		return &obj.Integer{
			Value: result,
		}
	default:
		fmt.Printf("Runtime Error : unknown operation performed.\nOperator used => [ %v ]", op)
	}
	return nil
}

func evalFltFltInfix(op string, l obj.Object, r obj.Object) obj.Object {
	switch op {
	case "+":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value + r.(*obj.Floating).Value,
		}
	case "-":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value - r.(*obj.Floating).Value,
		}
	case "*":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value * r.(*obj.Floating).Value,
		}
	case "/":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value / r.(*obj.Floating).Value,
		}
	default:
		fmt.Printf("Runtime Error : unknown operation performed.\nOperator used => [ %v ]", op)
	}
	return nil
}

func evalIntFltInfix(op string, l obj.Object, r obj.Object) obj.Object {
	switch op {
	case "+":
		return &obj.Floating{
			Value: float64(l.(*obj.Integer).Value) + r.(*obj.Floating).Value,
		}
	case "-":
		return &obj.Floating{
			Value: float64(l.(*obj.Integer).Value) - r.(*obj.Floating).Value,
		}
	case "*":
		return &obj.Floating{
			Value: float64(l.(*obj.Integer).Value) * r.(*obj.Floating).Value,
		}
	case "/":
		return &obj.Floating{
			Value: float64(l.(*obj.Integer).Value) / r.(*obj.Floating).Value,
		}
	default:
		fmt.Printf("Runtime Error : unknown operation performed.\nOperator used => [ %v ]", op)
	}
	return nil
}

func evalFltIntInfix(op string, l obj.Object, r obj.Object) obj.Object {
	switch op {
	case "+":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value + float64(r.(*obj.Integer).Value),
		}
	case "-":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value - float64(r.(*obj.Integer).Value),
		}
	case "*":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value * float64(r.(*obj.Integer).Value),
		}
	case "/":
		return &obj.Floating{
			Value: l.(*obj.Floating).Value / float64(r.(*obj.Integer).Value),
		}
	default:
		fmt.Printf("Runtime Error : unknown operation performed.\nOperator used => [ %v ]", op)
	}
	return nil
}

func evalStrStrInfix(op string, l obj.Object, r obj.Object) obj.Object {
	if op != "+" {
		fmt.Printf("Runtime Error : Cannot use this operator on string operands.\nOperator used => [ %v ]", op)
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

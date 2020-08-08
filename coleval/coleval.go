package coleval

import (
	ast "colon/colast"
	obj "colon/colobj"
	"fmt"
	"os"
)

// Storing values that are reused frequently
var (
	booleanTrue  = &obj.Boolean{Value: true}
	booleanFalse = &obj.Boolean{Value: false}
	// EMPTY : the Null / Nil equivalent object in colon
	EMPTY = &obj.Empty{}
)

// to know whether the evaluation happening at any moment is
// inside a loop or not
var inLoop bool = false

// Eval : evaluates the ast obtained after parsing
func Eval(node ast.Node, env *obj.Env) obj.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		rightExpression := Eval(node.RightExpression, env)
		return evalPrefixExpression(node.Operator, rightExpression, env)

	case *ast.InfixExpression:
		leftExpression := Eval(node.LeftExpression, env)
		rightExpression := Eval(node.RightExpression, env)
		return evalInfixExpression(node.Operator, leftExpression, rightExpression, env)

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
		return evalBlock(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		retVal := Eval(node.ReturnValue, env)
		return &obj.ReturnValue{Value: retVal}

	case *ast.VarStatement:
		varVal := Eval(node.Value, env)
		if varVal == EMPTY {
			reportRuntimeError(fmt.Sprintf("expression assigned to variable %q did not evaluate to a value of a legal datatype", node.Name.Value))
		}
		if inLoop && env.IsInside() {
			if _, inOuter := env.ContainedIn.Get(node.Name.Value); inOuter {
				env.ContainedIn.Set(node.Name.Value, varVal)
			} else {
				env.Set(node.Name.Value, varVal)
			}
		} else {
			env.Set(node.Name.Value, varVal)
		}

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionExpression:
		params := node.Params
		body := node.FuncBody
		return &obj.Function{
			Parameters: params,
			FuncBody:   body,
			Env:        env,
		}

	case *ast.FunctionCallExpression:
		function := Eval(node.Function, env)
		if function == EMPTY {
			reportRuntimeError(fmt.Sprintf("function %q not found.", node.Function.String()))
		}
		arguments := evalExpressions(node.Arguments, env)
		// i'm hoping that evalExpressions catches all the runtime errors
		return evalFunction(arguments, function)

	case *ast.LoopExpression:
		// condition := node.Condition
		// body := node.LoopBody
		// loopObj := &obj.Loop{
		// 	Condition: condition,
		// 	LoopBody: body,
		// 	Env: env,
		// }
		return evalLoopExpression(node, env)

	}
	return nil
}

func evalProgram(program *ast.Program, env *obj.Env) obj.Object {
	var res obj.Object
	for _, statement := range program.Statements {
		res = Eval(statement, env)
		// PostEvalOutput = append(PostEvalOutput, res)
		if retVal, ok := res.(*obj.ReturnValue); ok {
			return retVal.Value
		}
	}
	return res
}

func evalBlock(block *ast.Block, env *obj.Env) obj.Object {
	var res obj.Object
	for _, statement := range block.Statements {
		res = Eval(statement, env)
		// PostEvalOutput = append(PostEvalOutput, res)
		if res != nil && res.ObType() == obj.RETVAL {
			return res
		}
	}
	if res == nil {
		return EMPTY
	}
	return res
}

func evalPrefixExpression(operator string, rightExpression obj.Object, env *obj.Env) obj.Object {
	switch operator {
	case "!":
		return evalLogicalNotOf(rightExpression, env)
	case "-":
		return evalNumericNegation(rightExpression, env)
	default:
		reportRuntimeError(fmt.Sprintf("unknown prefix operation.\nOperator used => [ %v ]", operator))
	}
	return nil
}

func evalLogicalNotOf(rightExpression obj.Object, env *obj.Env) obj.Object {
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

func evalNumericNegation(rightExpression obj.Object, env *obj.Env) obj.Object {
	switch rightExpression.ObType() {
	case obj.INTEGER:
		return &obj.Integer{Value: -(rightExpression.(*obj.Integer).Value)}
	case obj.FLOATING:
		return &obj.Floating{Value: -(rightExpression.(*obj.Floating).Value)}
	default:
		reportRuntimeError(fmt.Sprintf("cannot perform NUMERIC_NEGATION operation on type %q\n", rightExpression.ObType()))
	}
	return nil
}

func evalInfixExpression(operator string, leftExpression obj.Object, rightExpression obj.Object, env *obj.Env) obj.Object {
	leftExprType := leftExpression.ObType()
	rightExprType := rightExpression.ObType()

	if leftExprType == obj.INTEGER && rightExprType == obj.INTEGER {
		return evalIntIntInfix(operator, leftExpression, rightExpression, env)
	} else if leftExprType == obj.FLOATING && rightExprType == obj.FLOATING {
		return evalFltFltInfix(operator, leftExpression, rightExpression, env)
	} else if leftExprType == obj.INTEGER && rightExprType == obj.FLOATING {
		return evalIntFltInfix(operator, leftExpression, rightExpression, env)
	} else if leftExprType == obj.FLOATING && rightExprType == obj.INTEGER {
		return evalFltIntInfix(operator, leftExpression, rightExpression, env)
	} else if leftExprType == obj.STRING && rightExprType == obj.STRING {
		return evalStrStrInfix(operator, leftExpression, rightExpression, env)
	} else if leftExprType == obj.BOOLEAN && rightExprType == obj.BOOLEAN {
		return evalBolBolInfix(operator, leftExpression, rightExpression, env)
	} else {
		reportRuntimeError(fmt.Sprintf("illegal infix expression encountered. Operator that operates on %q and %q at not found", leftExprType, rightExprType))
	}

	return nil
}

func evalIntIntInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
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

func evalFltFltInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
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

func evalIntFltInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
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

func evalFltIntInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
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

func evalStrStrInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
	if op == "+" {
		return &obj.String{
			Value: l.(*obj.String).Value + r.(*obj.String).Value,
		}
	} else if op == "==" {
		return &obj.Boolean{
			Value: l.(*obj.String).Value == r.(*obj.String).Value,
		}
	} else if op == "!=" {
		return &obj.Boolean{
			Value: l.(*obj.String).Value == r.(*obj.String).Value,
		}
	} else {
		reportRuntimeError(fmt.Sprintf("cannot use this operator on string operands.\nOperator used => [ %v ]", op))
	}
	return nil
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

func evalIfExpression(ife *ast.IfExpression, env *obj.Env) obj.Object {
	condition := Eval(ife.Condition, env)
	if condition.ObType() != obj.BOOLEAN {
		reportRuntimeError(fmt.Sprintf("Condition does not evaluate to `true` or `false`"))
	}
	if getBolValueFromObj(condition) {
		return Eval(ife.IfBody, env)
	} else if ife.ElseBody != nil {
		return Eval(ife.ElseBody, env)
	} else {
		return EMPTY
	}
}

func reportRuntimeError(msg string) {
	fmt.Printf("Runtime Error: %s\n", msg)
	os.Exit(22)
}

func evalBolBolInfix(op string, l obj.Object, r obj.Object, env *obj.Env) obj.Object {
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

func evalIdentifier(identifier *ast.Identifier, env *obj.Env) obj.Object {
	if boundVal, ok := env.Get(identifier.Value); ok {
		return boundVal
	}
	if bin, ok := builtin[identifier.Value]; ok {
		return bin
	}
	reportRuntimeError(fmt.Sprintf("variable %q not initialized. Cannot use uninitialized variables in expressions", identifier.Value))
	return nil
}

func evalExpressions(args []ast.Expression, env *obj.Env) []obj.Object {
	evaluatedEArgs := []obj.Object{}
	for _, v := range args {
		ev := Eval(v, env)
		// Not sure if there are any potential errors that can occur here
		if ev != nil {
			evaluatedEArgs = append(evaluatedEArgs, ev)
		}
	}
	return evaluatedEArgs
}

func evalFunction(arguments []obj.Object, function obj.Object) obj.Object {
	switch funct := function.(type) {
	case *obj.Function:
		functEnv := createNewSubEnv(arguments, funct)
		evaluatedFunct := Eval(funct.FuncBody, functEnv)
		return unwrapRetVal(evaluatedFunct)
	case *obj.BuiltIn:
		return funct.Bfunct(arguments...)
	default:
		reportRuntimeError(fmt.Sprintf("expression %q is not/ doesn't have a valid funcion definition", function.ObValue()))
	}
	return nil
}

// creates a "scope" for the function being called.
// this allows for creation of "local variables", i.e.,
// variables that are local to the function
func createNewSubEnv(arguments []obj.Object, function *obj.Function) *obj.Env {
	functEnv := obj.NewInnerEnv(function.Env)
	for k, v := range function.Parameters {
		functEnv.Set(v.Value, arguments[k])
	}
	return functEnv
}

func unwrapRetVal(EvalResult obj.Object) obj.Object {
	if retval, ok := EvalResult.(*obj.ReturnValue); ok {
		return retval.Value
	}
	return EvalResult
}

// @experimental
func evalLoopExpression(le *ast.LoopExpression, env *obj.Env) obj.Object {
	var loopResult obj.Object
	loopEnv := obj.NewInnerEnv(env)
	condition := Eval(le.Condition, loopEnv)
	if condition.ObType() != obj.BOOLEAN {
		reportRuntimeError(fmt.Sprintf("Condition does not evaluate to `true` or `false`"))
	}

	// to signify that the evaluation is going inside a loop
	inLoop = true

	for getBolValueFromObj(condition) {
		loopResult = Eval(le.LoopBody, loopEnv)
		condition = Eval(le.Condition, loopEnv)
		if condition.ObType() != obj.BOOLEAN {
			reportRuntimeError(fmt.Sprintf("Condition does not evaluate to `true` or `false`"))
		}
	}

	// to signify that the evaluation has come out of a loop
	inLoop = false

	return unwrapRetVal(loopResult)
}

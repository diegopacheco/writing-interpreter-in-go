package evaluator

import (
	"fmt"
	"reflect"

	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/ast"
	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/object"
	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/tracing"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	if tracing.IsDebugMode() {
		fmt.Printf(">>> Eval START: Received node type = %s\n", reflect.TypeOf(node))
	}
	var result object.Object

	switch node := node.(type) {
	case *ast.Program:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.Program")
		}
		result = evalProgram(node, env)

	case *ast.ExpressionStatement:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.ExpressionStatement")
		}
		result = Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.IntegerLiteral")
		}
		result = &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.Boolean")
		}
		result = nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.StringLiteral")
		}
		result = &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.PrefixExpression")
		}
		right := Eval(node.Right, env)
		if isError(right) {
			result = right
		} else {
			result = evalPrefixExpression(node.Operator, right)
		}

	case *ast.InfixExpression:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.InfixExpression")
		}
		left := Eval(node.Left, env)
		if isError(left) {
			result = left
		} else {
			right := Eval(node.Right, env)
			if isError(right) {
				result = right
			} else {
				result = evalInfixExpression(node.Operator, left, right) // Assign to result
			}
		}

	case *ast.BlockStatement:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.BlockStatement")
		}
		result = evalBlockStatement(node, env)

	case *ast.IfExpression:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.IfExpression")
		}
		result = evalIfExpression(node, env)

	case *ast.ReturnStatement:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.ReturnStatement")
		}
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			result = val
		} else {
			result = &object.ReturnValue{Value: val}
		}

	case *ast.LetStatement:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.LetStatement")
		}
		val := Eval(node.Value, env)
		if isError(val) {
			result = val
		} else {
			env.Set(node.Name.Value, val)
			result = val
		}

	case *ast.Identifier:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.Identifier")
		}
		result = evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.FunctionLiteral")
		}
		result = &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

	case *ast.CallExpression:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.CallExpression")
		}
		fn := Eval(node.Function, env)
		if isError(fn) {
			result = fn
		} else {
			args := evalExpressions(node.Arguments, env)
			if len(args) > 0 && isError(args[0]) {
				result = args[0]
			} else {
				result = applyFunction(fn, args)
			}
		}

	case *ast.ArrayLiteral:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.ArrayLiteral")
		}
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		result = &object.Array{Elements: elements}

	case *ast.IndexExpression:
		if tracing.IsDebugMode() {
			fmt.Println("    Eval: Handling *ast.IndexExpression")
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		result = evalIndexExpression(left, index)

	default:
		if tracing.IsDebugMode() {
			fmt.Printf("    Eval: Handling default case for type %T\n", node)
		}
		result = nil
	}

	if tracing.IsDebugMode() {
		fmt.Printf("<<< Eval END: Returning result type = %s, value = %+v\n", reflect.TypeOf(result), result)
	}
	return result
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch res := result.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		if operator == "+" {
			return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		}
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier,
	env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	if builtins := builtins[node.Value]; builtins != nil {
		return builtins
	}
	return newError("identifier not found: %s", node.Value)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if fn == nil {
		return newError("not a function: nil")
	}

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		builtin := fn.Fn(args...)
		if isError(builtin) {
			return builtin
		}
		return builtin
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

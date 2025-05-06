package evaluator

import (
	"fmt"
	"testing"

	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/lexer"
	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/object"
	"github.com/diegopacheco/writing-interpreter-in-go/objectsystem/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10) * 2", 30},
		{"(5 + 10) * (2 + 3)", 75},
		{"(5 + 10) * (2 + 3) * 2", 150},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		panic(fmt.Sprintf("parser had errors for input '%s': %v", input, p.Errors()))
	}
	if program == nil {
		panic(fmt.Sprintf("ParseProgram() returned nil for input '%s'", input))
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T,
	obj object.Object,
	expected int64) bool {

	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T,
	obj object.Object,
	expected bool) bool {

	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfelseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } }", 10},
		{`if (10 > 1) { if (10 > 1) { return 10; } return 1; }`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + (true + false);", "type mismatch: BOOLEAN + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "type mismatch: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "type mismatch: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{"if (10 > 1) { true + false; }", "type mismatch: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", "type mismatch: BOOLEAN + BOOLEAN"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = b + 5; c;", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	env := object.NewEnvironment()
	evaluated := Eval(program, env)

	if evaluated == nil {
		t.Fatalf("Eval returned nil directly")
	}

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClosure(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let newAdder = fn(x) { fn(y) { x + y; }; }; let addTwo = newAdder(2); addTwo(2);", 4},
		{"let newAdder = fn(x) { fn(y) { x + y; }; }; let addTwo = newAdder(2); addTwo(10);", 12},
		{"let newAdder = fn(x) { fn(y) { x + y; }; }; let addTwo = newAdder(2); let addThree = newAdder(3); addThree(addTwo(10));", 15},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Fatalf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Fatalf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"len(\"\")", 0},
		{"len(\"Hello World!\")", 12},
		{"len(\"Hello\" + \" \" + \"World!\")", 12},
		{"len(1)", "argument to `len` not supported, got INTEGER"},
		{"len(1, 2)", "wrong number of arguments. got=2"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Fatalf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		default:
			t.Fatalf("unknown expected type %T", expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(arr.Elements))
	}
	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let arr = [1, 2, 3]; len(arr)", 3},
		{"len(push([1, 2, 3], 4))", 4},
		{"first([1, 2, 3])", 1},
		{"last([1, 2, 3])", 3},
		{"len(rest([1, 2, 3]))", 2},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Fatalf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		default:
			t.Fatalf("unknown expected type %T", expected)
		}
	}
}

func TestHashLiteralsAdvanced(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		"two": 1 + 1,
		"three": 6 / 2,
		4:4,
		true: 5,
		false: 6,
	}
	`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("object is not Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		found := false
		for resultKey, pair := range result.Pairs {
			if expectedKey == resultKey {
				found = true
				testIntegerObject(t, pair.Value, expectedValue)
				break
			}
		}
		if !found {
			t.Errorf("no pair for given key. got=%v", expectedKey)
		}
	}

}

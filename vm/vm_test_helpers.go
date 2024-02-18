package vm

import (
	"cidoka/ast"
	"cidoka/compiler"
	"cidoka/lexer"
	"cidoka/object"
	"cidoka/parser"
	"fmt"
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null. got=%T (%+v)", actual, actual)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("object is not String. got=%T (%+v)", actual, actual)
		}
	case []int:
		err := testIntArray(expected, actual)
		if err != nil {
			t.Errorf("testIntArray failed: %s", err)
		}
	case map[object.HashKey]int64:
		err := testIntHash(expected, actual)
		if err != nil {
			t.Errorf("testIntHash failed: %s", err)
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)
			return
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message. expected=%q, got=%q",
				expected.Message, errObj.Message)
		}
	case float64:
		err := testFloatObject(expected, actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case float32:
		err := testFloatObject(float64(expected), actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
	}

	return nil
}

func testIntArray(expected []int, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)", actual, actual)
	}

	if len(result.Elements) != len(expected) {
		return fmt.Errorf("wrong number of elements. got=%d, want=%d",
			len(result.Elements), len(expected))
	}

	for i, integer := range expected {
		err := testIntegerObject(int64(integer), result.Elements[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func testIntHash(expected map[object.HashKey]int64, actual object.Object) error {
	result, ok := actual.(*object.Hash)
	if !ok {
		return fmt.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
	}

	if len(result.Pairs) != len(expected) {
		return fmt.Errorf("wrong number of pairs. got=%d, want=%d", len(result.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			return fmt.Errorf("no pair for given key in pairs")
		}

		err := testIntegerObject(expectedValue, pair.Value)
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %s", err)
		}
	}

	return nil
}

func testFloatObject(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}

	comparisonValue := math.Abs(expected-result.Value) <= float64EqualityThreshold
	if !comparisonValue {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}

	return nil
}

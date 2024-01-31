package vm

import (
	"boludolang/ast"
	"boludolang/compiler"
	"boludolang/lexer"
	"boludolang/object"
	"boludolang/parser"
	"fmt"
	"testing"
)

type VmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []VmTestCase) {
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
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
	}

	return nil
}

func testIntArray(expected []int, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)",
			actual, actual)
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

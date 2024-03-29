package compiler

import (
	"cidoka/code"
	"testing"
)

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", last.Opcode, code.OpSub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbolTable")
	}

	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.OpAdd)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d", last.Opcode, code.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != code.OpMul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d", previous.Opcode, code.OpMul)
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLessThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 <= 1",
			expectedConstants: []interface{}{1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLessOrEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 >= 1",
			expectedConstants: []interface{}{1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterOrEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true && true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpTrue),
				code.Make(code.OpAnd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "true || true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpTrue),
				code.Make(code.OpOr),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) { 10 }; 3333;
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNull),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			if (true) { 10 } else { 20 }; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 0010
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one = 1;
			let two = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDeclareGlobal, 1),
			},
		},
		{
			input: `
			let one = 1;
			one;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpDeclareGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []interface{}{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"mon" + "key"`,
			expectedConstants: []interface{}{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1, 2, 3]",
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpGetIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []interface{}{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpGetIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { 1; 2 }`,
			expectedConstants: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `fn() { 24 }();`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0), // The compiled function
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let noArg = fn() { 24 };
			noArg();
			`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0), // The literal "24"
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0), // The compiled function
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let oneArg = fn(a) { };
			oneArg(24);
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
				24,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let manyArg = fn(a, b, c) { };
			manyArg(24, 25, 26);
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let oneArg = fn(a) { a };
			oneArg(24);
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				24,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let manyArg = fn(a, b, c) { a; b; c };
			manyArg(24, 25, 26);
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 2),
					code.Make(code.OpReturnValue),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestLetStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let num = 55;
			fn() { num }
			`,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			fn() {
				let num = 55;
				num
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			fn() {
				let a = 55;
				let b = 77;
				a + b
			}
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDeclareLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			len([]);
			push([], 1);
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpGetBuiltin, 0),
				code.Make(code.OpArray, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpGetBuiltin, 5),
				code.Make(code.OpArray, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `fn() { len([]) }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetBuiltin, 0),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			fn(a) {
				fn(b) {
					a + b
				}
			}
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			fn(a) {
				fn(b) {
					fn(c) {
						a + b + c
					}
				}
			};
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 1, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let global = 55;

			fn() {
				let a = 66;

				fn() {
					let b = 77;

					fn() {
						let c = 88;

						global + a + b + c;
					}
				}
			}
			`,
			expectedConstants: []interface{}{
				55,
				66,
				77,
				88,
				[]code.Instructions{
					code.Make(code.OpConstant, 3),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 2),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 4, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 5, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpClosure, 6, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let countDown = fn(x) { countDown(x - 1); };
			countDown(1);
			`,
			expectedConstants: []interface{}{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSub),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) { countDown(x - 1); };
				countDown(1);
			};
			wrapper();
			`,
			expectedConstants: []interface{}{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSub),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
				[]code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 3, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFloatingPointArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1.2 + 3.4",
			expectedConstants: []interface{}{1.2, 3.4},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.2 - 3.4",
			expectedConstants: []interface{}{1.2, 3.4},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.2 * 3.4",
			expectedConstants: []interface{}{1.2, 3.4},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1.2 / 3.4",
			expectedConstants: []interface{}{1.2, 3.4},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestAssignment(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let a = 1;
			a = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			let b = a;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpDeclareGlobal, 1),
			},
		},
		{
			input: `
			let a = 1;
			let b = a;
			a = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpDeclareGlobal, 1),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestFunctionsAlreadyDeclaredArg(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let a = 1;
			fn() { let a = 2; };
			`,
			expectedConstants: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = fn() { let a = 2; a; };
			a();
			`,
			expectedConstants: []interface{}{
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpDeclareLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = fn(a) { a; };
			a(1);
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestForLoops(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				for (let i = 0; i < 10; i = i + 1) { i };
				`,
			expectedConstants: []interface{}{
				0,
				10,
				1,
				[]code.Instructions{
					// 0000 - init (let i = 0)
					code.Make(code.OpConstant, 0),
					// 0003
					code.Make(code.OpDeclareLocal, 0),
					// 0005 - condition (i < 10)
					code.Make(code.OpGetLocal, 0),
					// 0007
					code.Make(code.OpConstant, 1),
					// 0010
					code.Make(code.OpLessThan),
					// 0011 - exit loop if false
					code.Make(code.OpJumpNotTruthy, 29),
					// 0014 - loop body ({ i })
					code.Make(code.OpGetLocal, 0),
					// 0016
					code.Make(code.OpPop),
					// 0017 - loop increment (i = i + 1)
					code.Make(code.OpGetLocal, 0),
					// 0019
					code.Make(code.OpConstant, 2),
					// 0022
					code.Make(code.OpAdd),
					// 0023
					code.Make(code.OpSetLocal, 0),
					// 0025
					code.Make(code.OpPop),
					// 0026 - jump back to condition
					code.Make(code.OpJump, 5),
					// 0029 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpLoop, 3),
			},
		},
		{
			input: `
				let sum = 0;
				for (let i = 0; i < 10; i = i + 1) {
					sum = sum + i
				};
				sum;
				`,
			expectedConstants: []interface{}{
				0,
				0,
				10,
				1,
				[]code.Instructions{
					// 0000 - init (let i = 0)
					code.Make(code.OpConstant, 1),
					// 0003
					code.Make(code.OpDeclareLocal, 0),
					// 0005 - condition (i < 10)
					code.Make(code.OpGetLocal, 0),
					// 0007
					code.Make(code.OpConstant, 2),
					// 0010
					code.Make(code.OpLessThan),
					// 0011 - exit loop if false
					code.Make(code.OpJumpNotTruthy, 36),
					// 0014 - loop body ({ sum = sum + i })
					code.Make(code.OpGetGlobal, 0),
					// 0017
					code.Make(code.OpGetLocal, 0),
					// 0019
					code.Make(code.OpAdd),
					// 0020
					code.Make(code.OpSetGlobal, 0),
					// 0023
					code.Make(code.OpPop),
					// 0024 - loop increment (i = i + 1)
					code.Make(code.OpGetLocal, 0),
					// 0026
					code.Make(code.OpConstant, 3),
					// 0029
					code.Make(code.OpAdd),
					// 0030
					code.Make(code.OpSetLocal, 0),
					// 0032
					code.Make(code.OpPop),
					// 0033 - jump back to condition
					code.Make(code.OpJump, 5),
					// 0036 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpLoop, 4),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let sum = 0;
			for (let i = 0; i < 10; i = i + 1) {
				if (i == 5) {
					break;
				}
				sum = sum + i;
			};
			sum;
			`,
			expectedConstants: []interface{}{
				0,
				0,
				10,
				5,
				1,
				[]code.Instructions{
					// 0000 - init (let i = 0)
					code.Make(code.OpConstant, 1),
					// 0003
					code.Make(code.OpDeclareLocal, 0),
					// 0005 - condition (i < 10)
					code.Make(code.OpGetLocal, 0),
					// 0007
					code.Make(code.OpConstant, 2),
					// 0010
					code.Make(code.OpLessThan),
					// 0011 - exit loop if false
					code.Make(code.OpJumpNotTruthy, 51),
					// 0014 - loop body ({ if ... })
					code.Make(code.OpGetLocal, 0),
					// 0016
					code.Make(code.OpConstant, 3),
					// 0019
					code.Make(code.OpEqual),
					// 0020
					code.Make(code.OpJumpNotTruthy, 27),
					// 0023
					code.Make(code.OpBreak),
					// 0024
					code.Make(code.OpJump, 28),
					// 0027
					code.Make(code.OpNull),
					// 0028
					code.Make(code.OpPop),
					// 0029
					code.Make(code.OpGetGlobal, 0),
					// 0032
					code.Make(code.OpGetLocal, 0),
					// 0034
					code.Make(code.OpAdd),
					// 0035
					code.Make(code.OpSetGlobal, 0),
					// 0038
					code.Make(code.OpPop),
					// 0039 - loop increment (i = i + 1)
					code.Make(code.OpGetLocal, 0),
					// 0041
					code.Make(code.OpConstant, 4),
					// 0044
					code.Make(code.OpAdd),
					// 0045
					code.Make(code.OpSetLocal, 0),
					// 0047
					code.Make(code.OpPop),
					// 0048 - jump back to condition
					code.Make(code.OpJump, 5),
					// 0051 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpDeclareGlobal, 0),
				// 0006
				code.Make(code.OpLoop, 5),
				// 0009
				code.Make(code.OpGetGlobal, 0),
				// 0012
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let i = 0;
			for (;;) {
				if (i == 5) {
					break;
				}
				i = i + 1;
			}
			i;
			`,
			expectedConstants: []interface{}{
				0,
				5,
				1,
				[]code.Instructions{
					// 0000 - init
					code.Make(code.OpNull),
					// 0001
					code.Make(code.OpPop),
					// 0002 - condition
					code.Make(code.OpTrue),
					// 0003
					code.Make(code.OpJumpNotTruthy, 38),
					// 0006 - loop body
					code.Make(code.OpGetGlobal, 0),
					// 0009
					code.Make(code.OpConstant, 1),
					// 0012
					code.Make(code.OpEqual),
					// 0013
					code.Make(code.OpJumpNotTruthy, 20),
					// 0016
					code.Make(code.OpBreak),
					// 0017
					code.Make(code.OpJump, 21),
					// 0020
					code.Make(code.OpNull),
					// 0021
					code.Make(code.OpPop),
					// 0022
					code.Make(code.OpGetGlobal, 0),
					// 0025
					code.Make(code.OpConstant, 2),
					// 0028
					code.Make(code.OpAdd),
					// 0029
					code.Make(code.OpSetGlobal, 0),
					// 0032
					code.Make(code.OpPop),
					// 0033 - increment
					code.Make(code.OpNull),
					// 0034
					code.Make(code.OpPop),
					// 0035 - jump back to condition
					code.Make(code.OpJump, 2),
					// 0038 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpConstant, 0),
				// 0003
				code.Make(code.OpDeclareGlobal, 0),
				// 0006
				code.Make(code.OpLoop, 3),
				// 0009
				code.Make(code.OpGetGlobal, 0),
				// 0012
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestCompoundAssignment(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let a = 1;
			a += 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a -= 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a *= 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a /= 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a %= 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMod),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a += 2 * 3;
			`,
			expectedConstants: []interface{}{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpMul),
				code.Make(code.OpAdd),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestContinueStatement(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			while (true) {
				continue;
			}
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					// 0000 - init
					code.Make(code.OpNull),
					// 0001
					code.Make(code.OpPop),
					// 0002 - condition
					code.Make(code.OpTrue),
					// 0003
					code.Make(code.OpJumpNotTruthy, 14),
					// 0006 - loop body
					code.Make(code.OpJump, 9),
					// 0009 - update
					code.Make(code.OpNull),
					// 0010
					code.Make(code.OpPop),
					// 0011 - jump back to condition
					code.Make(code.OpJump, 2),
					// 0014 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpLoop, 0),
			},
		},
		{
			input: `
			let sum = 0;
			let i = 0;
			while (i < 10) {
				for (let j = 0; j < 10; j += 1) {
					if (j == 5) {
						continue;
					}
					sum += j;
				}
				i += 1;
			}
			sum;
			`,
			expectedConstants: []interface{}{
				0,
				0,
				10,
				0,
				10,
				5,
				1,
				[]code.Instructions{
					// 0000 - init (let j = 0)
					code.Make(code.OpConstant, 3),
					// 0003
					code.Make(code.OpDeclareLocal, 0),
					// 0005 - condition (j < 10)
					code.Make(code.OpGetLocal, 0),
					// 0007
					code.Make(code.OpConstant, 4),
					// 0010
					code.Make(code.OpLessThan),
					// 0011 - exit loop if false
					code.Make(code.OpJumpNotTruthy, 53),
					// 0014 - loop body ({ if ... })
					code.Make(code.OpGetLocal, 0),
					// 0016
					code.Make(code.OpConstant, 5),
					// 0019
					code.Make(code.OpEqual),
					// 0020
					code.Make(code.OpJumpNotTruthy, 29),
					// 0023
					code.Make(code.OpJump, 41),
					// 0026
					code.Make(code.OpJump, 30),
					// 0029
					code.Make(code.OpNull),
					// 0030
					code.Make(code.OpPop),
					// 0031
					code.Make(code.OpGetGlobal, 0),
					// 0034
					code.Make(code.OpGetLocal, 0),
					// 0036
					code.Make(code.OpAdd),
					// 0037
					code.Make(code.OpSetGlobal, 0),
					// 0040
					code.Make(code.OpPop),
					// 0041 - loop increment (j = j + 1)
					code.Make(code.OpGetLocal, 0),
					// 0043
					code.Make(code.OpConstant, 6),
					// 0046
					code.Make(code.OpAdd),
					// 0047
					code.Make(code.OpSetLocal, 0),
					// 0049
					code.Make(code.OpPop),
					// 0050 - jump back to condition
					code.Make(code.OpJump, 5),
					// 0053 - exit loop
					code.Make(code.OpBreak),
				},
				1,
				[]code.Instructions{
					// 0000 - init (nil)
					code.Make(code.OpNull),
					// 0001
					code.Make(code.OpPop),
					// 0002 - condition (i < 10)
					code.Make(code.OpGetGlobal, 1),
					// 0005
					code.Make(code.OpConstant, 2),
					// 0008
					code.Make(code.OpLessThan),
					// 0009 - exit loop if false
					code.Make(code.OpJumpNotTruthy, 31),
					// 0012 - loop body
					code.Make(code.OpLoop, 7),
					// 0015
					code.Make(code.OpGetGlobal, 1),
					// 0018
					code.Make(code.OpConstant, 8),
					// 0021
					code.Make(code.OpAdd),
					// 0022
					code.Make(code.OpSetGlobal, 1),
					// 0025
					code.Make(code.OpPop),
					// 0026 - loop increment (nil)
					code.Make(code.OpNull),
					// 0027
					code.Make(code.OpPop),
					// 0028 - jump back to condition
					code.Make(code.OpJump, 2),
					// 0031 - exit loop
					code.Make(code.OpBreak),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDeclareGlobal, 1),
				code.Make(code.OpLoop, 9),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestPostfixOperators(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let a = 1;
			a++;
			`,
			expectedConstants: []interface{}{1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let a = 1;
			a--;
			`,
			expectedConstants: []interface{}{1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpDeclareGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			20++;
			`,
			expectedConstants: []interface{}{20, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

package vm

import (
	"cidoka/compiler"
	"cidoka/object"
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
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
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
		{"if (true) {10} else if (true) {20} else {30}", 10},
		{"if (false) {10} else if (true) {20} else {30}", 20},
		{"if (false) {10} else if (false) {20} else {30}", 30},
		{"if (true) {10} else if (true) {20}", 10},
		{"if (false) {10} else if (true) {20}", 20},
		{"if (true) {10} else if (true) {20} else if (true) {30} else {40}", 10},
		{"if (false) {10} else if (true) {20} else if (true) {30} else {40}", 20},
		{"if (false) {10} else if (false) {20} else if (true) {30} else {40}", 30},
		{"if (false) {10} else if (false) {20} else if (false) {30} else {40}", 40},
		{"if (false) {10} else if (false) {20} else if (false) {30}", Null},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}
	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", Null},
		{"[1, 2, 3][99]", Null},
		{"[1][-1]", Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", Null},
		{"{}[0]", Null},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fivePlusTen = fn() { 5 + 10; };
			fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
			let one = fn() { 1; };
			let two = fn() { 2; };
			one() + two()
			`,
			expected: 3,
		},
		{
			input: `
			let a = fn() { 1 };
			let b = fn() { a() + 1 };
			let c = fn() { b() + 1 };
			c();
			`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let earlyExit = fn() { return 99; 100; };
			earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit = fn() { return 99; return 100; };
			earlyExit();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn = fn() { };
			noReturn();
			`,
			expected: Null,
		},
		{
			input: `
			let noReturn = fn() { };
			let noReturnTwo = fn() { noReturn(); };
			noReturn();
			noReturnTwo();
			`,
			expected: Null,
		},
	}

	runVmTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let returnsOne = fn() { 1; };
			let returnsOneReturner = fn() { returnsOne; };
			returnsOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
			let returnsOneReturner = fn() {
				let returnsOne = fn() { 1; };
				returnsOne;
			};
			returnsOneReturner()();
			`,
			expected: 1,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let one = fn() { let one = 1; one };
			one();
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num;
			}
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
			expected: 97,
		},
		{
			input: `
			let a = 1;
			let foo = fn() {
				let a = 2;
				a;
			}
			foo() + a;
			`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let identity = fn(a) { a; };
			identity(4);
			`,
			expected: 4,
		},
		{
			input: `
			let sum = fn(a, b) { a + b; };
			sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				c;
			};
			sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				c;
			};
			sum(1, 2) + sum(3, 4);
			`,
			expected: 10,
		},
		{
			input: `
			let sum = fn(a, b) {
				let c = a + b;
				c;
			};
			let outer = fn() {
				sum(1, 2) + sum(3, 4);
			};
			outer();
			`,
			expected: 10,
		},
		{
			input: `
			let globalNum = 10;

			let sum = fn(a, b) {
				let c = a + b;
				c + globalNum;
			};

			let outer = fn() {
				sum(1, 2) + sum(3, 4) + globalNum;
			};

			outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "wrong number of arguments. got=2, want=1",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`print("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, Null},
		{
			`first(1)`,
			&object.Error{
				Message: "argument to `first` must be ARRAY, got INTEGER",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`, Null},
		{
			`last(1)`,
			&object.Error{
				Message: "argument to `last` must be ARRAY, got INTEGER",
			},
		},
		{`tail([1, 2, 3])`, []int{2, 3}},
		{`tail([])`, Null},
		{`push([], 1)`, []int{1}},
		{
			`push(1, 1)`,
			&object.Error{
				Message: "argument to `push` must be ARRAY, got INTEGER",
			},
		},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let newClosure = fn(a) {
				fn() { a; };
			};
			let closure = newClosure(99);
			closure();
			`,
			expected: 99,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				fn(c) { a + b + c };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
			let newAdder = fn(a, b) {
				let c = a + b;
				fn(d) { c + d };
			};
			let adder = newAdder(1, 2);
			adder(8);
			`,
			expected: 11,
		},
		{
			input: `
			let newAdderOuter = fn(a, b) {
				let c = a + b;
				fn(d) {
					let e = d + c;
					fn(f) { e + f; };
				};
			};
			let newAdderInner = newAdderOuter(1, 2)
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let a = 1;
			let newAdderOuter = fn(b) {
				fn(c) {
					fn(d) { a + b + c + d };
				};
			};
			let newAdderInner = newAdderOuter(2)
			let adder = newAdderInner(3);
			adder(8);
			`,
			expected: 14,
		},
		{
			input: `
			let newClosure = fn(a, b) {
				let one = fn() { a; };
				let two = fn() { b; };
				fn() { one() + two(); };
			};
			let closure = newClosure(9, 90);
			closure();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1);
				}
			};
			let wrapper = fn() {
				countDown(1);
			};
			wrapper();
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				};
				countDown(1);
			};
			wrapper();
			`,
			expected: 0,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fibonacci = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					if (x == 1) {
						return 1;
					} else {
						fibonacci(x - 1) + fibonacci(x - 2);
					}
				}
			};
			fibonacci(15);
			`,
			expected: 610,
		},
	}

	runVmTests(t, tests)
}

func TestIterativeFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fibonacci = fn(x) {
				let sequence = [0, 1];
				for (let i = 2; i <= x; i += 1) {
					sequence = push(sequence, sequence[i - 1] + sequence[i - 2]);
				}
				return sequence[x];
			};
			fibonacci(15);
			`,
			expected: 610,
		},
	}

	runVmTests(t, tests)
}

func TestFloatingPointNumbers(t *testing.T) {
	tests := []vmTestCase{
		{"3.14", 3.14},
		{"0.1 + 0.2", 0.3},
		{"0.1 * 2.0", 0.2},
		{"0.1 / 2.0", 0.05},
		{"0.1 + 0.1 + 0.1", 0.3},
	}

	runVmTests(t, tests)
}

func TestForLoop(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let sum = 0;
			for (let i = 0; i < 10; i = i + 1) {
				sum = sum + i;
			}
			sum;
			`,
			expected: 45,
		},
		{
			input: `
			let sum = 0;
			for (let i = 0; i < 10; i = i + 1) {
				if (i == 5) {
					break;
				}
				sum = sum + i;
			}
			sum;
			`,
			expected: 10,
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
			expected: 5,
		},
		{
			input: `
			let x = 0;
			for (let i = 0; i < 10; i += 1) {
				let sum = 0;
				for (let j = 0; j < 10; j += 1) {
					sum = sum + j;
				}
				x = sum;
			}
			x;
			`,
			expected: 45,
		},
		{
			input: `
			let sum = 0;
			for (let i = 0; i < 10; i += 1) {
				for (let j = 0; j < 10; j += 1) {
					for (let k = 0; k < 10; k += 1) {
						sum += k + j + i;
					}
				}
			}
			sum;
			`,
			expected: 13500,
		},
	}

	runVmTests(t, tests)
}

func TestCompoundAssignment(t *testing.T) {
	tests := []vmTestCase{
		{"let a = 1; a += 2; a", 3},
		{"let a = 1; a -= 2; a", -1},
		{"let a = 1; a *= 2; a", 2},
		{"let a = 4; a /= 2; a", 2},
		{"let a = 5; a %= 2; a", 1},
		{"let a = 1; a += 2 * 3; a", 7},
		{"let a = 2; a *= 5 + 2; a", 14},
	}

	runVmTests(t, tests)
}

func TestWhileLoop(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let sum = 0;
			let i = 0;
			while (i < 10) {
				sum += i;
				i += 1;
			}
			sum;
			`,
			expected: 45,
		},
		{
			input: `
			let sum = 0;
			let i = 0;
			while (i < 10) {
				if (i == 5) {
					break;
				}
				sum += i;
				i += 1;
			}
			sum;
			`,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}

func TestContinue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let sum = 0;
			for (let i = 0; i < 10; i += 1) {
				if (i == 5) {
					continue;
				}
				sum += i;
			}
			sum;
			`,
			expected: 40,
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
			expected: 400,
		},
		{
			input: `
			let sum = 0;
			let i = 0;
			while (i < 10) {
				if (i == 5) {
					i += 1;
					continue;
				}
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
			expected: 360,
		},
	}

	runVmTests(t, tests)
}

func TestArrayReassignment(t *testing.T) {
	tests := []vmTestCase{
		{"let a = [1, 2, 3]; a[0] = 4; a[0]", 4},
		{"let a = [1, 2, 3]; a[0] = 4; a[1]", 2},
		{"let a = [1, 2, 3]; a[0] = 4; a[2]", 3},
		{"let a = [1, 2, 3]; a[0] = 4; a[3]", Null},
		{"let a = [1, 2, 3]; a[0] = 4; a[-1]", Null},
		{"let a = [1, 2, 3]; a[0] = 4; a[1] = 5; a[2] = 6; a", []int{4, 5, 6}},
	}

	runVmTests(t, tests)
}

func TestPostfixOperators(t *testing.T) {
	tests := []vmTestCase{
		{"let a = 1; a++; a", 2},
		{"let a = 1; a--; a", 0},
		{"let a = 1; a++; a++; a", 3},
		{"let a = 1; a--; a--; a", -1},
		{"let x = 5; let b = x++ + x--; b;", 11},
		{"let i = 0; let result = i++ * (i + 5); result;", 6},
	}

	runVmTests(t, tests)
}

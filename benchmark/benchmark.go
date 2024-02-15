package main

import (
	"cidoka/compiler"
	"cidoka/evaluator"
	"cidoka/lexer"
	"cidoka/object"
	"cidoka/parser"
	"cidoka/vm"
	"flag"
	"fmt"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = flag.String("input", "recursive", "recursive or iterative benchmark")

var recursive = `
let fibonacci = fn(x) {
	if (x == 0) {
		0
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	}
};
fibonacci(35);
`

var iterative = `
let fibonacci = fn(x) {
	let sequence = [0, 1];
	for (let i = 2; i <= x; i += 1) {
		sequence = push(sequence, sequence[i - 1] + sequence[i - 2]);
	}
	return sequence[x];
};
fibonacci(35);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object
	var benchmark string

	if *input == "recursive" {
		benchmark = recursive
	} else {
		benchmark = iterative
	}

	l := lexer.New(benchmark)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		machine := vm.New(comp.Bytecode())

		start := time.Now()

		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s, duration=%s\n", *engine, result.Inspect(), duration)
}

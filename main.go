package main

import (
	"cidoka/compiler"
	"cidoka/evaluator"
	"cidoka/lexer"
	"cidoka/object"
	"cidoka/parser"
	"cidoka/repl"
	"cidoka/vm"
	"flag"
	"fmt"
	"os"
	"os/user"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
var input = flag.String("input", "", "input file")

func main() {
	flag.Parse()

	if *input != "" {
		runFile(*input)
		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Cidoka programming language!\n", user.Username)
	fmt.Printf("Running in %s mode\n", *engine)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, *engine)
}

func runFile(input string) {
	f, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	var result object.Object

	l := lexer.New(string(f))
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

		err = machine.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		result = evaluator.Eval(program, env)
	}

	fmt.Printf("engine=%s, result=%s\n", *engine, result.Inspect())
}

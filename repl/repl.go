package repl

import (
	"cidoka/compiler"
	"cidoka/evaluator"
	"cidoka/lexer"
	"cidoka/object"
	"cidoka/parser"
	"cidoka/vm"
	"fmt"
	"io"

	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, engine string, liner *liner.State) {
	env := object.NewEnvironment()

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		scanned, err := liner.Prompt(PROMPT)
		if err != nil {
			return
		}

		liner.AppendHistory(scanned)

		lex := lexer.New(scanned)
		parser := parser.New(lex)

		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}

		if engine == "eval" {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		} else {
			comp := compiler.NewWithState(symbolTable, constants)
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
				continue
			}

			code := comp.Bytecode()
			constants = code.Constants

			machine := vm.NewWithGlobalsStore(code, globals)
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
				continue
			}

			lastPopped := machine.LastPoppedStackElem()
			io.WriteString(out, lastPopped.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

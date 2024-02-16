package repl

import (
	"cidoka/compiler"
	"cidoka/evaluator"
	"cidoka/lexer"
	"cidoka/object"
	"cidoka/parser"
	"cidoka/token"
	"cidoka/vm"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/peterh/liner"
)

var historyFile = filepath.Join(os.TempDir(), ".cidoka_lang_history")

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, engine string) {
	var env *object.Environment
	var constants []object.Object
	var globals []object.Object
	var symbolTable *compiler.SymbolTable

	liner := liner.NewLiner()
	defer liner.Close()

	liner.SetCtrlCAborts(true)
	liner.SetCompleter(completer)

	if f, err := os.Open(historyFile); err == nil {
		liner.ReadHistory(f)
		f.Close()
	}

	if engine == "eval" {
		env = object.NewEnvironment()
	} else {
		constants = []object.Object{}
		globals = make([]object.Object, vm.GlobalsSize)

		symbolTable = compiler.NewSymbolTable()

		for i, v := range object.Builtins {
			symbolTable.DefineBuiltin(i, v.Name)
		}
	}

	for {
		scanned, err := scanInput(liner)
		if err != nil {
			writeHistory(liner)
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

func completer(line string) (c []string) {
	for _, builtin := range object.Builtins {
		if len(line) < len(builtin.Name) && builtin.Name[:len(line)] == line {
			c = append(c, builtin.Name)
		}
	}

	for keyword := range token.Keywords {
		if len(line) < len(keyword) && keyword[:len(line)] == line {
			c = append(c, keyword)
		}
	}

	return
}

func scanInput(liner *liner.State) (string, error) {
	var scanned string

	for {
		scannedLine, err := liner.Prompt(PROMPT)
		if err != nil {
			return "", err
		}

		if scannedLine == "" {
			continue
		}

		scanned += scannedLine
		break
	}

	return scanned, nil
}

func writeHistory(liner *liner.State) {
	if f, err := os.Create(historyFile); err != nil {
		fmt.Printf("Error writing history file: %s\n", err)
	} else {
		liner.WriteHistory(f)
		f.Close()
	}
}

package repl

import (
	"cidoka/ast"
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

		program, err := setupProgram(scanned)
		if err != nil {
			fmt.Fprintf(out, "Error parsing program: %s\n", err)
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
			if lastPopped.Type() != object.NULL_OBJ {
				io.WriteString(out, lastPopped.Inspect())
				io.WriteString(out, "\n")
			}
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

func RunFile(file string, engine string) {
	input, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	var result object.Object

	program, err := setupProgram(string(input))
	if err != nil {
		fmt.Printf("Error setting up program: %s\n", err)
		return
	}

	if engine == "vm" {
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("Woops! Compilation failed:\n %s\n", err)
			return
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Printf("Woops! Executing bytecode failed:\n %s\n", err)
			return
		}

		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		result = evaluator.Eval(program, env)
	}

	fmt.Printf("engine=%s, result=%s\n", engine, result.Inspect())
}

func setupProgram(input string) (*ast.Program, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stdout, p.Errors())
		return nil, fmt.Errorf("parsing error")
	}

	return program, nil
}

package repl

import (
	"boludolang/lexer"
	"boludolang/token"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan() // read a line of input from the user
		if !scanned {
			return
		}

		line := scanner.Text() // get the line of input
		lex := lexer.New(line) // create a lexer for the input
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok) // print the tokens
		}
	}
}

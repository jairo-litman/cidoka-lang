package main

import (
	"cidoka/repl"
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
		repl.RunFile(*input, *engine)
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

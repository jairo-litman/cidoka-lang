package main

import (
	"cidoka/repl"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/peterh/liner"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

func main() {
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	historyFile := filepath.Join(homeDir, ".cidoka_lang_history")

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Cidoka programming language!\n", user.Username)
	fmt.Printf("Running in %s mode\n", *engine)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout, *engine, line)

	if f, err := os.Create(historyFile); err != nil {
		fmt.Printf("Error writing history file: %s\n", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

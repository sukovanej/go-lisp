package main

import (
	"bufio"
	"fmt"
	. "github.com/sukovanej/go-lisp/interpreter"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		env := SetupMainEnv()
		runRepl(env)
	} else if len(os.Args) == 2 {
		env := SetupMainEnv()
		EvalFile(os.Args[1], env)
	} else if len(os.Args) == 3 {
		if os.Args[1] == "-i" {
			env := SetupMainEnv()
			EvalFile(os.Args[2], env)
			runRepl(env)
		}
	}
}

func runRepl(env *Env) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("GISP repl (dev version)")

	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		result := Eval(bufio.NewReader(strings.NewReader(text)), env)
		PrintCallable([]Object{result}, env)
		fmt.Println()
	}
}

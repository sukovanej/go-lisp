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
		result := EvalFile(os.Args[1], env)
		if IsErrorObject(result) {
			PrintTraceback(result.(ErrorObject))
		}
	} else if len(os.Args) == 3 {
		if os.Args[1] == "-i" {
			env := SetupMainEnv()
			result := EvalFile(os.Args[2], env)
			if IsErrorObject(result) {
				PrintTraceback(result.(ErrorObject))
			}
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
		result := Eval(bufio.NewReader(strings.NewReader(text)), env, &BufferMetaInformation{0, 0, "<REPL>"})
		if IsErrorObject(result) {
			PrintTraceback(result.(ErrorObject))
		}
		PrintCallable([]Object{result}, env)
		fmt.Println()
	}
}

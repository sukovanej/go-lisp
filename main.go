package main

import (
	. "github.com/sukovanej/go-lisp/interpreter"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		EvalFile(os.Args[1], GetMainEnv())
	}
}

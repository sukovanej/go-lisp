package main

import (
	"bufio"
	. "github.com/sukovanej/go-lisp/interpreter"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	var expectedObject, outputObject Object
	var input *bufio.Reader
	env := GetMainEnv()

	input = bufio.NewReader(strings.NewReader("1"))

	expectedObject = NumberObject{1}
	outputObject = Eval(input, env)
	if !compareSyntax(expectedSyntax, outputSyntax) {
		t.Errorf("%v != %v.", expectedSyntax, outputSyntax)
	}
}

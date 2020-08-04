package main

import (
	evl "colon/coleval"
	lex "colon/collex"
	par "colon/colparc"
	"os"

	// tol "colon/coltools"
	"fmt"
)

func testLex() {
	code := `
	5
	true
	!false
	-21
	12 + 14
	12 - 14
	12 * 14
	12.0 / 14
	1.5 + 3.6
	"hello" + " " + "world!"
	`
	interpret(code)
}

func interpret(code string) {
	lexer := lex.CreateLexerState(code)
	tokens := lexer.Lex()
	parser := par.CreateParserState(tokens, lexer.SourceLines())
	program := parser.Parse()
	parseErrors := parser.ReportErrors()

	if parseErrors[0] != "None" {
		for _, v := range parseErrors {
			fmt.Println(v)
		}
		os.Exit(22)
	}

	evl.Eval(program)

	if evl.PostEvalOutput != nil {
		for _, v := range evl.PostEvalOutput {
			fmt.Println(v.ObValue())
		}
	}
}

func main() {
	testLex()
	// tol.NewRepl()
}

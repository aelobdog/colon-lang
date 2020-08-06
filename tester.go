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
	// code := `
	// 5
	// true
	// !false
	// -21
	// 12 + 14
	// 12 - 14
	// 12 * 14
	// 12.0 / 14
	// 1.5 + 3.6
	// "hello" + " " + "world!"
	// `

	line := `
	i(12 > 2):
		i(12 > 4):
			r: 12
		
		
		:i
		r: 14
	:i
	`

	interpret(line)
}

func interpret(code string) {
	lexer := lex.CreateLexerState(code)
	tokens := lexer.Lex()
	// for _, v := range tokens {
	// 	fmt.Println(v)
	// }
	parser := par.CreateParserState(tokens, lexer.SourceLines())
	program := parser.Parse()
	// for _, v := range program.Statements {
	// 	fmt.Println(v.String())
	// }
	parseErrors := parser.ReportErrors()

	if parseErrors[0] != "None" {
		for _, v := range parseErrors {
			fmt.Println(v)
		}
		os.Exit(22)
	}

	res := evl.Eval(program)
	if res != nil {
		fmt.Println(res.ObValue())
	}
	// if evl.PostEvalOutput != nil {
	// 	for _, v := range evl.PostEvalOutput {
	// 		fmt.Println(v.ObValue())
	// 	}
	// }
}

func main() {
	testLex()
	// tol.NewRepl()
}

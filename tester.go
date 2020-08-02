package main

import (
	lex "colon/collex"
	par "colon/colparc"

	// tol "colon/coltools"
	"fmt"
)

func testLex() {
	code := `
	
	`
	// interpret(code)
	errors := interpret(code)

	if len(errors) > 0 {
		for _, v := range errors {
			fmt.Println(v)
		}
	}
}

func interpret(code string) []string {
	lexer := lex.CreateLexerState(code)
	tokens := lexer.Lex()
	// for _, v := range tokens {
	// 	fmt.Println(v)
	// }
	parser := par.CreateParserState(tokens, lexer.SourceLines())
	// parser.Parse()
	temp := parser.Parse()
	// fmt.Println(len(temp.Statements))
	for _, v := range temp.Statements {
		fmt.Println(v)
	}
	return parser.ReportErrors()
}

func main() {
	testLex()
	// tol.NewRepl()
}

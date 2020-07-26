package main

import (
	lex "colon/collex"
	par "colon/colparc"
	tol "colon/coltools"
	"fmt"
)

func testLex() {
	code := `
	12
	14
	-12
	-true
	!12
	`
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
	parser := par.CreateParserState(tokens, lexer.SourceLines())
	parser.Parse(false)
	return parser.ReportErrors()
}

func main() {
	// testLex()
	tol.NewRepl()
}

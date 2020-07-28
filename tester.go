package main

import (
	lex "colon/collex"
	par "colon/colparc"
	tol "colon/coltools"
	"fmt"
)

func testLex() {
	code := `
	-12
	12 + 13 * 24 >= 14 + 14
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

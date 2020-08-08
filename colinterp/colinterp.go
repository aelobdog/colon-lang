package colinterp

import (
	evl "colon/coleval"
	lex "colon/collex"
	obj "colon/colobj"
	par "colon/colparc"
	"fmt"
	"os"
)

// Interpret : the colon interpreter
func Interpret(code string) {
	// LEXING
	lexer := lex.CreateLexerState(code)
	tokens := lexer.Lex()

	// PARSING
	parser := par.CreateParserState(tokens, lexer.SourceLines())
	program := parser.Parse()
	parseErrors := parser.ReportErrors()

	// PARSE-ERROR CHECKING
	if parseErrors[0] != "None" {
		for _, v := range parseErrors {
			fmt.Println(v)
		}
		os.Exit(22)
	}

	// EVALUATION
	env := obj.NewEnv()
	evl.Eval(program, env)
}

/*
	DEBUG STUFF:

	LEXER :
		for _, v := range tokens {
		fmt.Println(v)
		}

	PARSER :
		for _, v := range program.Statements {
			fmt.Println(v.String())
		}
*/

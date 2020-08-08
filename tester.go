package main

import (
	evl "colon/coleval"
	lex "colon/collex"
	obj "colon/colobj"
	par "colon/colparc"
	"os"

	// tol "colon/coltools"
	"fmt"
)

func testLex() {
	env := obj.NewEnv()

	line := `
	v: hello = f(word):
		i(word == "hello"):
			word + " world!"
		:i
	:f

	hello("hello")
	`
	interpret(line, env)
	// interpret(line1, env)

}

func interpret(code string, env *obj.Env) {
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

	// env := obj.NewEnv()

	res := evl.Eval(program, env)
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

package coltools

import (
	"bufio"
	lex "colon/collex"
	par "colon/colparc"
	"fmt"
	"os"
	"strings"
)

// NewRepl : starts a new colon REPL
func NewRepl() {
	prompt := "::> "
	input := ""
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("COLON v1.0.0")
	for strings.TrimSpace(input) != "q:" {
		fmt.Print(prompt)
		input, _ = reader.ReadString('\n')
		lexer := lex.CreateLexerState(input)
		tokens := lexer.Lex()
		// @DEBUG
		// for _, t := range tokens {
		// 	if t.TokType.String() != "EOF" {
		// 		fmt.Printf("%s\t%s\t%d\n", t.TokType.String(), t.Literal, t.Line)
		// 	}
		// }
		// fmt.Println()
		parser := par.CreateParserState(tokens, lexer.SourceLines())

		// @DEBUG
		temp := parser.Parse(true)
		fmt.Println(temp.Statements[0].String())
		// parser.Parse(true)

		fmt.Println(parser.ReportErrors())
	}
}

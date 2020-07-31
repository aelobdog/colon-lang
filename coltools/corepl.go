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
	for {
		fmt.Print(prompt)
		input, _ = reader.ReadString('\n')
		if strings.TrimSpace(input) == "q:" {
			break
		}
		lexer := lex.CreateLexerState(input)
		tokens := lexer.Lex()
		// @DEBUG
		// for _, t := range tokens {
		// 	fmt.Printf("%s\t%s\t%d\n", t.TokType.String(), t.Literal, t.Line)
		// }
		// fmt.Println()
		parser := par.CreateParserState(tokens, lexer.SourceLines())

		// @DEBUG
		parser.Parse()
		// fmt.Println(temp.Statements)
		// // parser.Parse(true)

		fmt.Println(parser.ReportErrors())
	}
}

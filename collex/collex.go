package collex

import (
	tok "colon/coltok"
	"fmt"
	"os"
	"strings"
)

// Lexer : Current state of the lexer
type Lexer struct {
	Source  string
	CurrPos int
	NextPos int
	Ch      byte
	line    int
}

// CreateLexerState : to create a new lexer state and initialize it
func CreateLexerState(source string) *Lexer {
	l := Lexer{Source: source, CurrPos: 0, NextPos: 1, Ch: source[0], line: 0}
	return &l
}

// ReadChar : to read the next character
func (l *Lexer) ReadChar() {
	if l.NextPos >= len(l.Source) {
		l.Ch = 0
	} else {
		l.Ch = l.Source[l.NextPos]
	}
	l.CurrPos = l.NextPos
	l.NextPos++
}

// PeekChar : to peek at the next character
func (l *Lexer) PeekChar() byte {
	if l.NextPos >= len(l.Source) {
		return 0
	}
	return l.Source[l.NextPos]
}

// NextToken : to get the next token from the source
func (l *Lexer) NextToken() tok.Token {
	var token tok.Token
	l.consumeWhiteSpace()
	switch l.Ch {
	case '\n':
		token = tok.NewToken(tok.EOL, "", l.line)
		l.line++
	case ',':
		token = tok.NewToken(tok.COM, string(l.Ch), l.line)
	case '+':
		token = tok.NewToken(tok.PLS, string(l.Ch), l.line)
	case '-':
		token = tok.NewToken(tok.MIN, string(l.Ch), l.line)
	case '*':
		token = tok.NewToken(tok.PRD, string(l.Ch), l.line)
	case '/':
		token = tok.NewToken(tok.DIV, string(l.Ch), l.line)
	case '%':
		token = tok.NewToken(tok.REM, string(l.Ch), l.line)
	case '^':
		token = tok.NewToken(tok.POW, string(l.Ch), l.line)
	case '(':
		token = tok.NewToken(tok.LPR, string(l.Ch), l.line)
	case ')':
		token = tok.NewToken(tok.RPR, string(l.Ch), l.line)
	case '=':
		if l.PeekChar() == '=' {
			token = tok.NewToken(tok.EQL, "==", l.line)
			l.ReadChar()
		} else {
			token = tok.NewToken(tok.ASN, string(l.Ch), l.line)
		}
	case '!':
		if l.PeekChar() == '=' {
			token = tok.NewToken(tok.NEQ, "!=", l.line)
			l.ReadChar()
		} else {
			token = tok.NewToken(tok.LNT, string(l.Ch), l.line)
		}
	case '>':
		if l.PeekChar() == '=' {
			token = tok.NewToken(tok.GRE, ">=", l.line)
			l.ReadChar()
		} else {
			token = tok.NewToken(tok.GRT, string(l.Ch), l.line)
		}
	case '<':
		if l.PeekChar() == '=' {
			token = tok.NewToken(tok.LSE, "<=", l.line)
			l.ReadChar()
		} else {
			token = tok.NewToken(tok.LST, string(l.Ch), l.line)
		}
	case '&':
		token = tok.NewToken(tok.LND, string(l.Ch), l.line)
	case '|':
		token = tok.NewToken(tok.LOR, string(l.Ch), l.line)
	case 0:
		token = tok.NewToken(tok.EOF, "", l.line)
	default:
		if tok.IsDigit(l.Ch) {
			token = l.readNumber()
		} else if tok.IsLetter(l.Ch) {
			token = l.readWord()
		} else if l.Ch == '"' {
			token = l.readString()
		} else {
			token = tok.NewToken(tok.ILG, string(l.Ch), l.line)
		}
	}

	if token.TokType == tok.ILG {
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("Error on line : ", l.line, " ILLEGAL_TOKEN [", token.Literal, "] found.")
		fmt.Println("--------------------------------------------------------------------------------")
		os.Exit(65)
	}
	return token
}

func (l *Lexer) readNumber() tok.Token {
	number := string(l.Ch)
	floating := false
	for tok.IsDigit(l.PeekChar()) || l.PeekChar() == '.' {
		if l.Ch == '.' {
			if !tok.IsDigit(l.PeekChar()) {
				return tok.NewToken(tok.ILG, "", l.line)
			}
			floating = true
		}
		number = number + string(l.PeekChar())
		l.ReadChar()
	}
	if floating == true {
		return tok.NewToken(tok.FLT, number, l.line)
	}
	return tok.NewToken(tok.INT, number, l.line)
}

func (l *Lexer) readWord() tok.Token {
	word := string(l.Ch)
	for tok.IsLetter(l.PeekChar()) {
		word = word + string(l.PeekChar())
		l.ReadChar()
	}
	if tok.IsKeyword(word) {
		return tok.NewToken(tok.Keywords[word], word, l.line)
	} else if word == "true" || word == "false" ||
		word == "TRUE" || word == "FALSE" {
		return tok.NewToken(tok.BOL, word, l.line)
	}
	return tok.NewToken(tok.IDN, word, l.line)
}

func (l *Lexer) readString() tok.Token {
	var str string
	l.ReadChar()
	for l.Ch != '"' {
		if l.Ch == '\\' && l.PeekChar() == '"' {
			str = str + string(l.Ch) + string(l.PeekChar())
			l.ReadChar()
		} else {
			str = str + string(l.Ch)
		}
		l.ReadChar()

		if l.Ch == 0 {
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Printf("Error on line %d, string literal may not be closed\n", l.line)
			fmt.Println("--------------------------------------------------------------------------------")
			os.Exit(22)
		}
	}
	return tok.NewToken(tok.STR, str, l.line)
}

func (l *Lexer) consumeWhiteSpace() {
	for l.Ch == ' ' || l.Ch == '\t' || l.Ch == '\r' {
		l.ReadChar()
	}
}

// Lex : returns a list of all tokens
func (l *Lexer) Lex() []tok.Token {
	var tokens []tok.Token
	for l.Ch != 0 {
		tokens = append(tokens, l.NextToken())
		l.ReadChar()
	}
	return tokens
}

// SourceLines : returns list of lines of text in the source code (For better error handling).
func (l *Lexer) SourceLines() []string {
	return strings.Split(l.Source, "\n")
}

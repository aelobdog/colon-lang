package colparc

import (
	ast "colon/colast"
	tok "colon/coltok"
	"fmt"
	"strconv"
)

// loc : To store the contents of the source file. For better error reporting
var loc []string

// defining a couple of function types.
type (
	prefixFunc func() ast.Expression
	infixFunc  func(ast.Expression) ast.Expression
)

// Operator precedence / Binding power values
const (
	LOWEST       int = iota
	ASSIGNMENT       // Assignment operator (currently [+=, -=, *=, /=, %=, ^=] are not supported)
	COMPARISON       // All relational operators have equal precendence [==, >=, <=, >, <]
	SIMPLEARITH      // Addition and Subtraction have equal precedence [+, -]
	COMPLEXARITH     // Multiplication, Division and Remainder have equal precedence [*, /, %]
	POWER            // To the power of, or multiply be self [n] times
	PREFIX           // Unary prefix operators have the equal precedence [!, -]
	FCALL            // Function calls
)

var precedenceTable = map[tok.TokenType]int{
	// relational operators
	tok.EQL: COMPARISON,
	tok.NEQ: COMPARISON,
	tok.LST: COMPARISON,
	tok.GRT: COMPARISON,
	tok.LSE: COMPARISON,
	tok.GRE: COMPARISON,
	// arithmetic operators
	tok.PLS: SIMPLEARITH,
	tok.MIN: SIMPLEARITH,
	tok.DIV: COMPLEXARITH,
	tok.PRD: COMPLEXARITH,
	tok.REM: COMPLEXARITH,
	tok.POW: POWER,
}

// Parser : Current state of the parser
type Parser struct {
	tokens          []tok.Token
	currentToken    int
	peekedToken     int
	errors          []string
	prefixFunctions map[tok.TokenType]prefixFunc
	infixFunctions  map[tok.TokenType]infixFunc
}

// CreateParserState : to create a new parser state and initialize it
func CreateParserState(toks []tok.Token, locs []string) *Parser {
	p := &Parser{tokens: toks}
	p.currentToken = 0
	p.peekedToken = 1
	p.errors = []string{}
	loc = locs
	p.prefixFunctions = make(map[tok.TokenType]prefixFunc)
	p.infixFunctions = make(map[tok.TokenType]infixFunc)

	// registering all the valid PREFIX tokens
	p.registerPrefixFunc(tok.IDN, p.parseIdentifier)
	p.registerPrefixFunc(tok.INT, p.parseIntegerLiteral)
	p.registerPrefixFunc(tok.FLT, p.parseFloatingLiteral)
	p.registerPrefixFunc(tok.BOL, p.parseBooleanLiteral)
	p.registerPrefixFunc(tok.STR, p.parseStringLiteral)
	p.registerPrefixFunc(tok.MIN, p.parsePrefixExpression)
	p.registerPrefixFunc(tok.LNT, p.parsePrefixExpression)

	// registering all the valid INFIX tokens
	p.registerInfixFunc(tok.EQL, p.parseInfixExpression)
	p.registerInfixFunc(tok.NEQ, p.parseInfixExpression)
	p.registerInfixFunc(tok.LSE, p.parseInfixExpression)
	p.registerInfixFunc(tok.GRE, p.parseInfixExpression)
	p.registerInfixFunc(tok.LST, p.parseInfixExpression)
	p.registerInfixFunc(tok.GRT, p.parseInfixExpression)
	p.registerInfixFunc(tok.PLS, p.parseInfixExpression)
	p.registerInfixFunc(tok.MIN, p.parseInfixExpression)
	p.registerInfixFunc(tok.PRD, p.parseInfixExpression)
	p.registerInfixFunc(tok.DIV, p.parseInfixExpression)
	p.registerInfixFunc(tok.REM, p.parseInfixExpression)
	p.registerInfixFunc(tok.POW, p.parseInfixExpression)

	return p
}

// nextToken : To move the current token and the peek token ahead by one
func (p *Parser) advanceToken() {
	p.currentToken = p.peekedToken
	p.peekedToken++
}

// Errors : Returns a list of errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Parse : function that parses a stream of tokens and returns the ast produced
func (p *Parser) Parse(interpreter bool) *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokIs(tok.EOF) {

		// @DEBUG
		// fmt.Println(p.tokens[p.currentToken])

		if p.currTokIs(tok.EOL) {
			p.advanceToken()
		}

		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		if interpreter == true {
			if p.currTokIs(tok.EOL) {
				return program
			}
		}

		p.advanceToken()
	}
	return program
}

// parseStatement : to parse statements
func (p *Parser) parseStatement() ast.Statement {
	switch p.tokens[p.currentToken].TokType {
	case tok.VAR:
		return p.parseVarStatement()
	case tok.RET:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

/* --------------------------------------------------------------------------
				Statement and Expression parsing functions
  --------------------------------------------------------------------------- */

// parseVarStatement : to parse variable declaration and/or initialization statements
func (p *Parser) parseVarStatement() *ast.VarStatement {
	statement := &ast.VarStatement{Token: p.tokens[p.currentToken]}

	// @DEBUG
	// fmt.Println(p.tokens[p.currentToken])
	// fmt.Println(p.tokens[p.peekedToken])

	if !p.NextTokenIs(tok.IDN) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.tokens[p.currentToken],
		Value: p.tokens[p.currentToken].Literal,
	}

	if !p.NextTokenIs(tok.ASN) {
		return nil
	}

	// DEAL WITH RHS EXPRESSION

	for !p.currTokIs(tok.EOL) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.tokens[p.currentToken]}
	p.advanceToken()

	// DEAL WITH RETURN EXPRESSION

	for !p.currTokIs(tok.EOL) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.tokens[p.currentToken]}
	statement.Expression = p.parseExpression(LOWEST)
	if p.peekTokIs(tok.EOL) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.tokens[p.currentToken],
		Value: p.tokens[p.currentToken].Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLit := &ast.IntegerLiteral{Token: p.tokens[p.currentToken]}
	value, err := strconv.ParseInt(p.tokens[p.currentToken].Literal, 0, 64)
	if err != nil {
		p.LiteralConversionError(p.tokens[p.currentToken].Literal, "integer")
		return nil
	}
	intLit.Value = value
	return intLit
}

func (p *Parser) parseFloatingLiteral() ast.Expression {
	fltLit := &ast.FloatingLiteral{Token: p.tokens[p.currentToken]}
	value, err := strconv.ParseFloat(p.tokens[p.currentToken].Literal, 64)
	if err != nil {
		p.LiteralConversionError(p.tokens[p.currentToken].Literal, "decimal-number")
		return nil
	}
	fltLit.Value = value
	return fltLit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	boolLit := &ast.BooleanLiteral{Token: p.tokens[p.currentToken]}
	value, err := strconv.ParseBool(p.tokens[p.currentToken].Literal)
	if err != nil {
		p.LiteralConversionError(p.tokens[p.currentToken].Literal, "boolean")
	}
	boolLit.Value = value
	return boolLit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.tokens[p.currentToken],
		Value: p.tokens[p.currentToken].Literal,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFunctions[p.tokens[p.currentToken].TokType]
	if prefix == nil {
		p.UndefinedPrefixExpressionError(p.tokens[p.currentToken].TokType)
		return nil
	}
	leftExpression := prefix()
	for !p.peekTokIs(tok.EOL) && precedence < p.peekPrecedence() {
		infix := p.infixFunctions[p.tokens[p.peekedToken].TokType]
		if infix == nil {
			return leftExpression
		}
		p.advanceToken()
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.tokens[p.currentToken],
		Operator: p.tokens[p.currentToken].Literal,
	}
	p.advanceToken()
	expression.RightExpression = p.parseExpression(PREFIX)
	/*
		type checking :
		 	has [-] been used with numeric literals ?
			has [!] been used with boolean literals ?
	*/
	if expression.Operator == "-" {
		_, okIE := expression.RightExpression.(*ast.IntegerLiteral)
		_, okFE := expression.RightExpression.(*ast.FloatingLiteral)
		if !(okFE || okIE) {
			// this will get executed when RightExpression is neither an integer literal nor a floating point literal
			p.WrongDataTypeWithOperatorError("integer, floating point", "-")
		}
	} else if expression.Operator == "!" {
		_, okBE := expression.RightExpression.(*ast.BooleanLiteral)
		if !okBE {
			// this will get executed when RightExpression is neither an integer literal nor a floating point literal
			p.WrongDataTypeWithOperatorError("Boolean", "!")
		}
	}

	return expression
}

func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:          p.tokens[p.currentToken],
		Operator:       p.tokens[p.currentToken].Literal,
		LeftExpression: leftExpression,
	}

	cPrecedence := p.currPrecedence()
	p.advanceToken()
	if cPrecedence == POWER {
		cPrecedence--
	}
	expression.RightExpression = p.parseExpression(cPrecedence)
	return expression
}

/* --------------------------------------------------------------------------
							Helper functions
  --------------------------------------------------------------------------- */

// currTokIs : helper function, checks if the current token being scanned is of the desired type or not
func (p *Parser) currTokIs(tokType tok.TokenType) bool {
	return p.tokens[p.currentToken].TokType == tokType
}

// peekTokIs : helper function, checks if the next token is of the desired type or not
func (p *Parser) peekTokIs(tokType tok.TokenType) bool {
	return p.tokens[p.peekedToken].TokType == tokType
}

// NextTokenIs : moves to the next token only it is of the desired token type
func (p *Parser) NextTokenIs(t tok.TokenType) bool {
	// if p.tokens[p.peekedToken].TokType == t {
	if p.peekTokIs(t) {
		p.advanceToken()
		return true
	}
	p.ExpectedTokenError(t)
	return false
}

// currPrecedence : to get the precedence / binding power value for the current token
func (p *Parser) currPrecedence() int {
	if p, ok := precedenceTable[p.tokens[p.currentToken].TokType]; ok {
		return p
	}
	return LOWEST
}

// peekPrecedence : to get the precedence / binding power value for the token after the current token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceTable[p.tokens[p.peekedToken].TokType]; ok {
		return p
	}
	return LOWEST
}

/* --------------------------------------------------------------------------
						Operator registration functions
  --------------------------------------------------------------------------- */

// registerPrefixFunc : adds an entry to the prefix functions map
func (p *Parser) registerPrefixFunc(toktype tok.TokenType, prefixFunction prefixFunc) {
	p.prefixFunctions[toktype] = prefixFunction
}

// registerInfixFunc : adds an entry to the infix functions map
func (p *Parser) registerInfixFunc(toktype tok.TokenType, infixFunction infixFunc) {
	p.infixFunctions[toktype] = infixFunction
}

/* --------------------------------------------------------------------------
						Error formatting functions
  --------------------------------------------------------------------------- */

// ExpectedTokenError : happens when the parser is expecting a particular token but recieves some other token
func (p *Parser) ExpectedTokenError(et tok.TokenType) {
	err := fmt.Sprintf("\nError on line %d : Expecting token of type %s but got %s instead", p.tokens[p.peekedToken].Line, et.String(), p.tokens[p.peekedToken].TokType.String())
	err += "\n\n\t" + loc[p.tokens[p.peekedToken].Line]
	p.errors = append(p.errors, err)
}

// LiteralConversionError : happens when parser is unable to convert a number to the intended target data-type
func (p *Parser) LiteralConversionError(literal, target string) {
	err := fmt.Sprintf("\nError on line %d : Could not parse %q as %q", p.tokens[p.currentToken].Line, literal, target)
	err += "\n\n\t" + loc[p.tokens[p.peekedToken].Line]
	p.errors = append(p.errors, err)
}

// UndefinedPrefixExpressionError : happens when an illegal token is encountered in place of a valid prefix token in an token in an expression
// for example, if the programmer has the expression -> (* 42) -> this makes no sense because '*' is not a valid prefix token
func (p *Parser) UndefinedPrefixExpressionError(t tok.TokenType) {
	err := fmt.Sprintf("\nError on line %d : %q is not a valid 'prefix' expression/token.", p.tokens[p.currentToken].Line, t.String())
	p.errors = append(p.errors, err)
}

// WrongDataTypeWithOperatorError : happens when an operator is used with operands that the operator does not operate on
func (p *Parser) WrongDataTypeWithOperatorError(expected, operator string) {
	err := fmt.Sprintf("\nError on line %d : Operator %q can be used with operands of type %s only.", p.tokens[p.currentToken].Line, operator, expected)
	p.errors = append(p.errors, err)
}

/* --------------------------------------------------------------------------
						Error Reporting function
  --------------------------------------------------------------------------- */

// ReportErrors : for formatting error messages
func (p *Parser) ReportErrors() []string {
	if len(p.errors) > 0 {
		report := []string{}
		for _, v := range p.errors {
			report = append(report, v)
		}
		// report = append(report, "\n")
		return report
	}
	return []string{"No errors to report :)"}
}

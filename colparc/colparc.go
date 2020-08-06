package colparc

import (
	ast "colon/colast"
	tok "colon/coltok"
	"fmt"
	"strconv"
	"strings"
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
	LOGICAL          // All relational operators have equal precedence (and no short-circuiting) [&, |]
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
	// function call operator
	tok.LPR: FCALL,
	// logical operators
	tok.LND: LOGICAL,
	tok.LOR: LOGICAL,
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
	p.registerPrefixFunc(tok.LPR, p.parseGroupedExpression)
	p.registerPrefixFunc(tok.IFB, p.parseIfExpression)
	p.registerPrefixFunc(tok.FNB, p.parseFunctionExpression)
	p.registerPrefixFunc(tok.LPB, p.parseLoopStatement)

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
	p.registerInfixFunc(tok.LND, p.parseInfixExpression)
	p.registerInfixFunc(tok.LOR, p.parseInfixExpression)
	p.registerInfixFunc(tok.LPR, p.parseFunctionCall)

	return p
}

// advanceToken : To move the current token and the peek token ahead by one
func (p *Parser) advanceToken() {
	if !(p.currentToken >= len(p.tokens)-1) {
		p.currentToken = p.peekedToken
		p.peekedToken++
	}
	// do nothing if currentToken is the last token
}

// Errors : Returns a list of errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Parse : function that parses a stream of tokens and returns the ast produced
func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currTokIs(tok.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
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
	case tok.EOL:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

/* --------------------------------------------------------------------------
				Statement and Expression parsing functions
  --------------------------------------------------------------------------- */

func (p *Parser) parseVarStatement() *ast.VarStatement {
	statement := &ast.VarStatement{Token: p.tokens[p.currentToken]}
	if !p.NextTokenIs(tok.BLK) {
		return nil
	}
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
	p.advanceToken()
	statement.Value = p.parseExpression(LOWEST)
	if p.peekTokIs(tok.EOL) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.tokens[p.currentToken]}
	p.advanceToken()
	statement.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokIs(tok.EOL) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.tokens[p.currentToken]}
	statement.Expression = p.parseExpression(LOWEST)
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken()
	expression := p.parseExpression(LOWEST)
	if !p.NextTokenIs(tok.RPR) {
		p.ClosedParenMissingError()
		return nil
	}
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.tokens[p.currentToken],
	}
	if !p.NextTokenIs(tok.LPR) {
		return nil
	}
	expression.Condition = p.parseGroupedExpression()
	if !p.NextTokenIs(tok.BLK) {
		return nil
	}
	expression.IfBody = p.parseBlock(tok.IFE)
	if p.tokens[p.currentToken].TokType == tok.ELB {
		if !p.NextTokenIs(tok.BLK) {
			return nil
		}
		expression.ElseBody = p.parseBlock(tok.ELE)
	}
	return expression
}

func (p *Parser) parseBlock(endToken tok.TokenType) *ast.Block {
	block := &ast.Block{
		Token: p.tokens[p.currentToken],
	}
	p.advanceToken()
	block.Statements = []ast.Statement{}

	for p.tokens[p.currentToken].TokType != endToken && p.tokens[p.currentToken].TokType != tok.EOF {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.advanceToken()
	}
	p.advanceToken()

	return block
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	expression := &ast.FunctionExpression{
		Token: p.tokens[p.currentToken],
	}
	if !p.NextTokenIs(tok.LPR) {
		return nil
	}
	expression.Params = p.parseFunctionParameters()
	// fmt.Println(p.parseFunctionParameters())
	if !p.NextTokenIs(tok.BLK) {
		return nil
	}
	expression.FuncBody = p.parseBlock(tok.FNE)
	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	// current token is "(" So checking if the next token is ")" , denoting that function has no params
	if p.peekTokIs(tok.RPR) {
		p.advanceToken() // current token is now ")" and the next token should be ":"
		return params
	}
	p.advanceToken()

	// Extracting the first parameter
	param := &ast.Identifier{
		Token: p.tokens[p.currentToken],
		Value: p.tokens[p.currentToken].Literal,
	}
	params = append(params, param)

	// Extracting any other parameters if present
	for p.peekTokIs(tok.COM) {
		// skipping over the comma
		p.advanceToken()
		p.advanceToken()
		param := &ast.Identifier{
			Token: p.tokens[p.currentToken],
			Value: p.tokens[p.currentToken].Literal,
		}
		params = append(params, param)
	}
	if !p.NextTokenIs(tok.RPR) {
		return nil
	}
	return params
}

func (p *Parser) parseFunctionCall(leftExpr ast.Expression) ast.Expression {
	expression := &ast.FunctionCallExpression{
		Token:    p.tokens[p.currentToken],
		Function: leftExpr,
	}
	expression.Arguments = p.parseFunctionCallArguments()
	return expression
}

func (p *Parser) parseFunctionCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokIs(tok.RPR) {
		p.advanceToken()
		return args
	}
	p.advanceToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokIs(tok.COM) {
		p.advanceToken()
		p.advanceToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.NextTokenIs(tok.RPR) {
		return nil
	}
	return args
}

func (p *Parser) parseLoopStatement() ast.Expression {
	loopExpression := &ast.LoopExpression{
		Token: p.tokens[p.currentToken],
	}
	if !p.NextTokenIs(tok.LPR) {
		return nil
	}
	loopExpression.Condition = p.parseGroupedExpression()
	if !p.NextTokenIs(tok.BLK) {
		return nil
	}
	loopExpression.LoopBody = p.parseBlock(tok.LPE)
	return loopExpression
}

/* --------------------------------------------------------------------------
							Helper functions
  --------------------------------------------------------------------------- */

// currTokIs : helper function, checks if the current token being scanned is of the desired type or not
func (p *Parser) currTokIs(tokType tok.TokenType) bool {
	if p.currentToken <= len(p.tokens)-1 {
		return p.tokens[p.currentToken].TokType == tokType
	}
	//fmt.Println("Exceeded token length")
	return false
}

// peekTokIs : helper function, checks if the next token is of the desired type or not
func (p *Parser) peekTokIs(tokType tok.TokenType) bool {
	if p.currentToken <= len(p.tokens)-2 {
		return p.tokens[p.peekedToken].TokType == tokType
	}
	//fmt.Println("Exceeded token length")
	return false
}

// NextTokenIs : moves to the next token only it is of the desired token type
func (p *Parser) NextTokenIs(t tok.TokenType) bool {
	if p.peekTokIs(t) {
		p.advanceToken()
		return true
	}
	p.ExpectedTokenError(t)
	return false
}

// currPrecedence : to get the precedence / binding power value for the current token
func (p *Parser) currPrecedence() int {
	if p.currentToken <= len(p.tokens)-1 {
		if p, ok := precedenceTable[p.tokens[p.currentToken].TokType]; ok {
			return p
		}
		return LOWEST
	}
	return -1
}

// peekPrecedence : to get the precedence / binding power value for the token after the current token
func (p *Parser) peekPrecedence() int {
	if p.currentToken <= len(p.tokens)-2 {
		if p, ok := precedenceTable[p.tokens[p.peekedToken].TokType]; ok {
			return p
		}
		return LOWEST
	}
	return -1
}

func (p *Parser) removeExtraNewLines() {
	for p.tokens[p.currentToken].TokType != tok.EOF &&
		p.tokens[p.currentToken].TokType == tok.EOL &&
		p.currentToken < len(p.tokens)-1 {
		p.advanceToken()
	}
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
	err += "\n\n\t" + strings.TrimSpace(loc[p.tokens[p.currentToken].Line])
	p.errors = append(p.errors, err)
}

// ClosedParenMissingError : happens when an expression is missing a right parenthesis
func (p *Parser) ClosedParenMissingError() {
	err := fmt.Sprintf("\nError on line %d : Closing parenthesis ')' expected but not found.", p.tokens[p.currentToken].Line)
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
	err += "\n\n\t" + loc[p.tokens[p.peekedToken].Line]
	p.errors = append(p.errors, err)
}

// WrongDataTypeWithOperatorError : happens when an operator is used with operands that the operator does not operate on
func (p *Parser) WrongDataTypeWithOperatorError(expected, operator string) {
	err := fmt.Sprintf("\nError on line %d : Operator %q can be used with operands of type %s only.", p.tokens[p.currentToken].Line, operator, expected)
	err += "\n\n\t" + loc[p.tokens[p.peekedToken].Line]
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
	return []string{"None"}
}

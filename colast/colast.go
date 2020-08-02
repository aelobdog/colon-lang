package collast

import (
	"bytes"
	tok "colon/coltok"
)

// Node :
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement :
type Statement interface {
	Node
	statementNode()
}

// Expression :
type Expression interface {
	Node
	expressionNode()
}

/*-------------------------------------------------------------------*/

// Program :
type Program struct {
	Statements []Statement
}

// TokenLiteral : Program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var str bytes.Buffer
	for _, s := range p.Statements {
		str.WriteString(s.String())
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

// Identifier : for names of variables, functions, etc.
type Identifier struct {
	Token tok.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral : Identifier
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

/*-------------------------------------------------------------------*/

// IntegerLiteral : for integer literals [signed 64-bit]
type IntegerLiteral struct {
	Token tok.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

// TokenLiteral : IntegerLiteral
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

/*-------------------------------------------------------------------*/

// FloatingLiteral : for floating point literals [64-bit]
type FloatingLiteral struct {
	Token tok.Token
	Value float64
}

func (f *FloatingLiteral) expressionNode() {}

// TokenLiteral : FloatingLiteral
func (f *FloatingLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatingLiteral) String() string {
	return f.Token.Literal
}

/*-------------------------------------------------------------------*/

// BooleanLiteral : for true and false literals
type BooleanLiteral struct {
	Token tok.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}

// TokenLiteral : BooleanLiteral
func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

/*-------------------------------------------------------------------*/

// StringLiteral : to store string literals
type StringLiteral struct {
	Token tok.Token
	Value string
}

func (b *StringLiteral) expressionNode() {}

// TokenLiteral : StringLiteral
func (b *StringLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *StringLiteral) String() string {
	return b.Token.Literal
}

/*-------------------------------------------------------------------*/

// VarStatement : Variable declaration and initialization
type VarStatement struct {
	Token tok.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStatement) statementNode() {}

// TokenLiteral : VarStatement
func (v *VarStatement) TokenLiteral() string {
	return v.Token.Literal
}

func (v *VarStatement) String() string {
	var str bytes.Buffer
	str.WriteString(v.TokenLiteral() + " " + v.Name.String() + " = ")
	if v.Value != nil {
		str.WriteString(v.Value.String())
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

// ReturnStatement : Returning expressions from functions
type ReturnStatement struct {
	Token       tok.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

// TokenLiteral : ReturnStatement
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var str bytes.Buffer
	str.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		str.WriteString(r.ReturnValue.String())
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

// ExpressionStatement : Expressions used as statements
type ExpressionStatement struct {
	Token      tok.Token // holds the first token in the expression-statement
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

// TokenLiteral : ExpressionStatement
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

/*-------------------------------------------------------------------*/

// PrefixExpression : Expressions used as statements
type PrefixExpression struct {
	Token           tok.Token // holds the prefix token
	Operator        string
	RightExpression Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral : ExpressionStatement
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var str bytes.Buffer
	str.WriteString("(" + pe.Operator + " " + pe.RightExpression.String() + ")")
	return str.String()
}

/*-------------------------------------------------------------------*/

// InfixExpression : Expressions used as statements
type InfixExpression struct {
	Token           tok.Token // holds the prefix token
	Operator        string
	LeftExpression  Expression
	RightExpression Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral : ExpressionStatement
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var str bytes.Buffer
	str.WriteString("(" + ie.LeftExpression.String() + " " + ie.Operator + " " + ie.RightExpression.String() + ")")
	return str.String()
}

/*-------------------------------------------------------------------*/

// IfExpression : If-else conditional expressions
type IfExpression struct {
	Token     tok.Token // holds the [i] token
	Condition Expression
	IfBody    *Block
	ElseBody  *Block
}

func (ife *IfExpression) expressionNode() {}

// TokenLiteral : IfExpression
func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

func (ife *IfExpression) String() string {
	var str bytes.Buffer
	str.WriteString("\nIF (begin) :")
	str.WriteString("\nCONDITION : " + ife.Condition.String())
	str.WriteString("\nBODY (if) :\n")
	str.WriteString(ife.IfBody.String())
	str.WriteString("\nIF (end)")
	if ife.ElseBody != nil {
		str.WriteString("\nELSE (begin) :")
		str.WriteString("\nBODY (else) :\n")
		str.WriteString(ife.ElseBody.String())
		str.WriteString("\nELSE (end)")
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

// Block : to represent blocks of data such as those within if-else, loops, functions
type Block struct {
	Token      tok.Token // holds the [:] block indicator token
	Statements []Statement
}

func (b *Block) statementNode() {}

// TokenLiteral : Block
func (b *Block) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Block) String() string {
	var str bytes.Buffer
	for _, v := range b.Statements {
		str.WriteString(v.String() + "\n")
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

// FunctionExpression : function definition expressions
type FunctionExpression struct {
	Token    tok.Token // the [f] token
	Params   []*Identifier
	FuncBody *Block
}

func (f *FunctionExpression) expressionNode() {}

// TokenLiteral : FunctionExpression
func (f *FunctionExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionExpression) String() string {
	var str bytes.Buffer
	str.WriteString("\nFUNC (begin):")
	str.WriteString("\nPARAMS : ")
	for _, v := range f.Params {
		str.WriteString(v.String() + ", ")
	}
	str.WriteString("\nFUNC BODY : \n")
	str.WriteString(f.FuncBody.String())
	str.WriteString("\nFUNC (end)")
	return str.String()
}

/*-------------------------------------------------------------------*/

// FunctionCallExpression : function call expressions
type FunctionCallExpression struct {
	Token     tok.Token  // for the ( of the function call
	Function  Expression // for storing a function definition or a function name
	Arguments []Expression
}

func (fc *FunctionCallExpression) expressionNode() {}

// TokenLiteral : FunctionCallExpression
func (fc *FunctionCallExpression) TokenLiteral() string {
	return fc.Token.Literal
}

func (fc *FunctionCallExpression) String() string {
	var str bytes.Buffer
	str.WriteString(fc.Function.String())
	str.WriteString("\nArguments: ")
	for _, v := range fc.Arguments {
		str.WriteString(v.String() + ", ")
	}
	return str.String()
}

/*-------------------------------------------------------------------*/

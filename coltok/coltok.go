package coltok

// TokenType : type of token
type TokenType int

// valid token types in COLON
const (
	INT TokenType = iota // INTEGER
	FLT                  // FLOAT
	STR                  // STRING
	BOL                  // BOOLEAN

	LPR // LEFT PARENTHESIS
	RPR // RIGHT PARENTHESIS
	LSB // LEFT SQUARE BRACKET
	RSB // RIGHT SQUARE BRACKET

	PLS // PLUS
	MIN // MINUS
	PRD // PRODUCT
	DIV // DIVISION
	REM // REMAINDER
	POW // POWER

	EQL // EQUAL
	NEQ // NOT EQUAL
	GRT // GREATER
	LST // LESSER
	GRE // GEATER OR EQUAL
	LSE // LESS OR EQUAL

	ASN // ASSIGNMENT
	COM // COMMA

	LND // LOGICAL_AND
	LOR // LOGICAL_OR
	LNT // LOGICAL_NOT

	IDN // IDENTIFIER

	VAR // VARIABLE
	IFB // IF BEGIN
	IFE // IF END
	ELB // ELSE BEGIN
	ELE // ELSE END
	LPB // LOOP BEGIN
	LPE // LOOP END
	FNB // FUNCTION BEGIN
	FNE // FUNCTION END
	RET // RETURN

	BLK // BLOCK
	EOL // END OF LINE
	EOF // END OF FILE
	ILG // ILLEGAL CHARACTER
)

// String [from tokentype]
func (t TokenType) String() string {
	switch t {
	case INT:
		return "INTEGER"
	case FLT:
		return "FLOATING"
	case STR:
		return "STRING"
	case BOL:
		return "BOOLEAN"
	case LPR:
		return "LEFT_PARENTHESES"
	case RPR:
		return "RIGHT_PARENTHESES"
	case PLS:
		return "PLUS"
	case MIN:
		return "MINUS"
	case PRD:
		return "PRODUCT"
	case DIV:
		return "DIVISION"
	case REM:
		return "REMAINDER"
	case POW:
		return "POWER"
	case EQL:
		return "EQUAL"
	case NEQ:
		return "NOT_EQUAL"
	case GRT:
		return "GREATER_THAN"
	case LST:
		return "LESS_THAN"
	case GRE:
		return "GREATER_EQUAL"
	case LSE:
		return "LESSER_EQUAL"
	case ASN:
		return "ASSIGNMENT"
	case LND:
		return "LOGICAL_AND"
	case LOR:
		return "LOGICAL_OR"
	case LNT:
		return "LOGICAL_NOT"
	case IDN:
		return "IDENTIFIER"
	case VAR:
		return "VARIABLE"
	case IFB:
		return "BEGIN: IF"
	case IFE:
		return "END: IF"
	case ELB:
		return "BEGIN: ELSE"
	case ELE:
		return "END: ELSE"
	case LPB:
		return "LPB"
	case LPE:
		return "LPE"
	case FNB:
		return "BEGIN: FUNCTION"
	case FNE:
		return "END FUNCTION"
	case RET:
		return "RETURN"
	case COM:
		return "COMMA"
	case LSB:
		return "LEFT SQ BRACKET"
	case RSB:
		return "RIGHT SQ BRACKET"
	case EOF:
		return "EOF"
	case EOL:
		return "EOL"
	case BLK:
		return "BLOCK"
	}
	return "ILLEGAL_TOKEN"
}

// Keywords map
var Keywords map[string]TokenType = map[string]TokenType{
	"v":  VAR,
	"i":  IFB,
	":i": IFE,
	"e":  ELB,
	":e": ELE,
	"l":  LPB,
	":l": LPE,
	"f":  FNB,
	":f": FNE,
	"r":  RET,
}

// Token : properties
type Token struct {
	TokType TokenType
	Literal string
	Line    int
}

// NewToken : to assemble a token 'object'
func NewToken(tokenType TokenType, lit string, line int) Token {
	return Token{TokType: tokenType, Literal: lit, Line: line}
}

// IsDigit : to check if the current character is a digit
func IsDigit(val byte) bool {
	if val >= '0' && val <= '9' {
		return true
	}
	return false
}

// IsLetter : to check if the current character is a letter
func IsLetter(val byte) bool {
	if (val >= 'a' && val <= 'z') || (val >= 'A' && val <= 'Z') || (val == '_') {
		return true
	}
	return false
}

// IsKeyword : to check if the given string is a keyword or not
func IsKeyword(word string) bool {
	if _, ok := Keywords[word]; ok {
		return true
	}
	return false
}

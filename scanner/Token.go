package scanner

import "fmt"

type Token struct {
	Type       TokenType // Type of Token (Single/Two character/Literals/Keywords)
	Lexeme     string    // The raw character blob (extracted from source)
	Literal    any       // The actual literal value (typed) .. For example lexeme could be "123", but literal would be the int 123
	LineNumber int       // The line number which this was found/parsed
}

func NewToken(tType TokenType, lexeme string, literal any, lineNumber int) *Token {
	return &Token{
		Type:       tType,
		Lexeme:     lexeme,
		Literal:    literal,
		LineNumber: lineNumber,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}

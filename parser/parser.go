package parser

import (
	"errors"
	"fmt"

	"github.com/sunguru98/glox/lib"
	"github.com/sunguru98/glox/scanner"
)

type Parser struct {
	tokens  []*scanner.Token // List of tokens scanned
	current int              // Pointer to next token
}

// ------------------------ INIT FUNCTION -----------------------------
func InitParser(tokens []*scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// --------------------------------------------------------------------

// ------------------------- UTILS ------------------------------------

func (p *Parser) checkTokenType(tType scanner.TokenType) bool {
	// If current index is pointing to EOF
	// Then it's not going to match any token
	if p.isCurrentEOF() {
		return false
	}

	// If not, check the type
	peekedToken := p.peek()
	return peekedToken.Type == tType
}

func (p *Parser) consume(tType scanner.TokenType, message string) (*scanner.Token, error) {
	isTokenMatching := p.checkTokenType(tType)
	if isTokenMatching {
		// If the token we are expecting to consume matches with
		// the token, current index is pointing, consume and move
		return p.consumeTokenAndAdvance(), nil
	}

	currentToken := p.peek()
	return nil, p.error(currentToken, message)
}

func (p *Parser) consumeTokenAndAdvance() *scanner.Token {
	// Check if the current index is not at the end
	if !p.isCurrentEOF() {
		// If not at the end, consume the index by 1
		p.current++
	}

	previousToken := p.peekPrevious()
	return previousToken
}

func (p *Parser) matchTokenAndAdvance(tTypes ...scanner.TokenType) bool {
	// Iterate through the tokens passed
	for _, tType := range tTypes {
		// Check if current iteration's token type
		// matches with current index
		isTokenMatching := p.checkTokenType(tType)
		if isTokenMatching {
			// If yes, consume it and move current index
			// And return
			p.consumeTokenAndAdvance()
			return true
		}
	}

	// Else, none of the passed token types match with current index
	// Returning false
	return false
}

func (p *Parser) peek() *scanner.Token {
	// Fetches the token pointed by current index
	return p.tokens[p.current]
}

func (p *Parser) peekPrevious() *scanner.Token {
	// Fetches the token pointed by current index - 1
	return p.tokens[p.current-1]
}

func (p *Parser) isCurrentEOF() bool {
	// Checks if the current index's token is an EOF
	peekedToken := p.peek()
	return peekedToken.Type == scanner.EOF
}

func (p *Parser) error(token *scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		// Error if token is of EOF type (with line number)
		lib.Report(token.LineNumber, " at end", message)
	} else {
		// Else we mention the raw lexeme (with line number)
		lib.Report(token.LineNumber, " at '"+token.Lexeme+"'", message)
	}

	return errors.New("Parser Error")
}

// ------------------------------------------------------------------

// ------------------------- PARSERS --------------------------------

// Functions are defined from top to bottom (lowest precedence to highest)

// Primary - Number / String / true / false / nil / grouping expression
func (p *Parser) parsePrimary() (Expression, error) {
	// Check if the current index points to a
	// 1. True token
	if p.matchTokenAndAdvance(scanner.True) {
		return CreateLiteralExpression(true), nil
	}

	// 2. False token
	if p.matchTokenAndAdvance(scanner.False) {
		return CreateLiteralExpression(false), nil
	}

	// 3. nil token
	if p.matchTokenAndAdvance(scanner.Nil) {
		return CreateLiteralExpression(nil), nil
	}

	// 4. A number or a string
	if p.matchTokenAndAdvance(scanner.String, scanner.Number) {
		// Whenever we match a token, the current index gets incremented
		// Hence technically, the token we want is previous
		// (oldCurrent = newCurrent - 1)
		previousPeekedToken := p.peekPrevious()
		return CreateLiteralExpression(previousPeekedToken.Literal), nil
	}

	// 5. A grouping expression "(expression)"
	if p.matchTokenAndAdvance(scanner.LeftParen) {
		// Call the lowest matching expression
		// Try consuming the right paren (check if exists, if not error out)
		_, err := p.consume(scanner.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}

		// Return the fetched expression with parens '(' and ')'
		return nil, nil
	}

	// Panic if none of the expected primary literals match
	panic(fmt.Sprintf("Unexpected token: %v", *p.peek()))
}

//-------------------------------------------------------------------

package parser

import (
	"errors"

	"github.com/sunguru98/glox/lib"
	s "github.com/sunguru98/glox/scanner"
)

type Parser struct {
	tokens  []s.Token // List of tokens scanned
	current int       // Pointer to next token
}

// ------------------------ PRIMARY FUNCTIONS -------------------------
func InitParser(tokens []s.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Statement, error) {
	statements := make([]Statement, 0)
	for {
		if p.isCurrentEOF() {
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

// --------------------------------------------------------------------

// ------------------------- UTILS ------------------------------------
func (p *Parser) checkTokenType(tType s.TokenType) bool {
	// If current index is pointing to EOF
	// Then it's not going to match any token
	if p.isCurrentEOF() {
		return false
	}

	// If not, check the type
	peekedToken := p.peek()
	return peekedToken.Type == tType
}

func (p *Parser) consume(tType s.TokenType, message string) (s.Token, error) {
	isTokenMatching := p.checkTokenType(tType)
	if isTokenMatching {
		// If the token we are expecting to consume matches with
		// the token, current index is pointing, consume and move
		return p.consumeTokenAndAdvance(), nil
	}

	currentToken := p.peek()
	return s.Token{}, p.error(currentToken, message)
}

func (p *Parser) consumeTokenAndAdvance() s.Token {
	// Check if the current index is not at the end
	if !p.isCurrentEOF() {
		// If not at the end, consume the index by 1
		p.current++
	}

	previousToken := p.peekPrevious()
	return previousToken
}

func (p *Parser) matchTokenAndAdvance(tTypes ...s.TokenType) bool {
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

func (p *Parser) peek() s.Token {
	// Fetches the token pointed by current index
	return p.tokens[p.current]
}

func (p *Parser) peekPrevious() s.Token {
	// Fetches the token pointed by current index - 1
	return p.tokens[p.current-1]
}

func (p *Parser) isCurrentEOF() bool {
	// Checks if the current index's token is an EOF
	peekedToken := p.peek()
	return peekedToken.Type == s.EOF
}

func (p *Parser) error(token s.Token, message string) error {
	if token.Type == s.EOF {
		// Error if token is of EOF type (with line number)
		lib.Report(token.LineNumber, " at end", message)
	} else {
		// Else we mention the raw lexeme (with line number)
		lib.Report(token.LineNumber, " at '"+token.Lexeme+"'", message)
	}

	return errors.New("Parser Error")
}

func (p *Parser) synchronizeFromError() {
	// When an error is met in any of the parser levels, we try to move on to the next token
	p.consumeTokenAndAdvance()

	// Once we move past the errored token, there could be two cases
	// 1. The errored line might be ended (through a semicolon)
	// 2. Or the error might have happened in mid statement, and another token exists
	// We check the both cases inside the loop

	for {
		// Check if we haven't met EOF
		if p.isCurrentEOF() {
			return
		}

		previousPeekedToken := p.peekPrevious()
		// Check if the token we just consumed is either a semicolon
		if previousPeekedToken.Type == s.Semicolon {
			return
		}

		// Or if the current token is one of the following
		currentToken := p.peek()
		switch currentToken.Type {
		case s.Class, s.Fun, s.Var, s.For, s.If, s.While, s.Print, s.Return:
			return
		}

		// We keep iterating till we either see semicolon/EOF/one of the tokens in switch
		p.consumeTokenAndAdvance()
	}
}

// ------------------------------------------------------------------

// ------------------------ STATEMENTS -----------------------------

func (p *Parser) parsePrintStatement() (Statement, error) {
	printValue, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(s.Semicolon, "Expecting ; after expression")
	if err != nil {
		return nil, err
	}

	return CreateNewPrintSt(printValue), nil
}

func (p *Parser) parseExpressionStatement() (Statement, error) {
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(s.Semicolon, "Expecting ; after expression")
	if err != nil {
		return nil, err
	}

	return CreateNewExpressionSt(expression), nil
}

func (p *Parser) parseStatement() (Statement, error) {
	if p.matchTokenAndAdvance(s.Print) {
		return p.parsePrintStatement()
	}

	return p.parseExpressionStatement()
}

// ------------------------------------------------------------------

// ------------------------- PARSERS --------------------------------
// Functions are defined from top to bottom (lowest precedence to highest)
// Each parsing precedence has it's own grammar

// Expression - equality operand
// This is the lowest most precedence operand
// Matching with Equality matches all possible cases
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

// Equality - comparison operand (('!=', '==') comparison operand)
func (p *Parser) parseEquality() (Expression, error) {
	// Fetch the left term expression
	comparisonLeftExpression, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var equalityExpression Expression = comparisonLeftExpression

	for {
		// Loop through until token type is neither != nor ==
		isTokenMatching := p.matchTokenAndAdvance(s.BangEqual, s.EqualEqual)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Followed by the right comparison expression
		comparisonRightExpression, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		// We then conjoin the initial left comparison expression
		// The operator and the right comparison expression
		equalityExpression = CreateBinaryExpression(comparisonLeftExpression, operator, comparisonRightExpression)
	}

	return equalityExpression, nil
}

// Comparison - term operand (('>' / '>=' / '<' / '<=') term operand)
func (p *Parser) parseComparison() (Expression, error) {
	// Fetch the left term expression
	termLeftExpression, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var comparisonExpression Expression = termLeftExpression

	for {
		// Loop through until token type is neither > nor >= nor < nor <=
		isTokenMatching := p.matchTokenAndAdvance(s.Greater, s.GreaterEqual, s.Less, s.LessEqual)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Followed by the right term expression
		termRightExpression, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		// We then conjoin the initial left term expression
		// The operator and the right term expression
		comparisonExpression = CreateBinaryExpression(termLeftExpression, operator, termRightExpression)
	}

	return comparisonExpression, nil
}

// Term - factor operand (('-' / '+') factor operand)
func (p *Parser) parseTerm() (Expression, error) {
	// Fetch the left factor expression
	factorLeftExpression, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var termExpression Expression = factorLeftExpression

	for {
		// Loop through until token type is neither - nor +
		isTokenMatching := p.matchTokenAndAdvance(s.Minus, s.Plus)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Followed by the right factor expression
		factorRightExpression, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		// We then conjoin the initial left factor expression
		// The operator and the right factor expression
		termExpression = CreateBinaryExpression(factorLeftExpression, operator, factorRightExpression)
	}

	return termExpression, nil
}

// Same as Term, except involves / and *
// Factor - Unary operand (('/' / '*') Unary operand)
func (p *Parser) parseFactor() (Expression, error) {
	// Fetch the expression for the left hand operand (unary)
	unaryLeftExpression, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var factorExpression Expression = unaryLeftExpression

	for {
		// Loop through until token type is neither / nor *
		isTokenMatching := p.matchTokenAndAdvance(s.Slash, s.Star)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Followed by the right unary operand
		unaryRightExpression, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		// We then conjoin the initial left unary expression
		// The operator and the right unary expression
		factorExpression = CreateBinaryExpression(unaryLeftExpression, operator, unaryRightExpression)
	}

	return factorExpression, nil
}

// Unary - Logical NOT or Minus [ (!/-) Unary operand ]
func (p *Parser) parseUnary() (Expression, error) {
	// Check if the token starts with a '!' or a '-'
	isTokenMatching := p.matchTokenAndAdvance(s.Bang, s.Minus)
	if isTokenMatching {
		// Fetch the specific operator token
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()
		// There is a possibility of chained unary expressions, hence we recurse
		// And fetch the corresponding expression
		expression, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		// We then emit this as a unary expression
		unaryExpression := CreateUnaryExpression(operator, expression)
		return unaryExpression, nil
	}

	// Else it's probably one the primary expressions mentioned below
	return p.parsePrimary()
}

// Primary - Number / String / true / false / nil / grouping expression / IDENTIFIER
func (p *Parser) parsePrimary() (Expression, error) {
	// Check if the current index points to a
	// 1. True token
	if p.matchTokenAndAdvance(s.True) {
		return CreateLiteralExpression(true), nil
	}

	// 2. False token
	if p.matchTokenAndAdvance(s.False) {
		return CreateLiteralExpression(false), nil
	}

	// 3. nil token
	if p.matchTokenAndAdvance(s.Nil) {
		return CreateLiteralExpression(nil), nil
	}

	// 4. A number or a string
	if p.matchTokenAndAdvance(s.String, s.Number) {
		// Whenever we match a token, the current index gets incremented
		// Hence technically, the token we want is previous
		// (oldCurrent = newCurrent - 1)
		previousPeekedToken := p.peekPrevious()
		return CreateLiteralExpression(previousPeekedToken.Literal), nil
	}

	// 5. A grouping expression "(expression)"
	// starts with left parentheses
	if p.matchTokenAndAdvance(s.LeftParen) {
		// Call the lowest matching expression
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// Try consuming the right paren (check if exists, if not error out)
		_, err = p.consume(s.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}

		// Return the fetched expression with parens '(' and ')'
		// As a grouped expression

		groupedExpression := CreateGroupingExpression(expression)
		return groupedExpression, nil
	}

	// We report an error if none of the above tokens match.
	// This means, we are at the lowest level, and still none matched to form an expression
	// Hence we report it.
	currentToken := p.peek()
	return nil, p.error(currentToken, "Expecting an expression")
}

//-------------------------------------------------------------------

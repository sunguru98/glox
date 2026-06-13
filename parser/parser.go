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

func (p *Parser) Parse() []Statement {
	// Create a statement array to parse every line
	statements := make([]Statement, 0)

	for {
		// The parser evaluates the file/command in REPL
		// Until EOF is met (a token of EOF)
		if p.isCurrentEOF() {
			break
		}

		// We then parse every declaration met
		declaration := p.parseDeclaration()
		// And add it to our parsed statements list
		statements = append(statements, declaration)
	}

	// Which is then returned
	return statements
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

func (p *Parser) parseIfStatement() (Statement, error) {
	// Check for a left paren token after 'if'
	_, err := p.consume(s.LeftParen, "Expect ( after if")
	if err != nil {
		return nil, err
	}

	// Parse the condition
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Check for a right paren after 'if'
	_, err = p.consume(s.RightParen, "Expect ) after if condition")
	if err != nil {
		return nil, err
	}

	// We then parse the if block statement
	thenBlockStatement, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// There can be an optional else block, hence we init with nil
	var elseBlockStatement Statement = nil
	if p.matchTokenAndAdvance(s.Else) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		elseBlockStatement = statement
	}

	// The parsed if statement is then created
	ifStatement := CreateIfSt(condition, thenBlockStatement, elseBlockStatement)
	return ifStatement, nil
}

func (p *Parser) parsePrintStatement() (Statement, error) {
	// The current index moves beyond 'print'
	// Hence we parse the expression
	printValue, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// We check if the statement ends with ';'
	_, err = p.consume(s.Semicolon, "Expect ; after expression")
	if err != nil {
		return nil, err
	}

	// If all succeeds, we create a Print statement node
	return CreatePrintSt(printValue), nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	// Consuming the left paren
	_, err := p.consume(s.LeftParen, "Expect '(' after 'while'")
	if err != nil {
		return nil, err
	}

	// Consuming the while condition expression
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Consuming the right paren
	_, err = p.consume(s.RightParen, "Expect ')' after 'while' condition")
	if err != nil {
		return nil, err
	}

	// Parsing the block statement
	blockStatement, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Creating a new while statement
	whileStatement := CreateWhileSt(condition, blockStatement)
	return whileStatement, nil
}

func (p *Parser) parseForStatement() (Statement, error) {
	// Consuming the left paren
	_, err := p.consume(s.LeftParen, "Expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	// A for loop 'usually' has three steps inside paren
	// 1. Init variable to a value
	var initializer Statement

	// If it starts with a semicolon, then there is no initializer
	if p.matchTokenAndAdvance(s.Semicolon) {
		initializer = nil
	}

	// If it starts with a "var" keyword
	// Then the initializer is a variable declaration
	if p.matchTokenAndAdvance(s.Var) {
		initializer, err = p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		// Else it's a regular expression
		// Meaning the variable declaration could be outside the loop
		// and just the initialization happens here
		initializer, err = p.parseExpressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// 2. Checking for a condition
	var condition Expression = nil
	// If there is a semicolon immediately after 1.
	// Then the condition expression is skipped
	// If not, we parse the condition
	if !p.checkTokenType(s.Semicolon) {
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// And then later, we consume the semicolon
	_, err = p.consume(s.Semicolon, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	// 3. Post condition updates to the variable
	var postConditionExpression Expression = nil
	// If there is a right paren immediately after 2.
	// Then the increment expression is skipped
	// If not, we parse the incrementing expression
	if !p.checkTokenType(s.RightParen) {
		postConditionExpression, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// And then later, we consume the right paren
	_, err = p.consume(s.RightParen, "Expect ')' after for clauses")
	if err != nil {
		return nil, err
	}

	// We now have individually parsed
	// 1. Initializer expression
	// 2. Condition expression
	// 3. Post condition expression
	// The for loop is simply a disguised while loop
	forStatement, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Usually a while loop's condition variable
	// Would be changed at the end of the iteration
	// Hence we consider the same here, if the
	// post condition expression is present
	if postConditionExpression != nil {
		// We make the already existing block statements execute first
		// Followed by the post condition expression statement
		statementsToExecute := []Statement{forStatement, CreateExpressionSt(postConditionExpression)}
		// And bundled together
		forStatement = CreateBlockSt(statementsToExecute)
	}

	// If the condition expression doesn't exist
	// That simply means a 'while(true)' a.k.a infinite loop
	if condition == nil {
		// We fill the condition value to be looping forever instead
		condition = CreateLiteralExpression(true)
	}

	// We have the bare minimum requirement for while (condition/body block)
	forStatement = CreateWhileSt(condition, forStatement)

	// Finally, if the initializer expression exists
	// That expression statement runs first before the above bundled forStatement
	if initializer != nil {
		// Hence we bundle again as a block, with the initializer expression occurring first
		statementsToExecute := []Statement{initializer, forStatement}
		forStatement = CreateBlockSt(statementsToExecute)
	}

	return forStatement, nil
}

func (p *Parser) parseExpressionStatement() (Statement, error) {
	// We parse the expression
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// We check if the statement ends with ';'
	_, err = p.consume(s.Semicolon, "Expect ; after expression")
	if err != nil {
		return nil, err
	}

	// If all succeeds, we create a generic statement node
	return CreateExpressionSt(expression), nil
}

func (p *Parser) parseBlockStatement() ([]Statement, error) {
	// Creating the statements array
	statements := make([]Statement, 0)
	for {
		// Until the Right brace is met '}' or EOF token is met
		if p.checkTokenType(s.RightBrace) || p.isCurrentEOF() {
			break
		}

		// We store declarations that are met inside the block
		declarationStatement := p.parseDeclaration()
		statements = append(statements, declarationStatement)
	}

	// The above loop could have terminated due to EOF without Right brace
	// Hence we check that
	_, err := p.consume(s.RightBrace, "Expect '}' after the block")
	if err != nil {
		return nil, err
	}

	// If the loop gracefully ended due to }
	// We then return the parsed block statements
	return statements, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	// If the statement starts with the 'if' keyword
	// Consider it as an if statement
	if p.matchTokenAndAdvance(s.If) {
		return p.parseIfStatement()
	}

	// If the statement starts with the 'print' keyword
	// Consider it as a print statement
	if p.matchTokenAndAdvance(s.Print) {
		return p.parsePrintStatement()
	}

	// If the statement starts with the 'while' keyword
	// Consider it as a while statement
	if p.matchTokenAndAdvance(s.While) {
		return p.parseWhileStatement()
	}

	// If the statement starts with the 'for' keyword
	// Consider it as a for statement
	if p.matchTokenAndAdvance(s.For) {
		return p.parseForStatement()
	}

	// If the statement starts with a left brace '{'
	// It would define a block declaration environment
	if p.matchTokenAndAdvance(s.LeftBrace) {
		blockStatements, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}

		return CreateBlockSt(blockStatements), nil
	}

	// Else it's a generic expression statement
	// Expression with a semicolon
	return p.parseExpressionStatement()
}

// ------------------------------------------------------------------

// ------------------------ DECLARATION -----------------------------

func (p *Parser) parseVariableDeclaration() (Statement, error) {
	// Current index points after the var keyword
	// Hence we parse the variable name/IDENTIFIER
	identifierToken, err := p.consume(s.Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	// We initially assume the variable is uninitialized
	var initializer Expression = nil
	// Hence we check if there exists an equal sign
	if p.matchTokenAndAdvance(s.Equal) {
		expressionResult, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// If it exists, that means there exists an initializer
		// Which equates to the result of the expression
		initializer = expressionResult
	}

	// We then consume the semicolon character to finish the var statement
	_, err = p.consume(s.Semicolon, "Expect ; after variable declaration")
	if err != nil {
		return nil, err
	}

	// Followed by generating a statement node
	return CreateVariableSt(identifierToken, initializer), nil
}

func (p *Parser) parseDeclaration() Statement {
	// We first try checking for a 'var' keyword
	if p.matchTokenAndAdvance(s.Var) {
		// If yes, we parse it as a variable statement
		variableStatement, err := p.parseVariableDeclaration()
		if err != nil {
			// If there exists an error, we synchronize
			// Check p.synchronizeFromError
			p.synchronizeFromError()
			return nil
		}

		// Else we return the evaluated variable statement
		return variableStatement
	}

	// Every declaration statement is a subset of a statement
	// Hence we parse that
	statement, err := p.parseStatement()
	if err != nil {
		p.synchronizeFromError()
		return nil
	}

	return statement
}

// ------------------------------------------------------------------

// ------------------------- EXPRESSION -----------------------------
// Functions are defined from top to bottom (lowest precedence to highest)
// Each parsing precedence has it's own grammar

// Expression - assignment operand
// This is the lowest most precedence operand
// Matching with Assignment matches all possible cases
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAssignment()
}

// Assignment - IDENTIFER = (assignment | logic_or)
// Lox bases assignments as expressions and not statements (just like C)
func (p *Parser) parseAssignment() (Expression, error) {
	// Fetch the left term expression
	equalityLeftExpression, err := p.parseLogicOr()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var assignmentExpression Expression = equalityLeftExpression

	// Check if there is an assignment (ie: an equal sign)
	if p.matchTokenAndAdvance(s.Equal) {
		equalToken := p.peekPrevious()

		// Fetch the right term expression (value)
		// Since assignments are right associative, we recursively call assignment()
		valueExpression, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}

		// Check if the assignment expression is of type Variable
		// That is, a variable that is already declared/initialized with a value
		variableExpression, ok := assignmentExpression.(*Variable)
		if !ok {
			// If not, then the value cannot be assigned
			return nil, p.error(equalToken, "Invalid assignment target.")
		}

		// Else, fetch the variable name and return as an assignment expression
		// With the new value
		variableName := variableExpression.Name
		return CreateAssignmentExpression(variableName, valueExpression), nil
	}

	return assignmentExpression, nil
}

// Logic OR (||) - logic_and ('or' logic_and)*
// ()* means that section can repeat multiple times
func (p *Parser) parseLogicOr() (Expression, error) {
	// Fetch the left term expression
	logicalAndLeft, err := p.parseLogicAnd()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var logicalOrExpression Expression = logicalAndLeft

	for {
		// Loop through until token type is not || (or)
		isTokenMatching := p.matchTokenAndAdvance(s.Or)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Fetch the right term expression
		logicalAndRight, err := p.parseLogicAnd()
		if err != nil {
			return nil, err
		}

		// Re-initialize the new Logical OR expression
		logicalOrExpression = CreateLogicalExpression(logicalAndLeft, operator, logicalAndRight)
	}

	return logicalOrExpression, nil
}

// Logic AND (&&) - equality ('and' equality)*
// ()* means that section can repeat multiple times
func (p *Parser) parseLogicAnd() (Expression, error) {
	// Fetch the left term expression
	equalityLeft, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	// The expression to return
	var logicalAndExpression Expression = equalityLeft

	for {
		// Loop through until token type is not && (and)
		isTokenMatching := p.matchTokenAndAdvance(s.And)
		if !isTokenMatching {
			break
		}

		// Fetch the specific operator
		// Previous because matchTokenAndAdvance moves the current index
		// Hence oldCurrent (the operator token) = newCurrent - 1
		operator := p.peekPrevious()

		// Fetch the right term expression
		equalityRight, err := p.parseEquality()
		if err != nil {
			return nil, err
		}

		// Re-initialize the new Logical AND expression
		logicalAndExpression = CreateLogicalExpression(equalityLeft, operator, equalityRight)
	}

	return logicalAndExpression, nil
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

	// 5. An Identifier
	if p.matchTokenAndAdvance(s.Identifier) {
		// Every identifier comes with a variable statement
		previousPeekedToken := p.peekPrevious()
		return CreateVariableExpression(previousPeekedToken), nil
	}

	// 6. A grouping expression "(expression)"
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
	return nil, p.error(currentToken, "Expect an expression")
}

//-------------------------------------------------------------------

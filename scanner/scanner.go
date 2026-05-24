package scanner

import (
	"fmt"
	"strconv"

	lib "github.com/sunguru98/glox/lib"
)

var Keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

type Scanner struct {
	source  string  // The raw source code
	tokens  []Token // The parsed tokens from lexemes
	start   int     // Start is offset of the first character in the current scanning lexeme
	current int     // Current is the current character offset in the whole source code
	line    int     // Line is the line number from the source code
}

func InitScanner(source string) *Scanner {
	// Creating a scanner
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanToken() {
	// Get the current character
	currentCharacter := s.consumeCharAndAdvance()

	// And parse the token based on the fetched character
	switch currentCharacter {
	// Single character lexemes
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(Comma)
	case '.':
		s.addToken(Dot)
	case '-':
		s.addToken(Minus)
	case '+':
		s.addToken(Plus)
	case ';':
		s.addToken(Semicolon)
	case '*':
		s.addToken(Star)

	// The following lexemes can either be single or double
	case '!':
		var tType TokenType = Bang
		if s.matchCharAndAdvance('=') {
			tType = BangEqual
		}
		s.addToken(tType)
	case '=':
		var tType TokenType = Equal
		if s.matchCharAndAdvance('=') {
			tType = EqualEqual
		}
		s.addToken(tType)
	case '<':
		var tType TokenType = Less
		if s.matchCharAndAdvance('=') {
			tType = LessEqual
		}
		s.addToken(tType)
	case '>':
		var tType TokenType = Greater
		if s.matchCharAndAdvance('=') {
			tType = GreaterEqual
		}
		s.addToken(tType)

	// '/' can either be division or comment '//'.
	case '/':
		if s.matchCharAndAdvance('/') {
			// The whole line is just a comment, hence advance till a new line ('\n') is met
			for {
				// If either we meet with the new line character
				// Or the end of file, break.
				// For new line, we have a seperate case handling it
				// Hence the '\n' is not consumed but peeked instead
				if s.peek() == '\n' || s.isAtEnd() {
					break
				}

				// Until then, we advance and consume those characters through current index increment
				s.consumeCharAndAdvance()
			}
		} else {
			// Else it's normal division lexeme
			s.addToken(Slash)
		}

	// No tokens are added for anything whitespace related
	case ' ', '\r', '\t':
		break

	// For new line, the line value is incremented
	case '\n':
		s.line++

	// String literals
	case '"':
		s.matchString()

	default:
		// Check if it's a numeric literal
		if s.isDigit(currentCharacter) {
			s.matchDigit()
			break
		}

		// Or if it's an alpha (alphabet/underscore)
		if s.isAlpha(currentCharacter) {
			s.matchIdentifider()
			break
		}

		// When none of the assumed character combinations match, invoke an error
		lib.Error(s.line, "Unexpected character.")
	}
}

func (s *Scanner) ScanTokens() []Token {
	for {
		// If reached end of source break
		if s.isAtEnd() {
			break
		}

		// We start from the beginning character of the lexeme. s.current is where the beginning is
		s.start = s.current
		// Scan and add token to list
		s.scanToken()
	}

	// Since we've reached EOF, need to add EOF token
	eofToken := NewToken(EOF, "", nil, s.line)
	s.tokens = append(s.tokens, *eofToken)

	return s.tokens
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	// We fetch the raw character blob (lexeme) and mark the corresponding token type and the parsed literal (if any)
	lexeme := s.source[s.start:s.current]
	// Creating and adding it into the token list
	token := NewToken(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, *token)
}

func (s *Scanner) isAtEnd() bool {
	// Check if the cursor of "current" offset has reached end of source code
	return s.current >= len(s.source)
}

func (s *Scanner) peek() byte {
	// If we've reached EOF, then return null terminator
	if s.isAtEnd() {
		return '\x00'
	}

	// Otherwise peek (look at current index's character of the source)
	peekedCharacter := s.source[s.current]
	return peekedCharacter
}

func (s *Scanner) peekNext() byte {
	// Next character is current index + 1
	// We aren't mutating s.current, since this is peek, not consume
	nextCharacter := s.current + 1
	// Same like peek() .. Return null terminator if EOF
	if nextCharacter >= len(s.source) {
		return '\x00'
	}

	peekedCharacter := s.source[nextCharacter]
	return peekedCharacter
}

func (s *Scanner) consumeCharAndAdvance() byte {
	// Ingest/Fetch/Consume the current character and move the cursor by 1 byte/character
	character := s.source[s.current]
	s.current++

	return character
}

func (s *Scanner) matchCharAndAdvance(expectedChar byte) bool {
	// If the current cursor is already at the source end
	// then there are no more characters to check
	if s.isAtEnd() {
		return false
	}

	// Check if the current character matches the expected character
	currentCharacter := s.source[s.current]
	if currentCharacter != expectedChar {
		return false
	}

	// If it does, consume the character and move forward by 1 character/byte
	s.current++

	return true
}

func (s *Scanner) isAlpha(charToBeChecked byte) bool {
	// The character should either be within 'a' to 'z'
	// Or 'A' - 'Z' or an underscore (since a variable/identifier can have underscores)
	return (charToBeChecked >= 'a' && charToBeChecked <= 'z') ||
		(charToBeChecked >= 'A' && charToBeChecked <= 'Z') ||
		charToBeChecked == '_'
}

func (s *Scanner) matchString() {
	for {
		// The loop advances until the current index is at end of source
		// Or current index points the character '"' (closing quote for the string)
		if s.isAtEnd() || s.peek() == '"' {
			break
		}

		// The programming language supports multiline strings
		// Hence we also update the line number state accordingly
		if s.peek() == '\n' {
			s.line++
		}

		s.consumeCharAndAdvance()
	}

	// The loop might be terminated due to two causes
	// 1. If current index is at end of source, without closing quote (")
	// Which means it's not a valid string
	if s.isAtEnd() {
		lib.Error(s.line, "Unterminated string")
		return
	}

	// 2. The healthy loop termination is when it has peeked the closing quote
	// Which means, we consume '"' and advance
	s.consumeCharAndAdvance()

	// Before this function, start and current would be in the same place.
	// Now start and current are in beginning and after end of the string.
	// For an example string like "apple" .. start would be at '"' (open quote)
	// And current would be after the closing quote. (since we consumed it)
	// Hence we have to trim the quotes for the actual value
	strStart := s.start + 1     // Go one char after open quote
	strCurrent := s.current - 1 // Go one char before, to land exactly at the closing quote (since end index is not inclusive)

	// Difference between literal and raw lexeme here is. Lexeme would contain with quotes
	// Literal value just captures the true type. Here it's the string value inside the quotes.
	strLiteralValue := s.source[strStart:strCurrent]

	// We then create this as a token and add it to the list
	s.addTokenWithLiteral(String, strLiteralValue)
}

func (s *Scanner) isDigit(charToBeChecked byte) bool {
	return charToBeChecked >= '0' && charToBeChecked <= '9'
}

func (s *Scanner) consumeDigits() {
	for {
		// The loop advances till the current index isn't pointing a number
		peekedCharacter := s.peek()
		if !s.isDigit(peekedCharacter) {
			break
		}

		s.consumeCharAndAdvance()
	}
}

func (s *Scanner) matchDigit() {
	// A number can be either multiple digits (1234)
	// Or with decimals (123.456)
	s.consumeDigits()

	// The loop inside `consumeDigits` might be terminated midway if current index points to '.'
	// Hence we check if that is the case and also check if current + 1 index points to a number
	currentPeekedCharacter := s.peek()
	nextPeekedCharacter := s.peekNext()
	if currentPeekedCharacter == '.' && s.isDigit(nextPeekedCharacter) {
		// Consume the decimal dot
		s.consumeCharAndAdvance()

		// Consume the remaining digits after the decimal dot
		s.consumeDigits()
	}

	// Since there are no quotes surrounding a number, we can take start and current as is
	// The programming language implicitly considers all numbers as float64 for interpreter simplicity
	// Hence we parse it as an float64 and then create a token and add it to list
	valueStr := s.source[s.start:s.current]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		lib.Error(s.line, fmt.Sprintf("Failed to format digit. Error: %v", err))
	}

	s.addTokenWithLiteral(Number, value)
}

func (s *Scanner) isAlphaNumeric(charToBeChecked byte) bool {
	// To check whether the character is either a number or alphabet or underscore
	return s.isAlpha(charToBeChecked) || s.isDigit(charToBeChecked)
}

func (s *Scanner) matchIdentifider() {
	for {
		// As long as the current index points to a valid digit/alphabet/underscore
		// We can consume the characters
		peekedCharacter := s.peek()
		// Once it's not the case, then we can break
		// To get the range of the identifier to capture (start to current)
		if !s.isAlphaNumeric(peekedCharacter) {
			break
		}

		s.consumeCharAndAdvance()
	}

	// A reserved keyword is also a type of identifier
	identifierRawStr := s.source[s.start:s.current]

	// We try to check if the consumed range is part of the keyword map
	// If yes, we create a token under the Keyword TokenType
	keywordType, ok := Keywords[identifierRawStr]
	if !ok {
		// Else it's marked as a normal Identifer
		keywordType = Identifier
	}

	// We then create a token and add to the list
	s.addToken(keywordType)
}

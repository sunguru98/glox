package scanner

import "github.com/sunguru98/glox/lib"

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

func (s *Scanner) ScanToken() {
	// Get the current character
	currentCharacter := s.FetchCharAndAdvance()

	// And parse the token based on the fetched character
	switch currentCharacter {
	case '(':
		s.AddToken(LeftParen)
	case ')':
		s.AddToken(RightParen)
	case '{':
		s.AddToken(LeftBrace)
	case '}':
		s.AddToken(RightBrace)
	case ',':
		s.AddToken(Comma)
	case '.':
		s.AddToken(Dot)
	case '-':
		s.AddToken(Minus)
	case '+':
		s.AddToken(Plus)
	case ';':
		s.AddToken(Semicolon)
	case '*':
		s.AddToken(Star)
	// When none of the assumed characters match, invoke an error
	default:
		lib.Error(s.line, "Unexpected character.")
		return
	}
}

func (s *Scanner) ScanTokens() {
	for {
		// If reached end of source break
		if s.IsAtEnd() {
			break
		}

		// We start from the beginning character of the lexeme. s.current is where the beginning is
		s.start = s.current
		// Scan and add token to list
		s.ScanToken()
	}

	// Since we've reached EOF, need to add EOF token
	eofToken := NewToken(EOF, "", nil, s.line)
	s.tokens = append(s.tokens, *eofToken)
}

func (s *Scanner) IsAtEnd() bool {
	// Check if the cursor of "current" offset has reached end of source code
	return s.current >= len(s.source)
}

func (s *Scanner) FetchCharAndAdvance() byte {
	// Ingest/Fetch/Consume the current character and move the cursor by 1 byte/character
	character := s.source[s.current]
	s.current++

	return character
}

func (s *Scanner) AddToken(tokenType TokenType) {
	s.AddTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) AddTokenWithLiteral(tokenType TokenType, literal any) {
	// We fetch the raw character blob (lexeme) and mark the corresponding token type and the parsed literal (if any)
	lexeme := s.source[s.start:s.current]
	// Creating and adding it into the token list
	token := NewToken(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, *token)
}

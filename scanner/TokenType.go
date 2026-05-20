package scanner

type TokenType int

const (
	LeftParen  = iota // (
	RightParen        // )
	LeftBrace         // {
	RightBrace        // }
	Comma             // ,
	Dot               // .
	Semicolon         // ;

	Minus // -
	Plus  // +
	Star  // *
	Slash // /

	Bang         // !
	BangEqual    // !=
	Equal        // =
	EqualEqual   // ==
	Greater      // >
	GreaterEqual // >=
	Less         // <
	LessEqual    // <=

	Identifier
	String
	Number

	And    // &&
	Class  // class
	Else   // else
	False  // false
	True   // true
	Fun    // fun (function)
	For    // for
	While  // while
	If     // if
	Nil    // nil (null)
	Or     // ||
	Print  // print (printf)
	Return // return (function)
	Super  // class inheritance related super
	This   // instance object "this"
	Var    // Variable declaration

	EOF // End Of File
)

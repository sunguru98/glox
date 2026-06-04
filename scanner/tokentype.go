package scanner

import "fmt"

type TokenType int

const (
	LeftParen  TokenType = iota // (
	RightParen                  // )
	LeftBrace                   // {
	RightBrace                  // }
	Comma                       // ,
	Dot                         // .
	Semicolon                   // ;

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
	Or     // ||
	If     // if
	Else   // else
	True   // true
	False  // false
	Class  // class
	Fun    // fun (function)
	For    // for
	While  // while
	Nil    // nil (null)
	Print  // print (printf)
	Return // return (function)
	Super  // class inheritance related super
	This   // instance object "this"
	Var    // Variable declaration

	EOF // End Of File
)

func (t TokenType) String() string {
	switch t {
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Semicolon:
		return "Semicolon"

	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"

	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case Equal:
		return "Equal"
	case EqualEqual:
		return "EqualEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Less:
		return "Less"
	case LessEqual:
		return "LessEqual"

	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Number:
		return "Number"

	case And:
		return "And"
	case Or:
		return "Or"
	case If:
		return "If"
	case Else:
		return "Else"
	case True:
		return "True"
	case False:
		return "False"
	case Class:
		return "Class"
	case Fun:
		return "Fun"
	case For:
		return "For"
	case While:
		return "While"
	case Nil:
		return "Nil"
	case Print:
		return "Print"
	case Return:
		return "Return"
	case Super:
		return "Super"
	case This:
		return "This"
	case Var:
		return "Var"

	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("Unidentified tokentype: %d", t)
	}
}

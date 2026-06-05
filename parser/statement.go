package parser

type Statement interface {
	stmnt()
}

type ExpressionSt struct {
	Expr Expression
}

func (*ExpressionSt) stmnt() {}

type PrintSt struct {
	Expr Expression
}

func (*PrintSt) stmnt() {}

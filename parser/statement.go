package parser

type Statement interface {
	stmnt()
}

type ExpressionSt struct {
	Expr Expression
}

func (*ExpressionSt) stmnt() {}

func CreateNewExpressionSt(expression Expression) *ExpressionSt {
	return &ExpressionSt{
		Expr: expression,
	}
}

type PrintSt struct {
	Expr Expression
}

func (*PrintSt) stmnt() {}

func CreateNewPrintSt(expression Expression) *PrintSt {
	return &PrintSt{
		Expr: expression,
	}
}

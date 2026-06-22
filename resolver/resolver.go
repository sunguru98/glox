package resolver

import (
	"github.com/sunguru98/glox/lib"
	p "github.com/sunguru98/glox/parser"
	s "github.com/sunguru98/glox/scanner"
)

// --------------------------------------------------------------------------------------------------------
type Resolver struct {
	Interpreter *p.Interpreter
	Scopes      []map[string]bool
}

func InitResolver(interpreter *p.Interpreter) *Resolver {
	return &Resolver{
		Interpreter: interpreter,
		Scopes:      make([]map[string]bool, 0),
	}
}

// --------------------------------------------------------------------------------------------------------

func (r *Resolver) resolve(blockSt p.BlockSt) {
	r.beginScope()
	r.resolveStatements(blockSt.Statements)
	r.endScope()
}

func (r *Resolver) resolveStatements(statements []p.Statement) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveStatement(statement p.Statement) {
	switch st := statement.(type) {
	case *p.VariableSt:
		r.declare(st.Name)
		if st.Initializer != nil {
			r.resolveExpression(st.Initializer)
		}
		r.define(st.Name)

	case *p.FunctionSt:
		r.declare(st.Name)
		r.define(st.Name)
		r.resolveFunction(st)

	case *p.ExpressionSt:
		r.resolveExpression(st.Expr)

	case *p.IfSt:
		r.resolveExpression(st.Condition)
		r.resolveStatement(st.ThenBranch)
		if st.ElseBranch != nil {
			r.resolveStatement(st.ElseBranch)
		}

	case *p.PrintSt:
		r.resolveExpression(st.Expr)

	case *p.ReturnSt:
		if st.Value != nil {
			r.resolveExpression(st.Value)
		}

	case *p.WhileSt:
		r.resolveExpression(st.Condition)
		r.resolveStatement(st.Body)
	}
}

func (r *Resolver) resolveExpression(expression p.Expression) {
	switch expr := expression.(type) {
	case *p.Variable:
		scopesLen := len(r.Scopes)
		scopeTop := r.Scopes[scopesLen-1]
		token := expr.Name

		if scopesLen != 0 && scopeTop[token.Lexeme] == false {
			lib.Report(token.LineNumber, " at '"+token.Lexeme+"'", "Can't read local variable in own initializer")
		}

		r.resolveLocal(expr, token)

	case *p.Call:
		r.resolveExpression(expr.Callee)
		for _, argument := range expr.Arguments {
			r.resolveExpression(argument)
		}

	case *p.Assignment:
		r.resolveExpression(expr.Value)
		r.resolveLocal(expr, expr.Name)

	case *p.Binary:
		r.resolveExpression(expr.Left)
		r.resolveExpression(expr.Right)

	case *p.Logical:
		r.resolveExpression(expr.Left)
		r.resolveExpression(expr.Right)

	case *p.Unary:
		r.resolveExpression(expr.Right)

	case *p.Grouping:
		r.resolveExpression(expr.Expression)

	case *p.Literal:
		return
	}
}

func (r *Resolver) resolveLocal(expr p.Expression, name s.Token) {
	for index := len(r.Scopes) - 1; index >= 0; index -= 1 {
		scope := r.Scopes[index]
		_, ok := scope[name.Lexeme]
		if ok {
			// TODO: Interpreter resolve
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *p.FunctionSt) {
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.resolveStatements(function.Body)
	r.endScope()
}

// --------------------------------------------------------------------------------------------------------

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	scopesLen := len(r.Scopes)
	if scopesLen >= 1 {
		r.Scopes = r.Scopes[:scopesLen]
	}
}

func (r *Resolver) declare(name s.Token) {
	scopesLen := len(r.Scopes)
	if scopesLen == 0 {
		return
	}

	scopeTop := r.Scopes[scopesLen-1]
	scopeTop[name.Lexeme] = false
}

func (r *Resolver) define(name s.Token) {
	scopesLen := len(r.Scopes)
	if scopesLen == 0 {
		return
	}

	scopeTop := r.Scopes[scopesLen-1]
	scopeTop[name.Lexeme] = true
}

// --------------------------------------------------------------------------------------------------------

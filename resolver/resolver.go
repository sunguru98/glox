package resolver

import (
	"fmt"

	p "github.com/sunguru98/glox/parser"
)

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

func (r *Resolver) beginScope() {
	r.Scopes = append(r.Scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	scopesLen := len(r.Scopes)
	if scopesLen >= 1 {
		r.Scopes = r.Scopes[:scopesLen]
	}
}

func (r *Resolver) resolve(blockSt p.BlockSt) {
	r.beginScope()
	for _, statement := range blockSt.Statements {
		r.resolveStatement(statement)
	}
	r.endScope()
}

func (r *Resolver) resolveStatement(statement p.Statement) {
	switch st := statement.(type) {
	default:
		fmt.Println(st)
	}
}

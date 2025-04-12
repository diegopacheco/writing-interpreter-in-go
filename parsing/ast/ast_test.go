package ast

import (
	"testing"

	"github.com/diegopacheco/writing-interpreter-in-go/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.TokenLiteral() != "let" {
		t.Errorf("program.TokenLiteral() wrong. got=%q", program.TokenLiteral())
	}

	letStmt := program.Statements[0].(*LetStatement)
	if letStmt.TokenLiteral() != "let" {
		t.Errorf("letStmt.TokenLiteral() wrong. got=%q", letStmt.TokenLiteral())
	}

	if letStmt.Name.TokenLiteral() != "myVar" {
		t.Errorf("letStmt.Name.TokenLiteral() wrong. got=%q", letStmt.Name.TokenLiteral())
	}

	if letStmt.Name.Value != "myVar" {
		t.Errorf("letStmt.Name.Value wrong. got=%q", letStmt.Name.Value)
	}
}

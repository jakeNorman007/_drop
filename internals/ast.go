package internals

import (
	"bytes"
	"strings"
	"_drop/internals"
)

type Node interface {
	Token_literal() string
	String()        string
}

type Statement interface {
	Node
	statement_node()
}

type Expression interface {
	Node
	expression_node()
}

type Program struct {
	Statements []Statement
}

func (program *Program) Token_literal() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].Token_literal()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var output bytes.Buffer

	for _, stmnt := range program.Statements {
		output.WriteString(stmnt.String())
	}
	return ""
}

type Identifier struct {
	Token Token // IDENTIFIER
	Value string
}

func (ident *Identifier) Token_literal() string { return ident.Token.Literal }
func (ident *Identifier) String() string { return ident.Value }

type Parameter struct {
	Name *Identifier
	Type *Identifier
}

type BlockStatement struct {
	Token Token 
	Statements []Statement
}

func (block *BlockStatement) expression_node() {}
func (block *BlockStatement) Token_literal() string { return block.Token.Literal}
func (block *BlockStatement) String() string {
	var output bytes.Buffer

	for _, stmnt := range block.Statements {
		output.WriteString(stmnt.String())
	}

	return output.String()
}

type DropFunction struct {
	Token			 Token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *Identifier
	Body       *BlockStatement
}

func (drop_func *DropFunction) statement_node() {}
func (drop_func *DropFunction) Token_literal() string { return drop_func.Token.Literal }
func (drop_func *DropFunction) String() string {
	var output bytes.Buffer

	parameters := []string{}
	for _, param := range drop_func.Parameters {
		parameters= append(parameters, param.Type.String())
	}

	output.WriteString(drop_func.Token_literal() + " ")
	output.WriteString(drop_func.Name.String())
	output.WriteString("(")
	output.WriteString(strings.Join(parameters, ", "))
	output.WriteString(") ")

	if drop_func.ReturnType != nil {
		output.WriteString(drop_func.ReturnType.String())
	}

	if drop_func.Body != nil {
		output.WriteString(drop_func.Body.String())
	}

	return output.String()
}

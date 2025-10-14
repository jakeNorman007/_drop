package internals

import (
	"bytes"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String()        string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var output bytes.Buffer

	for _, stmnt := range program.Statements {
		output.WriteString(stmnt.String())
	}
	return output.String()
}

type Identifier struct {
	Token Token
	Value string
}

func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string { return ident.Value }

type DropFunction struct {
	Token      Token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *Identifier
	Body       *BlockStatement
}

func (drop_func *DropFunction) statementNode() {}
func (drop_func *DropFunction) TokenLiteral() string { return drop_func.Token.Literal }
func (drop_func *DropFunction) String() string {
	var output bytes.Buffer

	parameters := []string{}
	for _, param := range drop_func.Parameters {
		parameters= append(parameters, param.Name.String() + " " + param.Type.String())
	}

	output.WriteString(drop_func.TokenLiteral() + " ")
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

type Parameter struct {
	Name *Identifier
	Type *Identifier
}

type BlockStatement struct {
	Token Token 
	Statements []Statement
}

func (block *BlockStatement) statementNode() {}
func (block *BlockStatement) TokenLiteral() string { return block.Token.Literal}
func (block *BlockStatement) String() string {
	var output bytes.Buffer

	for _, stmnt := range block.Statements {
		output.WriteString(stmnt.String())
	}

	return output.String()
}

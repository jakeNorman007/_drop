package internals

//comment from the unknown
import (
	"bytes"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String()       string
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

type Identifier struct {
	Token Token
	Value string
}

type DropFunction struct {
	Token      Token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *Identifier
}

type Parameter struct {
	Name *Identifier
	Type *Identifier
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

func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string { return ident.Value }

func (dropFunc *DropFunction) statementNode() {}
func (dropFunc *DropFunction) TokenLiteral() string { return dropFunc.Token.Literal }
func (dropFunc *DropFunction) String() string {
	var output bytes.Buffer

	parameters := []string{}

	for _, param := range dropFunc.Parameters {
		parameters= append(parameters, param.Name.String() + " " + param.Type.String())
	}

	output.WriteString(dropFunc.TokenLiteral() + " ")
	output.WriteString(dropFunc.Name.String())
	output.WriteString("(")
	output.WriteString(strings.Join(parameters, ", "))
	output.WriteString(") ")

	if dropFunc.ReturnType != nil {
		output.WriteString(dropFunc.ReturnType.String())
	}

	return output.String()
}

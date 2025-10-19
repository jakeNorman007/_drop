package internals

import "fmt"

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

// NEED TO ORGANIZE THESE HELPER FUNCTIONS= =============================
func (p *Parser) Errors() []string {
    return p.errors
}
	
func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) currentTokenIs(tokn TokenType) bool {
	return parser.currentToken.Type == tokn
}

func (parser *Parser) peekTokenIs(tokn TokenType) bool {
	return parser.peekToken.Type == tokn
}

func (parser *Parser) expectPeek(tokn TokenType) bool {
	if parser.peekTokenIs(tokn) {
		parser.nextToken()

		return true
	} else {
		parser.peekError(tokn)

		return false
	}
}

func (parser *Parser) peekError(tokn TokenType) {
	message := fmt.Sprintf("expected the next token to be: %s, got %s instead", tokn, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) ParseProgram() *Program {
	program := &Program {}
	program.Statements = []Statement{}

	for parser.currentToken.Type != EOF {
		stmnt := parser.parseStatement()

		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}

		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() Statement {
	switch parser.currentToken.Type {
	case DROP_START:
		return parser.ParseDropFunction()
	default:
		return nil
	}
}
//= =====================================================================

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser { lexer: lexer }

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) ParseDropFunction() *DropFunction {
	function := &DropFunction { Token: parser.currentToken }

	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}

	function.Name = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }

	if !parser.expectPeek(LEFT_PARENTHESIS) {
		return nil
	}

	function.Parameters = parser.parseDropFunctionParameters()

	if parser.peekTokenIs(IDENTIFIER) {
		parser.nextToken()
		function.ReturnType = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }
	}

	if !parser.expectPeek(LEFT_CURLY_BRACE) {
		return nil
	}

	return function
}

func (parser *Parser) parseDropFunctionParameters() []*Parameter {
	parameters := []*Parameter {}

	if parser.peekTokenIs(RIGHT_PARENTHESIS) {
		parser.nextToken()

		return parameters
	}

	parser.nextToken()

	parameter := &Parameter { Name: &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }}

	if !parser.expectPeek(IDENTIFIER) {
		return nil
	}

	parameter.Type = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }
	parameters = append(parameters, parameter)

	for parser.peekTokenIs(COMMA) {
		parser.nextToken()
		parser.nextToken()

		parameter := &Parameter { 
			Name: &Identifier { 
				Token: parser.currentToken, 
				Value: parser.currentToken.Literal,
			},
		}

		if !parser.expectPeek(IDENTIFIER) {
			return nil
		}

		parameter.Type = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }
		parameters = append(parameters, parameter)
	}

	if !parser.expectPeek(RIGHT_PARENTHESIS) {
		return nil
	}

	return parameters
}

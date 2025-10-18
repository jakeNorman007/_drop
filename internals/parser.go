package internals

import "fmt"

type Parser struct {
	lexer         *Lexer
	currentToken Token
	peekToken    Token
	errors        []string
}

// NEED TO ORGANIZE THESE HELPER FUNCTIONS= =============================
func (p *Parser) Errors() []string {
    return p.errors
}
	
func (parser *Parser) nextToken() {
	fmt.Printf("[nextToken] current=%s, peek=%s\n", parser.currentToken.Type, parser.peekToken.Type)
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
	fmt.Printf("→ advanced to current=%s, peek=%s\n", parser.currentToken.Type, parser.peekToken.Type)
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
	fmt.Println("[ParseDropFunction] Starting to parse drop function...")

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

	function.Body = parser.parseBlockStatement()

	return function
}

func (parser *Parser) parseDropFunctionParameters() []*Parameter {
	fmt.Println("[parseDropFunctionParameters] Parsing parameters...")
	parameters := []*Parameter {}

	if parser.peekTokenIs(RIGHT_PARENTHESIS) {
		fmt.Println("[parseDropFunctionParameters] No parameters found.")
		parser.nextToken()

		return parameters
	}

	parser.nextToken()
	fmt.Printf("[parseDropFunctionParameters] Parsing first parameter name: %s\n", parser.currentToken.Literal)

	parameter := &Parameter { 
		Name: &Identifier { 
			Token: parser.currentToken, 
			Value: parser.currentToken.Literal,
		},
	}

	if !parser.expectPeek(IDENTIFIER) {
		fmt.Println("[parseDropFunctionParameters] Expected parameter type after name but got:", parser.peekToken.Type)
		return nil
	}

	parameter.Type = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }
	parameters = append(parameters, parameter)

	fmt.Printf("[parseDropFunctionParameters] → First parameter: %s %s\n", parameter.Name.Value, parameter.Type.Value)

	for parser.peekTokenIs(COMMA) {
		parser.nextToken()
		parser.nextToken()

		fmt.Printf("[parseDropFunctionParameters] Found comma, parsing next parameter (%d)...\n", len(parameters)+1)
		fmt.Printf("  Current token: %s (%s)\n", parser.currentToken.Type, parser.currentToken.Literal)

		parameter := &Parameter { 
			Name: &Identifier { 
				Token: parser.currentToken, 
				Value: parser.currentToken.Literal,
			},
		}

		if !parser.expectPeek(IDENTIFIER) {
			fmt.Println("[parseDropFunctionParameters] Expected type after comma-separated parameter name.")
			return nil
		}

		parameter.Type = &Identifier { Token: parser.currentToken, Value: parser.currentToken.Literal }
		fmt.Printf("[parseDropFunctionParameters] → Parsed parameter #%d: %s %s\n", len(parameters)+1, parameter.Name.Value, parameter.Type.Value)
		parameters = append(parameters, parameter)
	}

	if !parser.expectPeek(RIGHT_PARENTHESIS) {
		fmt.Println("[parseDropFunctionParameters] Missing closing parenthesis after parameters.")
		return nil
	}

	fmt.Printf("[parseDropFunctionParameters] Done. Parsed %d parameter(s).\n", len(parameters))
	return parameters
}

func (parser *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement { Token: parser.currentToken, Statements: []Statement{} }
	block.Statements = []Statement{}

	parser.nextToken()

	for !parser.currentTokenIs(RIGHT_CURLY_BRACE) && !parser.currentTokenIs(EOF) {
		stmnt := parser.parseStatement()

		if stmnt != nil {
			block.Statements = append(block.Statements, stmnt)
		}

		parser.nextToken()
	}

	return block
}

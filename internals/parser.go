package internals

import (
	"fmt"
	"_drop/internals"
)

type Parser struct {
	lexer         *Lexer
	current_token Token
	peek_token    Token
	errors        []string
}

// NEED TO ORGANIZE THESE HELPER FUNCTIONS= =============================
func (parser *Parser) next_token() {
	parser.current_token = parser.peek_token
	parser.peek_token = parser.lexer.Next_token()
}

func (parser *Parser) current_token_is(tokn TokenType) bool {
	return parser.current_token.Type == tokn
}

func (parser *Parser) peek_token_is(tokn TokenType) bool {
	return parser.peek_token.Type == tokn
}

func (parser *Parser) expect_peek(tokn TokenType) bool {
	if parser.peek_token_is(tokn) {
		parser.next_token()

		return true
	} else {
		parser.peek_error(tokn)

		return false
	}
}

func (parser *Parser) peek_error(tokn TokenType) {
	message := fmt.Sprintf("expected the next token to be: %s, got %s instead", tokn, parser.peek_token.Type)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) Parse_program() *Program {
	program := &Program {}
	program.Statements = []Statement{}

	for parser.current_token.Type != EOF {
		stmnt := parser.parse_statement()

		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}

		parser.next_token()
	}

	return program
}

func (parser *Parser) parse_statement() Statement {
	switch parser.current_token.Type {
	case DROP_START:
		return parser.Parse_drop_function()
	default:
		return nil
	}
}
//= =====================================================================

func New_lexer(lexer *Lexer) *Parser {
	parser := &Parser {
		lexer: lexer,
	}

	parser.next_token()
	parser.next_token()

	return parser
}

func (parser *Parser) Parse_drop_function() *DropFunction { // serves as entry point into _drop template
	function := &DropFunction {
		Token: parser.current_token,
	}

	parser.next_token()
	function.Name = &Identifier {
		Token: parser.current_token,
		Value: parser.current_token.Literal,
	}

	parser.next_token()
	function.Parameters = parser.parse_drop_function_parameters()

	if parser.peek_token_is(IDENTIFIER) && !parser.peek_token_is(LEFT_CURLY_BRACE) {

		parser.next_token()
		function.ReturnType = &Identifier {
			Token: parser.current_token,
			Value: parser.current_token.Literal,
		}
	}

	if !parser.expect_peek(LEFT_CURLY_BRACE) {
		return nil
	}

	function.Body = parser.parse_block_statement()

	return function
}

func (parser *Parser) parse_drop_function_parameters() []*Parameter {
	parameters := []*Parameter {}

	if parser.peek_token_is(RIGHT_PARENTHESES) {
		parser.next_token()

		return parameters
	}

	parser.next_token()
	for !parser.current_token_is(RIGHT_PARENTHESES) {
		name := &Identifier {
			Token: parser.current_token,
			Value: parser.current_token.Literal,
		}

		parser.next_token()
		_type := &Identifier {
			Token: parser.current_token,
			Value: parser.current_token.Literal,
		}
		
		parameters = append(parameters, &Parameter {
			Name: name,
			Type: _type,
		})

		if parser.peek_token_is(COMMA) {
			parser.next_token()
			parser.next_token()
		} else {
			break
		}
	}

	parser.expect_peek(RIGHT_PARENTHESES)

	return parameters
}

func (parser *Parser) parse_block_statement() *BlockStatement {
	block := &BlockStatement {
		Token: parser.current_token,
	}

	for !parser.current_token_is(RIGHT_CURLY_BRACE) && parser.current_token_is(EOF){
		stmnt := parser.parse_statement()
		if stmnt != nil {
			block.Statements = append(block.Statements, stmnt)
		}

		parser.next_token()
	}

	return block
}

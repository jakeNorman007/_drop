package internals

import (
	"strings"
)

type Lexer struct {
	input            string
	current_position int
	read_position    int
	char             byte
}

func New_Lexer(input string) *Lexer {
	lex := &Lexer {
		input: input,
	}

	lex.read_char()

	return lex
}

func (lex *Lexer) Next_token() Token {
	var tokn Token

	lex.consume_whitespace()

	switch lex.char {
	case '\x28':
		tokn = new_token(LEFT_PARENTHESES, lex.char)
	case '\x29':
		tokn = new_token(RIGHT_PARENTHESES, lex.char)
	case '\x7B':
		if lex.peek_char() == '{' {
			char := lex.char
			lex.read_char()
			tokn = Token{ Type: DOUBLE_LEFT_CURLY_BRACE, Literal: string(char) + "{" }
		} else {
			tokn = new_token(LEFT_CURLY_BRACE, lex.char)
		}
	case '\x7D':
		if lex.peek_char() == '}' {
			char := lex.char
			lex.read_char()
			tokn = Token{ Type: DOUBLE_RIGHT_CURLY_BRACE, Literal: string(char) + "}" }
		} else {
			tokn = new_token(RIGHT_CURLY_BRACE, lex.char)
		}
	case '\x3C':
			tokn = new_token(LESS_THAN, lex.char)
	case '\x3E':
		tokn = new_token(GREATER_THAN, lex.char)
	case '\x2F':
		tokn = new_token(SLASH, lex.char)
	case '\x3D':
		tokn = new_token(EQUALS, lex.char)
	case '\x21':
		tokn = new_token(EXCLAMATION, lex.char)
	case '\x22':
		tokn.Type = STRING
		tokn.Literal = lex.read_string()
	case '\x27':
		tokn = new_token(SINGLE_QUOTE, lex.char)
	case '\x2C':
		tokn = new_token(COMMA, lex.char)
	case '\x5F':
		lex.read_char()
		drop_start := lex.read_identifier()

		if drop_start == "drop" {
			tokn = Token{ Type: DROP_START, Literal: "_" + string(drop_start) }
		} else {
			tokn = new_token(UNDERSCORE, lex.char)
		}
	case '0':
		tokn.Literal = ""
		tokn.Type = EOF
	default:
		if is_letter(lex.char) {
			tokn.Literal = lex.read_identifier()
			tokn = Token{ Type: IDENTIFIER, Literal: tokn.Literal }
			return tokn
		} else if is_closing_tag(lex.char) {
			tokn.Literal = lex.read_identifier()
			tokn = Token{ Type: IDENTIFIER, Literal: tokn.Literal }
			return tokn
		} else {
			tokn = new_token(ILLEGAL, lex.char)
		}

		if lex.char != '<' && lex.char != '{' {
			literal := lex.read_html_text()

			tokn = Token{ Type: HTML_TEXT, Literal: strings.TrimSpace(literal) }
			return tokn
		}
	}

	lex.read_char()

	return tokn
}

// helper functions

func new_token(tokenType TokenType, char byte) Token {
	return Token { Type: tokenType, Literal: string(char)  }
}

func (lex *Lexer) read_char() {
	if lex.read_position >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.read_position]
	}

	lex.current_position = lex.read_position
	lex.read_position += 1
}

func (lex *Lexer) peek_char() byte { // not used yet, but am keeping just in case
	if lex.read_position >= len(lex.input) {
		return 0
	} else {
		return lex.input[lex.read_position]
	}
}

func (lex *Lexer) read_identifier() string {
	position := lex.current_position

	for is_letter(lex.char) {
		lex.read_char()
	}

	return lex.input[position:lex.current_position]
}

func (lex *Lexer) read_string() string {
	position := lex.current_position + 1

	for {
		lex.read_char()

		if lex.char == '"' || lex.char == 0 {
			break
		}
	}

	return lex.input[position:lex.current_position]
}

func (lex *Lexer) read_html_text() string {
	position := lex.current_position

	for lex.char != LESS_THAN && lex.char !=  LEFT_CURLY_BRACE && lex.char != 0 {
		lex.read_char()
	}

	return lex.input[position:lex.current_position]
}

func is_letter(char byte) bool {
	// here we can treat UNDERSCORE as a letter for function casing
	return LOWER_A <= char && char <= LOWER_Z || UPPER_A <= char && UPPER_Z <= char || char == UNDERSCORE
}

func is_closing_tag(char byte) bool {
	return char == SLASH
}

func (lex *Lexer) consume_whitespace() {
	for lex.char == SPACE || lex.char == TAB || lex.char == NEW_LINE || lex.char == CARRIAGE_RETURN {
		lex.read_char()
	}
}

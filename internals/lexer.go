package internals

type Lexer struct {
	input            string
	currentPosition int
	readPosition    int
	char             byte
}

func NewLexer(input string) *Lexer {
	lex := &Lexer {
		input: input,
	}

	lex.readChar()

	return lex
}

func (lex *Lexer) NextToken() Token {
	var tokn Token

	lex.consumeWhitespace()

	switch lex.char {
	case '\x28':
		tokn = newToken(LEFT_PARENTHESIS, lex.char)
	case '\x29':
		tokn = newToken(RIGHT_PARENTHESIS, lex.char)
	case '\x7B':
		if lex.peekChar() == '{' {
			char := lex.char
			lex.readChar()
			tokn = Token{ Type: DOUBLE_LEFT_CURLY_BRACE, Literal: string(char) + "{" }
		} else {
			tokn = newToken(LEFT_CURLY_BRACE, lex.char)
		}
	case '\x7D':
		if lex.peekChar() == '}' {
			char := lex.char
			lex.readChar()
			tokn = Token{ Type: DOUBLE_RIGHT_CURLY_BRACE, Literal: string(char) + "}" }
		} else {
			tokn = newToken(RIGHT_CURLY_BRACE, lex.char)
		}
	case '\x3C':
			tokn = newToken(LESS_THAN, lex.char)
	case '\x3E':
		tokn = newToken(GREATER_THAN, lex.char)
	case '\x2F':
		tokn = newToken(SLASH, lex.char)
	case '\x3D':
		tokn = newToken(EQUALS, lex.char)
	case '\x21':
		tokn = newToken(EXCLAMATION, lex.char)
	case '\x22':
		tokn.Type = STRING
		tokn.Literal = lex.readString()
	case '\x27':
		tokn = newToken(SINGLE_QUOTE, lex.char)
	case '\x2C':
		tokn = newToken(COMMA, lex.char)
	case '\x5F':
		drop_start := lex.readIdentifier()

		var token_type TokenType

		if drop_start == "_drop" {
			token_type = DROP_START
		} else {
			token_type = IDENTIFIER
		}

		tokn = Token { Type: token_type, Literal: drop_start }
	default:
		if isLetter(lex.char) || lex.char == '_' {
			tokn.Literal = lex.readIdentifier()
			tokn = Token{ Type: IDENTIFIER, Literal: tokn.Literal }
			return tokn
		} else if isClosingTag(lex.char) {
			tokn.Literal = lex.readIdentifier()
			tokn = Token{ Type: IDENTIFIER, Literal: tokn.Literal }
			return tokn
		} else {
			tokn = newToken(ILLEGAL, lex.char)
		}

		if lex.char != '<' && lex.char != '{' {
				literal := lex.readHtmlText()

			if literal != "" {
				tokn = Token{ Type: HTML_TEXT, Literal: literal }
			} else if lex.char == 0 {
				lex.readChar()
				tokn = Token { Type: EOF, Literal: "" }
			} else {
				char := lex.char
				lex.readChar()
				tokn = newToken(ILLEGAL, char)
			}
			return tokn
		}
	}

	lex.readChar()

	return tokn
}

// helper functions

func newToken(tokenType TokenType, char byte) Token {
	return Token { Type: tokenType, Literal: string(char) }
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPosition]
	}

	lex.currentPosition = lex.readPosition
	lex.readPosition += 1
}

func (lex *Lexer) peekChar() byte {
	if lex.readPosition >= len(lex.input) {
		return 0
	} else {
		return lex.input[lex.readPosition]
	}
}

func (lex *Lexer) readIdentifier() string {
	position := lex.currentPosition

	for isLetter(lex.char) || isDigit(lex.char){
		lex.readChar()
	}

	return lex.input[position:lex.currentPosition]
}

func (lex *Lexer) readString() string {
	position := lex.currentPosition + 1

	for {
		lex.readChar()

		if lex.char == '"' || lex.char == 0 {
			break
		}
	}

	return lex.input[position:lex.currentPosition]
}

func (lex *Lexer) readHtmlText() string {
	 position := lex.currentPosition
    for lex.char != '<' && lex.char != 0 {
        lex.readChar()
        if lex.readPosition > len(lex.input) {
            break
        }
    }
    if position >= len(lex.input) {
        return ""
    }
    if lex.currentPosition > len(lex.input) {
        lex.currentPosition = len(lex.input)
    }
    return lex.input[position:lex.currentPosition]
}

func isLetter(char byte) bool {
	return LOWER_A <= char && char <= LOWER_Z || UPPER_A <= char && UPPER_Z <= char || char == UNDERSCORE
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isClosingTag(char byte) bool {
	return char == SLASH
}

func (lex *Lexer) consumeWhitespace() {
	for lex.char == SPACE || lex.char == TAB || lex.char == NEW_LINE || lex.char == CARRIAGE_RETURN {
		lex.readChar()
	}
}

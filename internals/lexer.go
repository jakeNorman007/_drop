package internals

type Lexer struct {
	input           string
	currentPosition int
	readPosition    int
	char            byte
}

func NewLexer(input string) *Lexer {
	lex := &Lexer { input: input }

	lex.readChar()

	return lex
}

func newToken(tokenType TokenType, char byte) Token {
	return Token { Type: tokenType, Literal: string(char) }
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
		if lex.peekChar() == '\x7B' {
			char := lex.char
			lex.readChar()
			tokn = Token { Type: DOUBLE_LEFT_CURLY_BRACE, Literal: string(char) + "\x7B" }
		} else {
			tokn = newToken(LEFT_CURLY_BRACE, lex.char)
		}
	case '\x7D':
		if lex.peekChar() == '\x7D' {
			char := lex.char
			lex.readChar()
			tokn = Token { Type: DOUBLE_RIGHT_CURLY_BRACE, Literal: string(char) + "\x7D" }
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
		start, end := lex.readStringIdx()
		tokn.Type = STRING
		tokn.Literal = lex.input[start:end]
	case '\x27':
		tokn = newToken(SINGLE_QUOTE, lex.char)
	case '\x2C':
		tokn = newToken(COMMA, lex.char)
	case '\x5F':
		start, end := lex.readIdentifierIdx()
		literal := lex.input[start:end]

		if literal == "_drop" {
			tokn = Token {Type: DROP_START, Literal: literal } 
		} else {
			tokn = Token { Type: IDENTIFIER, Literal: literal } 
		}
	default:
		if isLetter(lex.char) || lex.char == '\x5F' {
			start, end := lex.readIdentifierIdx()
			tokn = Token { Type: IDENTIFIER, Literal: lex.input[start:end] }
			return tokn
		} else if isClosingTag(lex.char) {
			start, end := lex.readIdentifierIdx()
			tokn = Token { Type: IDENTIFIER, Literal: lex.input[start:end] }
			return tokn
		} else {
			tokn = newToken(ILLEGAL, lex.char)
		}

		if lex.char != '\x3C' && lex.char != '\x7B' {
			start, end := lex.readHtmlTextIdx()

			if start != end {
				tokn = Token { Type: HTML_TEXT, Literal: lex.input[start:end] }
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

func (lex *Lexer) readIdentifierIdx() (start, end int) {
	start = lex.currentPosition

	for isLetter(lex.char) || isDigit(lex.char){
		lex.readChar()
	}

	end = lex.currentPosition

	return
}

func (lex *Lexer) readStringIdx() (start, end int) {
	start = lex.currentPosition + 1

	for {
		lex.readChar()

		if lex.char == '\x22' || lex.char == 0 {
			break
		}
	}

	end = lex.currentPosition

	if lex.char == '\x22' {
		end = lex.currentPosition
	}

	if end > len(lex.input) {
		end = len(lex.input)
	}

	if end < start {
		end = start
	}

	return
}

func (lex *Lexer) readHtmlTextIdx() (start, end int) {
	start = lex.currentPosition

	for lex.char != '\x3C' && lex.char != 0 {
		lex.readChar()

		if lex.readPosition > len(lex.input) {
			break
		}
	}

	if lex.currentPosition > len(lex.input) {
		lex.currentPosition = len(lex.input)
	}

	end = lex.currentPosition

	return
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

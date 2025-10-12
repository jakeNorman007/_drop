package internals_test

import (
	"testing"
	"_drop/internals"
)

func TestNextToken(t *testing.T) {
	input := `_drop my_function() {
							<div style="background: red;">
								<p class="my_class">
									<p id="my_id"> {{ }} </p>
								</p>
							</div>
						}`

	tests := []struct {
		expectedType			internals.TokenType
		expectedLiteral 	string
	}{
		{ internals.DROP_START, "_drop" },
		{ internals.IDENTIFIER, "my_function" },
		{ internals.LEFT_PARENTHESES, "(" },
		{ internals.RIGHT_PARENTHESES, ")" },
		{ internals.LEFT_CURLY_BRACE, "{" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "div" },
		{ internals.IDENTIFIER, "style" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "background: red;" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "p" },
		{ internals.IDENTIFIER, "class" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "my_class" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "p" },
		{ internals.IDENTIFIER, "id" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "my_id" },
		{ internals.GREATER_THAN, ">" },
		{ internals.DOUBLE_LEFT_CURLY_BRACE, "{{" },
		{ internals.DOUBLE_RIGHT_CURLY_BRACE, "}}" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "p" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "p" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "div" },
		{ internals.GREATER_THAN, ">" },
		{ internals.RIGHT_CURLY_BRACE, "}" },
	}

	lex := internals.New_Lexer(input)

	for i, tt := range tests {
		tokn := lex.Next_token()

		if tokn.Type != tt.expectedType {
			t.Fatalf("tests[%d] -> token type incorrect. expected=(%d, %q) got=(%d, %q)", i, tt.expectedType, tt.expectedType, tokn.Type, tokn.Type)
		}

		if tokn.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] -> literal incorrect. expected=%s got=%s", i, tt.expectedLiteral, tokn.Literal)

			/* KEEPING THIS ERROR FORMAT JIC, HELPED DETERMINE NON EATEN WHITESPACE */
			//t.Errorf("tests[%d] -> literal incorrect.\nexpected=%q\n     got=%q\nbytes expected=%v\nbytes got=%v", 
    		//i, tt.expectedLiteral, tokn.Literal, []byte(tt.expectedLiteral), []byte(tokn.Literal))
		}
	}
}

func TestParseHtmlSkeleton(t *testing.T) {
	// TODO: fix this so both DOCTYPE and doctype can be scanned by lexer, works for now w/ lowercase doctype
	input := `<!doctype html>
						<html lang="en">
							<head>
								<meta charset="UTF-8">
								<meta name="viewport" content="width=device-width, initial-scale=1.0">
								<title>Test Title</title>
							</head>
							<body>
							</body>
						</html>
						`

	tests := []struct {
		expectedType			internals.TokenType
		expectedLiteral 	string
	}{
		{ internals.LESS_THAN, "<" },
		{ internals.EXCLAMATION, "!" },
		{ internals.IDENTIFIER, "doctype" },
		{ internals.IDENTIFIER, "html" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "html" },
		{ internals.IDENTIFIER, "lang" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "en" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "head" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "meta" },
		{ internals.IDENTIFIER, "charset" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "UTF-8" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "meta" },
		{ internals.IDENTIFIER, "name" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "viewport" },
		{ internals.IDENTIFIER, "content" },
		{ internals.EQUALS, "=" },
		{ internals.STRING, "width=device-width, initial-scale=1.0" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "title" },
		{ internals.GREATER_THAN, ">" },
		{ internals.HTML_TEXT, "Test Title" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "title" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "head" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.IDENTIFIER, "body" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "body" },
		{ internals.GREATER_THAN, ">" },
		{ internals.LESS_THAN, "<" },
		{ internals.SLASH, "/" },
		{ internals.IDENTIFIER, "html" },
		{ internals.GREATER_THAN, ">" },
	}

	lex := internals.New_Lexer(input)

	for i, tt := range tests {
		tokn := lex.Next_token()

		if tokn.Type != tt.expectedType {
			t.Fatalf("tests[%d] -> token type incorrect. expected=(%d, %q) got=(%d, %q)", i, tt.expectedType, tt.expectedType, tokn.Type, tokn.Type)
		}

		if tokn.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] -> literal incorrect. expected=%s got=%s", i, tt.expectedLiteral, tokn.Literal)
		}
	}
}

// TODO:
	// create repo [ ]
	// parse strings for inline css styles [x]
	// write test to parse HTML skeleton [x]
	// clean up/refactor etc. [ ]
	// figure out the whitespsce thing with interpolation [ ]


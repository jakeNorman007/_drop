package internals

import "fmt"

type TokenType int

type Token struct {
	Type			TokenType
	Literal		string
}

const (
	ILLEGAL	TokenType 	= iota // 
	EOF											   // end of file

	IDENTIFIER

	DROP_START						 		 // _drop

	STRING										 // any content inside a string

	DOUBLE_LEFT_CURLY_BRACE    // interpolation begin {{
	DOUBLE_RIGHT_CURLY_BRACE	 // interpolation end }}

	LEFT_PARENTHESES 		= 0x28 // (
	RIGHT_PARENTHESES 	= 0x29 // )
	LEFT_CURLY_BRACE 		= 0x7B // {
	RIGHT_CURLY_BRACE 	= 0x7D // }
	LESS_THAN           = 0x3C // <
	GREATER_THAN			  = 0x3E // >
	SLASH								= 0x2F // /
	EQUALS							= 0x3D // =
	SINGLE_QUOTE	      = 0x27 // '
	COMMA								= 0x2C // ,
	EXCLAMATION					= 0x21 // !
	LOWER_A						  = 0x61 // a	
	UPPER_A							= 0X41 // A
	LOWER_Z             = 0x7A // z
	UPPER_Z             = 0X5A // Z
	UNDERSCORE					= 0x5F // _
	TAB									= 0x09 // \t
	SPACE								= 0x20 // space
	NEW_LINE            = 0x0A // \n
	CARRIAGE_RETURN			= 0x0D // \r

	HTML_TEXT						= iota + 9999
)

var keywords = map[string]TokenType {
	"_drop": DROP_START,
}

func Lookup_identifier(identifier string) TokenType {
    if tok, ok := keywords[identifier]; ok {
        return tok
    }

    return IDENTIFIER
}

func(tokn TokenType) String() string {
	switch tokn {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case DROP_START:
		return "DROP_START"
	case LEFT_PARENTHESES:
		return "("
	case RIGHT_PARENTHESES:
		return ")"
	case LEFT_CURLY_BRACE:
		return "{"
	case RIGHT_CURLY_BRACE:
		return "}"
	case LESS_THAN:
		return "<"
	case GREATER_THAN:
		return ">"
	case SLASH:
		return "/" 
	case STRING:
		return "STRING"
	case EQUALS:
		return "="
	case SINGLE_QUOTE:
		return "'" 
	case COMMA:
		return ","
	case EXCLAMATION:
		return "!"
	case LOWER_A:
		return "a"
	case UPPER_A:
		return "A"
	case LOWER_Z:
		return "z"
	case UPPER_Z:
		return "Z"
	case UNDERSCORE:
		return "_"
	case TAB:
		return "\t"
	case SPACE:
		return " "
	case NEW_LINE:
		return "\n"
	case CARRIAGE_RETURN:
		return "\r"
	case HTML_TEXT:
		return "HTML_TEXT"
	case DOUBLE_LEFT_CURLY_BRACE:
		return "DOUBLE_LEFT_CURLY_BRACE"
	case DOUBLE_RIGHT_CURLY_BRACE:
		return "DOUBLE_RIGHT_CURLY_BRACE"
	}

	return fmt.Sprintf("UNKNOWN(%d)", int(tokn))
}

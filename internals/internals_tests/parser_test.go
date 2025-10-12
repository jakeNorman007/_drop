package internals_test

import (
	"testing"
	"_drop/internals"
)

func TestParseDropFunction(t *testing.T) {

	//
	// The issue seems to be when I add arguments
	//

	input := `_drop my_function() string {}`

	lexer := internals.New_Lexer(input)
	parser := internals.New_Parser(lexer)
	program := parser.Parse_program()
	check_parser_errors(t, parser)

	if program == nil {
		t.Fatalf("Parse_program() returned nil")
	}

	if len(program.Statements) != 1 {
        t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*internals.DropFunction)

    if !ok {
        t.Fatalf("Expected *internals.DropFunction, got %T", program.Statements[0])
    }

    if stmt.Name.Value != "my_function" {
        t.Fatalf("Expected function name 'my_function', got %q", stmt.Name.Value)
    }

		if stmt.ReturnType.Value != "string" {
        t.Fatalf("Expected return type 'string', got %q", stmt.ReturnType.Value)
		}

    output := stmt.String()
    expected := "_drop my_function() string"

    if output != expected {
        t.Errorf("String() output mismatch.\nexpected=%q\ngot=%q", expected, output)
    }

}

func check_parser_errors(t *testing.T, parser *internals.Parser) {
    errors := parser.Errors()

    if len(errors) == 0 {
        return
    }

    t.Errorf("parser encountered %d errors.", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }

    t.FailNow()
}

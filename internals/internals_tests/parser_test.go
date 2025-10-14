package internals_test

import (
	"testing"
	"_drop/internals"
)

func TestParseDropFunction(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{ input: "_drop my_function(arg string, arg2 string) {}", expectedParams: []string{"arg", "string", "arg2", "string"} },
	}

	for _, tt := range tests {
		lexer := internals.NewLexer(tt.input)
		parser := internals.NewParser(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("No program statements. got=%d", len(program.Statements))
		}

		stmnt := program.Statements[0]

		dropFunction, ok := stmnt.(*internals.DropFunction)
		if !ok {
			t.Fatalf("Statement is not _drop function. got=%T", stmnt)
		}

	}
}

func checkParserErrors(t *testing.T, parser *internals.Parser) {
    errors := parser.Errors()

    if len(errors) == 0 {
        return
    }

    t.Errorf("Parser encountered %d errors.", len(errors))
    for _, msg := range errors {
        t.Errorf("Parser error: %q", msg)
    }

    t.FailNow()
}

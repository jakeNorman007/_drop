package internals_test

import (
	"testing"
	"_drop/internals"
)

func TestParseDropFunction(t *testing.T) {
	tests := []struct {
        input          string
        expectedName   string
        expectedParams []string
        expectedTypes  []string
        expectedReturn string
    }{
				// TODO: GOT TO FIGURE OUT HOW TO PARSE MULTIPLE ARGUMENTS
				// ONE ARGUMENT AND RETURN TYPES WORK FINE
        {
            input:          "_drop my_function() {}",
            expectedName:   "my_function",
            expectedParams: []string{},
            expectedTypes:  []string{},
            expectedReturn: "",
        },
        {
            input:          "_drop fn() int {}",
            expectedName:   "fn",
            expectedParams: []string{},
            expectedTypes:  []string{},
            expectedReturn: "int",
        },
        {
            input:          "_drop example(name string, email string) string {}",
            expectedName:   "example",
            expectedParams: []string{"name", "email"},
            expectedTypes:  []string{"string", "string"},
            expectedReturn: "string",
        },
    }

    for _, tt := range tests {
        l := internals.NewLexer(tt.input)
        p := internals.NewParser(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements wrong. got=%d", len(program.Statements))
        }

        stmnt := program.Statements[0]
        dropFn, ok := stmnt.(*internals.DropFunction)
        if !ok {
            t.Fatalf("stmnt not *internals.DropFunction. got=%T", stmnt)
        }

        if dropFn.Name.Value != tt.expectedName {
            t.Errorf("expected name=%q, got=%q", tt.expectedName, dropFn.Name.Value)
        }

        if len(dropFn.Parameters) != len(tt.expectedParams) {
            t.Fatalf("parameter length mismatch. expected=%d, got=%d", len(tt.expectedParams), len(dropFn.Parameters))
        }

        for i, param := range dropFn.Parameters {
            if param.Name.Value != tt.expectedParams[i] {
                t.Errorf("param[%d] name wrong. expected=%q, got=%q", i, tt.expectedParams[i], param.Name.Value)
            }

            if dropFn.Parameters[i].Type.Value != tt.expectedTypes[i] {
                t.Errorf("param[%d] type wrong. expected=%q, got=%q", i, tt.expectedTypes[i], dropFn.Parameters[i].Type.Value)
            }
        }

        if dropFn.ReturnType != nil && dropFn.ReturnType.Value != tt.expectedReturn {
            t.Errorf("expected return=%q, got=%q", tt.expectedReturn, dropFn.ReturnType.Value)
        }

        if dropFn.Body == nil {
            t.Errorf("expected body to be parsed, got nil")
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

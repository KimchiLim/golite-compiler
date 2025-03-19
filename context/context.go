package context

import ("fmt")

type CompilerPhase int

const (LEXER CompilerPhase = iota ; PARSER ; SEMANTIC)

type CompilerError struct {
	Line   int           // The line where the error occurred
	Column int           // The column where the error occurred
	Msg    string        // A user friendly message to display
	Phase  CompilerPhase // The phase of the compiler were the error was generated
}

func (err *CompilerError) String() string {
	if err.Phase == LEXER {
		return fmt.Sprintf("lexer error (%d:%d) : %s", err.Line, err.Column, err.Msg)
	} else if err.Phase == PARSER {
		return fmt.Sprintf("syntax error (%d:%d) : %s", err.Line, err.Column, err.Msg)
	} else if err.Phase == SEMANTIC {
		return fmt.Sprintf("semantic error (%d:%d) : %s", err.Line, err.Column, err.Msg)
	}
	panic("Invalid phase found")
}

func HasErrors(errs []*CompilerError) bool {
	if len(errs) > 0 {
		for _, err := range errs { fmt.Println(err) }
		return true
	}
	return false
}

package lexer

import ("golite/context" ; "fmt" ; "github.com/antlr/antlr4/runtime/Go/antlr/v4")

type Lexer interface {
	GetTokenStream() *antlr.CommonTokenStream
	GetErrors() []*context.CompilerError
	PrintTokens()
}

type lexerWrapper struct {
	*antlr.DefaultErrorListener // Embed default which ensures we fit the interface
	antrlLexer *GoliteLexer
	stream     *antlr.CommonTokenStream
	errors     []*context.CompilerError
}

func (lexer *lexerWrapper) GetTokenStream() *antlr.CommonTokenStream {
	return lexer.stream
}

func (lexer *lexerWrapper) GetErrors() []*context.CompilerError {
	return lexer.errors
}

func (lexer *lexerWrapper) PrintTokens() {
	lexer.stream.Fill()
	for _, token := range lexer.stream.GetAllTokens() {
		if token.GetTokenType() != antlr.TokenEOF {
			if (lexer.antrlLexer.SymbolicNames[token.GetTokenType()] == "IDENTIFIER") {
				fmt.Printf("TOKEN.IDENTIFIER\t(%v, %v)\n", token.GetLine(), token.GetText())
			} else if (lexer.antrlLexer.SymbolicNames[token.GetTokenType()] == "NUMBER") {
				fmt.Printf("TOKEN.NUMBER\t\t(%v, %v)\n", token.GetLine(), token.GetText())
			} else if (lexer.antrlLexer.SymbolicNames[token.GetTokenType()] == "STRING") {
				fmt.Printf("TOKEN.STRING\t\t(%v, %v)\n", token.GetLine(), token.GetText())
			} else {
				fmt.Printf("TOKEN.%v\t\t(%v)\n", lexer.antrlLexer.SymbolicNames[token.GetTokenType()], token.GetLine())
			}
		}
	}
}

func (lexer *lexerWrapper) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	lexer.errors = append(lexer.errors, &context.CompilerError{
		Line:   line,
		Column: column,
		Msg:    msg,
		Phase:  context.LEXER,
	})
}

func NewLexer(inputSourcePath string) Lexer {
	input, _ := antlr.NewFileStream(inputSourcePath)
	lexer := &lexerWrapper{antlr.NewDefaultErrorListener(), nil, nil, nil}
	antlrLexer := NewGoliteLexer(input)
	antlrLexer.RemoveErrorListeners()
	antlrLexer.AddErrorListener(lexer)
	lexer.antrlLexer = antlrLexer
	lexer.stream = antlr.NewCommonTokenStream(antlrLexer, 0)
	return lexer
}

package parser

import ("golite/ast" ; "golite/context" ; "golite/lexer" ; "golite/token" ; "golite/types" ; 
		"fmt" ; "strconv" ; "github.com/antlr/antlr4/runtime/Go/antlr/v4")

type Parser interface {
	Parse() *ast.Program
	GetErrors() []*context.CompilerError
}

type parserWrapper struct {
	*antlr.DefaultErrorListener // Embed default which ensures we fit the interface
	*BaseGoliteParserListener
	antlrParser *GoliteParser
	lexer       lexer.Lexer
	errors      []*context.CompilerError
	nodes       map[string]interface{}
}

func (parser *parserWrapper) GetErrors() []*context.CompilerError {
	return parser.errors
}

func (parser *parserWrapper) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	parser.errors = append(parser.errors, &context.CompilerError{
		Line:   line,
		Column: column,
		Msg:    msg,
		Phase:  context.PARSER,
	})
}

func NewParser(lexer lexer.Lexer) Parser {
	parser := &parserWrapper{antlr.NewDefaultErrorListener(), &BaseGoliteParserListener{}, nil, lexer, nil, make(map[string]interface{})}
	antlrParser := NewGoliteParser(lexer.GetTokenStream())
	antlrParser.RemoveErrorListeners()
	antlrParser.AddErrorListener(parser)
	parser.antlrParser = antlrParser
	return parser
}

func (parser *parserWrapper) Parse() *ast.Program {
	ctx := parser.antlrParser.Program()

	if context.HasErrors(parser.lexer.GetErrors()) || context.HasErrors(parser.GetErrors()) { return nil }
	antlr.ParseTreeWalkerDefault.Walk(parser, ctx)
	_, _, key := GetTokenInfo(ctx)
	return parser.nodes[key].(*ast.Program)
}

/******************** Implementation of the Listeners ********************/

// GetTokenInfo gerates a unique key for each node in the ParserTree
func GetTokenInfo(ctx antlr.ParserRuleContext) (int, int, string) {
	key := fmt.Sprintf("%d,%d", ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())
	return ctx.GetStart().GetLine(), ctx.GetStart().GetColumn(), key
}

func (parser *parserWrapper) ExitProgram(ctx *ProgramContext) {
	line, column, key := GetTokenInfo(ctx)
	var types []*ast.TypeDeclaration
	var declarations []*ast.Declaration
	var functions []*ast.Function

	_, _, typesKey := GetTokenInfo(ctx.Types())
	if value, exists := parser.nodes[typesKey].([]*ast.TypeDeclaration); exists { types = value }

	_, _, declarationsKey := GetTokenInfo(ctx.Declarations())
	if value, exists := parser.nodes[declarationsKey].([]*ast.Declaration); exists { declarations = value }

	_, _, functionsKey := GetTokenInfo(ctx.Functions())
	if value, exists := parser.nodes[functionsKey].([]*ast.Function); exists { functions = value }

	parser.nodes[key] = ast.NewProgram(types, declarations, functions, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitTypes(ctx *TypesContext) {
	_, _, key := GetTokenInfo(ctx)
	typeDeclContexts := ctx.AllTypeDeclaration()
	var typeDecls []*ast.TypeDeclaration

	for _, typeDeclCtx := range typeDeclContexts {
		_, _, typeDeclKey := GetTokenInfo(typeDeclCtx)
		typeDecls = append(typeDecls, parser.nodes[typeDeclKey].(*ast.TypeDeclaration))
	}

	parser.nodes[key] = typeDecls
}

func (parser *parserWrapper) ExitTypeDeclaration(ctx *TypeDeclarationContext) {
	line, column, key := GetTokenInfo(ctx)

	variable := ast.NewVariable(ctx.IDENTIFIER().GetText(), token.NewToken(line, column))
	var fields []*ast.Decl

	_, _, fieldsKey := GetTokenInfo(ctx.Fields())
	if value, exists := parser.nodes[fieldsKey].([]*ast.Decl); exists { fields = value} 

	parser.nodes[key] = ast.NewTypeDeclaration(variable, fields, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitFields(ctx *FieldsContext) {
	_, _, key := GetTokenInfo(ctx)
	declContexts := ctx.AllDecl()
	var decls []*ast.Decl

	for _, declCtx := range declContexts {
		_, _, declKey := GetTokenInfo(declCtx)
		decls = append(decls, parser.nodes[declKey].(*ast.Decl))
	}

	parser.nodes[key] = decls
}

func (parser *parserWrapper) ExitDecl(ctx *DeclContext) {
	line, column, key := GetTokenInfo(ctx)
	
	variable := ast.NewVariable(ctx.IDENTIFIER().GetText(), token.NewToken(line, column))
	var ty types.Type

	_, _, typeKey := GetTokenInfo(ctx.Type_())
	if value, exists := parser.nodes[typeKey].(types.Type); exists { ty = value }

	parser.nodes[key] = ast.NewDecl(variable, ty, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitType(ctx *TypeContext) {
	_, _, key := GetTokenInfo(ctx)

	if intType := ctx.INT(); intType != nil {
		parser.nodes[key] = types.IntTySig

	} else if boolType := ctx.BOOL(); boolType != nil {
		parser.nodes[key] = types.BoolTySig

	} else if pointerType := ctx.ASTERISK(); pointerType != nil {
		parser.nodes[key] = &types.PointerTy{&types.StructTy{ctx.IDENTIFIER().GetText()}}
	}
}

func (parser *parserWrapper) ExitDeclarations(ctx *DeclarationsContext) {
	_, _, key := GetTokenInfo(ctx)
	declarationContexts := ctx.AllDeclaration()
	var declarations []*ast.Declaration

	for _, declarationCtx := range declarationContexts {
		_, _, declarationKey := GetTokenInfo(declarationCtx)
		declarations = append(declarations, parser.nodes[declarationKey].(*ast.Declaration))
	}

	parser.nodes[key] = declarations
}

func (parser *parserWrapper) ExitDeclaration(ctx *DeclarationContext) {
	line, column, key := GetTokenInfo(ctx)
	var variables []*ast.Variable
	var ty types.Type

	_, _, idsKey := GetTokenInfo(ctx.Ids())
	if value, exists := parser.nodes[idsKey].([]*ast.Variable); exists { variables = value }

	_, _, typeKey := GetTokenInfo(ctx.Type_())
	if value, exists := parser.nodes[typeKey].(types.Type); exists { ty = value }

	parser.nodes[key] = ast.NewDeclaration(variables, ty, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitIds(ctx *IdsContext) {
	line, column, key := GetTokenInfo(ctx)
	identifierContexts := ctx.AllIDENTIFIER()
	var variables []*ast.Variable
	
	for _, idCtx := range identifierContexts {
		variable := ast.NewVariable(idCtx.GetText(), token.NewToken(line, column))
		variables = append(variables, variable)
	}

	parser.nodes[key] = variables
}

func (parser *parserWrapper) ExitFunctions(ctx *FunctionsContext) {
	_, _, key := GetTokenInfo(ctx)
	functionContexts := ctx.AllFunction()
	var functions []*ast.Function

	for _, funcCtx := range functionContexts {
		_, _, funcKey := GetTokenInfo(funcCtx)
		functions = append(functions, parser.nodes[funcKey].(*ast.Function))
	}

	parser.nodes[key] = functions
}

func (parser *parserWrapper) ExitFunction(ctx *FunctionContext) {
	line, column, key := GetTokenInfo(ctx)
	variable := ast.NewVariable(ctx.IDENTIFIER().GetText(), token.NewToken(line, column))

	var ty types.Type
	var parameters []*ast.Decl
	var declarations []*ast.Declaration
	var statements []ast.Statement

	_, _, paramsKey := GetTokenInfo(ctx.Parameters())
	if value, exists := parser.nodes[paramsKey].([]*ast.Decl); exists { parameters = value }

	if ctx.ReturnType() != nil {
		_, _, typeKey := GetTokenInfo(ctx.ReturnType())
		if value, exists := parser.nodes[typeKey].(types.Type); exists { ty = value }
	}

	_, _, declsKey := GetTokenInfo(ctx.Declarations()) 
	if value, exists := parser.nodes[declsKey].([]*ast.Declaration); exists { declarations = value }
		
	_, _, stmtsKey := GetTokenInfo(ctx.Statements())
	if value, exists := parser.nodes[stmtsKey].([]ast.Statement); exists { statements = value }

	parser.nodes[key] = ast.NewFunction(variable, parameters, ty, declarations, statements, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitParameters(ctx *ParametersContext) {
	_, _, key := GetTokenInfo(ctx)
	declContexts := ctx.AllDecl()
	var decls []*ast.Decl

	for _, declCtx := range declContexts {
		_, _, declKey := GetTokenInfo(declCtx)
		decls = append(decls, parser.nodes[declKey].(*ast.Decl))
	}

	parser.nodes[key] = decls
}

func (parser *parserWrapper) ExitReturnType(ctx *ReturnTypeContext) {
	_, _, key := GetTokenInfo(ctx)
	_, _, typeKey := GetTokenInfo(ctx.Type_())
	if value, exists := parser.nodes[typeKey].(types.Type); exists { parser.nodes[key] = value }
}

func (parser *parserWrapper) ExitStatements(ctx *StatementsContext) {
	_, _, key := GetTokenInfo(ctx)
	statementContexts := ctx.AllStatement()
	var statements []ast.Statement

	for _, stmtCtx := range statementContexts {
		_, _, stmtKey := GetTokenInfo(stmtCtx)
		statements = append(statements, parser.nodes[stmtKey].(ast.Statement))
	}

	parser.nodes[key] = statements
}

func (parser *parserWrapper) ExitStatement(ctx *StatementContext) {
	_, _, key := GetTokenInfo(ctx)
	
	if idAssignment := ctx.Assignment(); idAssignment != nil {
		_, _, assignmentKey := GetTokenInfo(idAssignment)
		parser.nodes[key] = parser.nodes[assignmentKey]

	} else if idPrint := ctx.Print_(); idPrint != nil {
		_, _, printKey := GetTokenInfo(idPrint)
		parser.nodes[key] = parser.nodes[printKey]

	} else if idRead := ctx.Read(); idRead != nil {
		_, _, readKey := GetTokenInfo(idRead)
		parser.nodes[key] = parser.nodes[readKey]

	} else if idDelete := ctx.Delete_(); idDelete != nil {
		_, _, deleteKey := GetTokenInfo(idDelete)
		parser.nodes[key] = parser.nodes[deleteKey]

	} else if idConditional := ctx.Conditional(); idConditional != nil {
		_, _, conditionalKey := GetTokenInfo(idConditional)
		parser.nodes[key] = parser.nodes[conditionalKey]

	} else if idLoop := ctx.Loop(); idLoop != nil {
		_, _, loopKey := GetTokenInfo(idLoop)
		parser.nodes[key] = parser.nodes[loopKey]

	} else if idReturn := ctx.Return_(); idReturn != nil {
		_, _, returnKey := GetTokenInfo(idReturn)
		parser.nodes[key] = parser.nodes[returnKey]

	} else if idInvocation := ctx.Invocation(); idInvocation != nil {
		_, _, invocationKey := GetTokenInfo(idInvocation)
		parser.nodes[key] = parser.nodes[invocationKey]
	}
}

func (parser *parserWrapper) ExitBlock(ctx *BlockContext) {
	_, _, key := GetTokenInfo(ctx)
	_, _, stmtsKey := GetTokenInfo(ctx.Statements())
	if value, exists := parser.nodes[stmtsKey].([]ast.Statement); exists { parser.nodes[key] = value }
}

func (parser *parserWrapper) ExitDelete(ctx *DeleteContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression

	_, _, exprKey := GetTokenInfo(ctx.Expression())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	parser.nodes[key] = ast.NewDelete(expression, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitRead(ctx *ReadContext) {
	line, column, key := GetTokenInfo(ctx)
	var lvalue *ast.LValue

	_, _, lvalKey := GetTokenInfo(ctx.LValue())
	if value, exists := parser.nodes[lvalKey].(*ast.LValue); exists { lvalue = value }

	parser.nodes[key] = ast.NewRead(lvalue, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitAssignment(ctx *AssignmentContext) {
	line, column, key := GetTokenInfo(ctx)
	var lvalue *ast.LValue
	var expression ast.Expression

	_, _, lvalKey := GetTokenInfo(ctx.LValue())
	if value, exists := parser.nodes[lvalKey].(*ast.LValue); exists { lvalue = value }

	_, _, exprKey := GetTokenInfo(ctx.Expression())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	parser.nodes[key] = ast.NewAssignment(lvalue, expression, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitPrint(ctx *PrintContext) {
	line, column, key := GetTokenInfo(ctx)
	expressions := ctx.AllExpression()
	var exprs []ast.Expression

	for _, expression := range expressions {
		_, _, exprKey := GetTokenInfo(expression)
		exprs = append(exprs, parser.nodes[exprKey].(ast.Expression))
	}

	parser.nodes[key] = ast.NewPrint(ctx.STRING().GetText(), exprs, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitConditional(ctx *ConditionalContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression
	var ifBlock []ast.Statement
	var elseBlock []ast.Statement

	_, _, exprKey := GetTokenInfo(ctx.Expression())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	_, _, ifBlockKey := GetTokenInfo(ctx.Block(0))
	if value, exists := parser.nodes[ifBlockKey].([]ast.Statement); exists { ifBlock = value }

	if ctx.Block(1) != nil {
		_, _, elseBlockKey := GetTokenInfo(ctx.Block(1))
		if value, exists := parser.nodes[elseBlockKey].([]ast.Statement); exists { elseBlock = value }
	}

	parser.nodes[key] = ast.NewConditional(expression, ifBlock, elseBlock, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitLoop(ctx *LoopContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression
	var forBlock []ast.Statement

	_, _, exprKey := GetTokenInfo(ctx.Expression())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	_, _, blockKey := GetTokenInfo(ctx.Block())
	if value, exists := parser.nodes[blockKey].([]ast.Statement); exists { forBlock = value }

	parser.nodes[key] = ast.NewLoop(expression, forBlock, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitReturn(ctx *ReturnContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression

	if ctx.Expression() != nil {
		_, _, exprKey := GetTokenInfo(ctx.Expression())
		if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }
	}
	
	parser.nodes[key] = ast.NewReturn(expression, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitInvocation(ctx *InvocationContext) {
	line, column, key := GetTokenInfo(ctx)

	variable := ast.NewVariable(ctx.IDENTIFIER().GetText(), token.NewToken(line, column))
	var arguments []ast.Expression

	_, _, argsKey := GetTokenInfo(ctx.Arguments())
	if value, exists := parser.nodes[argsKey].([]ast.Expression); exists { arguments = value }

	parser.nodes[key] = ast.NewInvocation(variable, arguments, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitArguments(ctx *ArgumentsContext) {
	_, _, key := GetTokenInfo(ctx)
	expressionContexts := ctx.AllExpression()
	var expressions []ast.Expression
	
	for _, exprCtx := range expressionContexts {
		_, _, exprKey := GetTokenInfo(exprCtx)
		expressions = append(expressions, parser.nodes[exprKey].(ast.Expression))
	}

	parser.nodes[key] = expressions
}

func (parser *parserWrapper) ExitLValue(ctx *LValueContext) {
	line, column, key := GetTokenInfo(ctx)
	identifiers := ctx.AllIDENTIFIER()
	var fields []*ast.Variable

	variable := ast.NewVariable(identifiers[0].GetText(), token.NewToken(line, column))
    
	for _, identifier := range identifiers[1:] {
		field := ast.NewVariable(identifier.GetText(), token.NewToken(line, column))
		fields = append(fields, field)
	}
	
	parser.nodes[key] = ast.NewLValue(variable, fields, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitExpression(ctx *ExpressionContext) {
	line, column, key := GetTokenInfo(ctx)
	boolTerms := ctx.AllBoolTerm()

	_, _, exprKey := GetTokenInfo(boolTerms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)

	for _, boolTerm := range boolTerms[1:] {
		_, _, nextExprKey := GetTokenInfo(boolTerm)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)
		currExpr = ast.NewBinaryExpr(ast.OR, currExpr, nextExpr, token.NewToken(line, column))
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitBoolTerm(ctx *BoolTermContext) {
	line, column, key := GetTokenInfo(ctx)
	equalTerms := ctx.AllEqualTerm()

	_, _, exprKey := GetTokenInfo(equalTerms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)

	for _, equalTerm := range equalTerms[1:] {
		_, _, nextExprKey := GetTokenInfo(equalTerm)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)
		currExpr = ast.NewBinaryExpr(ast.AND, currExpr, nextExpr, token.NewToken(line, column))
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitEqualTerm(ctx *EqualTermContext) {
	line, column, key := GetTokenInfo(ctx)
	relationTerms := ctx.AllRelationTerm()

	_, _, exprKey := GetTokenInfo(relationTerms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)
	var a, b int

	for _, relationTerm := range relationTerms[1:] {
		_, _, nextExprKey := GetTokenInfo(relationTerm)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)

		if eqOperator := ctx.EQ(a); eqOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.EQ, currExpr, nextExpr, token.NewToken(line, column))

		} else if neqOperator := ctx.NEQ(b); neqOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.NEQ, currExpr, nextExpr, token.NewToken(line, column))
		}
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitRelationTerm(ctx *RelationTermContext) {
	line, column, key := GetTokenInfo(ctx)
	simpleTerms := ctx.AllSimpleTerm()

	_, _, exprKey := GetTokenInfo(simpleTerms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)
	var a, b, c, d int

	for _, simpleTerm := range simpleTerms[1:] {
		_, _, nextExprKey := GetTokenInfo(simpleTerm)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)

		if gtOperator := ctx.GT(a); gtOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.GT, currExpr, nextExpr, token.NewToken(line, column))

		} else if ltOperator := ctx.LT(b); ltOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.LT, currExpr, nextExpr, token.NewToken(line, column))

		} else if geqOperator := ctx.GEQ(c); geqOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.GEQ, currExpr, nextExpr, token.NewToken(line, column))

		} else if leqOperator := ctx.LEQ(d); leqOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.LEQ, currExpr, nextExpr, token.NewToken(line, column))
		}
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitSimpleTerm(ctx *SimpleTermContext) {
	line, column, key := GetTokenInfo(ctx)
	terms := ctx.AllTerm()

	_, _, exprKey := GetTokenInfo(terms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)
	var a, b int

	for _, term := range terms[1:] {
		_, _, nextExprKey := GetTokenInfo(term)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)

		if plusOperator := ctx.PLUS(a); plusOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.ADD, currExpr, nextExpr, token.NewToken(line, column))
			a++

		} else if minusOperator := ctx.MINUS(b); minusOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.SUBTRACT, currExpr, nextExpr, token.NewToken(line, column))
			b++
		}
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitTerm(ctx *TermContext) {
	line, column, key := GetTokenInfo(ctx)
	unaryTerms := ctx.AllUnaryTerm()
	var a, b int

	_, _, exprKey := GetTokenInfo(unaryTerms[0])
	currExpr := parser.nodes[exprKey].(ast.Expression)

	for _, unaryTerm := range unaryTerms[1:] {
		_, _, nextExprKey := GetTokenInfo(unaryTerm)
		nextExpr := parser.nodes[nextExprKey].(ast.Expression)

		if multiplyOperator := ctx.ASTERISK(a); multiplyOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.MULTIPLY, currExpr, nextExpr, token.NewToken(line, column))
			a++

		} else if divideOperator := ctx.FSLASH(b); divideOperator != nil {
			currExpr = ast.NewBinaryExpr(ast.DIVIDE, currExpr, nextExpr, token.NewToken(line, column))
			b++
		}
	}

    parser.nodes[key] = currExpr
}

func (parser *parserWrapper) ExitUnaryTerm(ctx *UnaryTermContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression

	_, _, exprKey := GetTokenInfo(ctx.SelectorTerm())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	if notOperator := ctx.NOT(); notOperator != nil {
		parser.nodes[key] = ast.NewUnaryExpr(ast.NOT, expression, token.NewToken(line, column))

	} else if minusOperator := ctx.MINUS(); minusOperator != nil {
		parser.nodes[key] = ast.NewUnaryExpr(ast.SUBTRACT, expression, token.NewToken(line, column))

	} else { parser.nodes[key] = expression }
}

func (parser *parserWrapper) ExitSelectorTerm(ctx *SelectorTermContext) {
	line, column, key := GetTokenInfo(ctx)
	var expression ast.Expression
	var fields []*ast.Variable

	_, _, exprKey := GetTokenInfo(ctx.Factor())
	if value, exists := parser.nodes[exprKey].(ast.Expression); exists { expression = value }

	for _, identifier := range ctx.AllIDENTIFIER() {
		variable := ast.NewVariable(identifier.GetText(), token.NewToken(line, column))
		fields = append(fields, variable)
	}

	parser.nodes[key] = ast.NewFieldExpr(expression, fields, token.NewToken(line, column))
}

func (parser *parserWrapper) ExitFactor(ctx *FactorContext) {
	line, column, key := GetTokenInfo(ctx)

	if exprFactor := ctx.Expression(); exprFactor != nil {
		_, _, exprKey := GetTokenInfo(exprFactor)
		expression := parser.nodes[exprKey].(ast.Expression)
		parser.nodes[key] = expression

	} else if idFactor := ctx.IDENTIFIER(); idFactor != nil {
		variable := ast.NewVariable(idFactor.GetText(), token.NewToken(line, column))

		if argsFactor := ctx.Arguments(); argsFactor != nil {
			_, _, argsKey := GetTokenInfo(argsFactor)
			arguments := parser.nodes[argsKey].([]ast.Expression)
			parser.nodes[key] = ast.NewCallExpr(variable, arguments, token.NewToken(line, column))

		} else if newFactor := ctx.NEW(); newFactor != nil {
			parser.nodes[key] = ast.NewNewExpr(variable, token.NewToken(line, column))

		} else { parser.nodes[key] = variable }
	
	} else if intFactor := ctx.NUMBER(); intFactor != nil {
		integer, _ := strconv.Atoi(intFactor.GetText())
		parser.nodes[key] = ast.NewIntLiteral(int64(integer), token.NewToken(line, column))

	} else if trueFactor := ctx.TRUE(); trueFactor != nil {
		parser.nodes[key] = ast.NewBoolLiteral(true, token.NewToken(line, column))

	} else if falseFactor := ctx.FALSE(); falseFactor != nil {
		parser.nodes[key] = ast.NewBoolLiteral(false, token.NewToken(line, column))

	} else if nilFactor := ctx.NIL(); nilFactor != nil {
		parser.nodes[key] = ast.NewNilLiteral(token.NewToken(line, column))
	}
}

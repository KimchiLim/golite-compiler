// Generated from /home/ahuang02/golite-lets-byte/grammars/GoliteParser.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.tree.ParseTreeListener;

/**
 * This interface defines a complete listener for a parse tree produced by
 * {@link GoliteParser}.
 */
public interface GoliteParserListener extends ParseTreeListener {
	/**
	 * Enter a parse tree produced by {@link GoliteParser#program}.
	 * @param ctx the parse tree
	 */
	void enterProgram(GoliteParser.ProgramContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#program}.
	 * @param ctx the parse tree
	 */
	void exitProgram(GoliteParser.ProgramContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#types}.
	 * @param ctx the parse tree
	 */
	void enterTypes(GoliteParser.TypesContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#types}.
	 * @param ctx the parse tree
	 */
	void exitTypes(GoliteParser.TypesContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#typeDeclaration}.
	 * @param ctx the parse tree
	 */
	void enterTypeDeclaration(GoliteParser.TypeDeclarationContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#typeDeclaration}.
	 * @param ctx the parse tree
	 */
	void exitTypeDeclaration(GoliteParser.TypeDeclarationContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#fields}.
	 * @param ctx the parse tree
	 */
	void enterFields(GoliteParser.FieldsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#fields}.
	 * @param ctx the parse tree
	 */
	void exitFields(GoliteParser.FieldsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#decl}.
	 * @param ctx the parse tree
	 */
	void enterDecl(GoliteParser.DeclContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#decl}.
	 * @param ctx the parse tree
	 */
	void exitDecl(GoliteParser.DeclContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#type}.
	 * @param ctx the parse tree
	 */
	void enterType(GoliteParser.TypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#type}.
	 * @param ctx the parse tree
	 */
	void exitType(GoliteParser.TypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#declarations}.
	 * @param ctx the parse tree
	 */
	void enterDeclarations(GoliteParser.DeclarationsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#declarations}.
	 * @param ctx the parse tree
	 */
	void exitDeclarations(GoliteParser.DeclarationsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#declaration}.
	 * @param ctx the parse tree
	 */
	void enterDeclaration(GoliteParser.DeclarationContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#declaration}.
	 * @param ctx the parse tree
	 */
	void exitDeclaration(GoliteParser.DeclarationContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#ids}.
	 * @param ctx the parse tree
	 */
	void enterIds(GoliteParser.IdsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#ids}.
	 * @param ctx the parse tree
	 */
	void exitIds(GoliteParser.IdsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#functions}.
	 * @param ctx the parse tree
	 */
	void enterFunctions(GoliteParser.FunctionsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#functions}.
	 * @param ctx the parse tree
	 */
	void exitFunctions(GoliteParser.FunctionsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#function}.
	 * @param ctx the parse tree
	 */
	void enterFunction(GoliteParser.FunctionContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#function}.
	 * @param ctx the parse tree
	 */
	void exitFunction(GoliteParser.FunctionContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#parameters}.
	 * @param ctx the parse tree
	 */
	void enterParameters(GoliteParser.ParametersContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#parameters}.
	 * @param ctx the parse tree
	 */
	void exitParameters(GoliteParser.ParametersContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#returnType}.
	 * @param ctx the parse tree
	 */
	void enterReturnType(GoliteParser.ReturnTypeContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#returnType}.
	 * @param ctx the parse tree
	 */
	void exitReturnType(GoliteParser.ReturnTypeContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#statements}.
	 * @param ctx the parse tree
	 */
	void enterStatements(GoliteParser.StatementsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#statements}.
	 * @param ctx the parse tree
	 */
	void exitStatements(GoliteParser.StatementsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#statement}.
	 * @param ctx the parse tree
	 */
	void enterStatement(GoliteParser.StatementContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#statement}.
	 * @param ctx the parse tree
	 */
	void exitStatement(GoliteParser.StatementContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#block}.
	 * @param ctx the parse tree
	 */
	void enterBlock(GoliteParser.BlockContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#block}.
	 * @param ctx the parse tree
	 */
	void exitBlock(GoliteParser.BlockContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#delete}.
	 * @param ctx the parse tree
	 */
	void enterDelete(GoliteParser.DeleteContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#delete}.
	 * @param ctx the parse tree
	 */
	void exitDelete(GoliteParser.DeleteContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#read}.
	 * @param ctx the parse tree
	 */
	void enterRead(GoliteParser.ReadContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#read}.
	 * @param ctx the parse tree
	 */
	void exitRead(GoliteParser.ReadContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#assignment}.
	 * @param ctx the parse tree
	 */
	void enterAssignment(GoliteParser.AssignmentContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#assignment}.
	 * @param ctx the parse tree
	 */
	void exitAssignment(GoliteParser.AssignmentContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#print}.
	 * @param ctx the parse tree
	 */
	void enterPrint(GoliteParser.PrintContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#print}.
	 * @param ctx the parse tree
	 */
	void exitPrint(GoliteParser.PrintContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#conditional}.
	 * @param ctx the parse tree
	 */
	void enterConditional(GoliteParser.ConditionalContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#conditional}.
	 * @param ctx the parse tree
	 */
	void exitConditional(GoliteParser.ConditionalContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#loop}.
	 * @param ctx the parse tree
	 */
	void enterLoop(GoliteParser.LoopContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#loop}.
	 * @param ctx the parse tree
	 */
	void exitLoop(GoliteParser.LoopContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#return}.
	 * @param ctx the parse tree
	 */
	void enterReturn(GoliteParser.ReturnContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#return}.
	 * @param ctx the parse tree
	 */
	void exitReturn(GoliteParser.ReturnContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#invocation}.
	 * @param ctx the parse tree
	 */
	void enterInvocation(GoliteParser.InvocationContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#invocation}.
	 * @param ctx the parse tree
	 */
	void exitInvocation(GoliteParser.InvocationContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#arguments}.
	 * @param ctx the parse tree
	 */
	void enterArguments(GoliteParser.ArgumentsContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#arguments}.
	 * @param ctx the parse tree
	 */
	void exitArguments(GoliteParser.ArgumentsContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#lValue}.
	 * @param ctx the parse tree
	 */
	void enterLValue(GoliteParser.LValueContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#lValue}.
	 * @param ctx the parse tree
	 */
	void exitLValue(GoliteParser.LValueContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#expression}.
	 * @param ctx the parse tree
	 */
	void enterExpression(GoliteParser.ExpressionContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#expression}.
	 * @param ctx the parse tree
	 */
	void exitExpression(GoliteParser.ExpressionContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#boolTerm}.
	 * @param ctx the parse tree
	 */
	void enterBoolTerm(GoliteParser.BoolTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#boolTerm}.
	 * @param ctx the parse tree
	 */
	void exitBoolTerm(GoliteParser.BoolTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#equalTerm}.
	 * @param ctx the parse tree
	 */
	void enterEqualTerm(GoliteParser.EqualTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#equalTerm}.
	 * @param ctx the parse tree
	 */
	void exitEqualTerm(GoliteParser.EqualTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#relationTerm}.
	 * @param ctx the parse tree
	 */
	void enterRelationTerm(GoliteParser.RelationTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#relationTerm}.
	 * @param ctx the parse tree
	 */
	void exitRelationTerm(GoliteParser.RelationTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#simpleTerm}.
	 * @param ctx the parse tree
	 */
	void enterSimpleTerm(GoliteParser.SimpleTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#simpleTerm}.
	 * @param ctx the parse tree
	 */
	void exitSimpleTerm(GoliteParser.SimpleTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#term}.
	 * @param ctx the parse tree
	 */
	void enterTerm(GoliteParser.TermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#term}.
	 * @param ctx the parse tree
	 */
	void exitTerm(GoliteParser.TermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#unaryTerm}.
	 * @param ctx the parse tree
	 */
	void enterUnaryTerm(GoliteParser.UnaryTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#unaryTerm}.
	 * @param ctx the parse tree
	 */
	void exitUnaryTerm(GoliteParser.UnaryTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#selectorTerm}.
	 * @param ctx the parse tree
	 */
	void enterSelectorTerm(GoliteParser.SelectorTermContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#selectorTerm}.
	 * @param ctx the parse tree
	 */
	void exitSelectorTerm(GoliteParser.SelectorTermContext ctx);
	/**
	 * Enter a parse tree produced by {@link GoliteParser#factor}.
	 * @param ctx the parse tree
	 */
	void enterFactor(GoliteParser.FactorContext ctx);
	/**
	 * Exit a parse tree produced by {@link GoliteParser#factor}.
	 * @param ctx the parse tree
	 */
	void exitFactor(GoliteParser.FactorContext ctx);
}
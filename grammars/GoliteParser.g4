parser grammar GoliteParser;

options {
    tokenVocab = GoliteLexer;
}

program: types declarations functions EOF;
types: typeDeclaration*;
typeDeclaration: TYPE IDENTIFIER STRUCT LBRACE fields RBRACE SEMICOLON;
fields: decl SEMICOLON (decl SEMICOLON)*;
decl: IDENTIFIER type;
type: INT | BOOL | ASTERISK IDENTIFIER;
declarations: declaration*;
declaration: VAR ids type SEMICOLON;
ids: IDENTIFIER (COMMA IDENTIFIER)*;
functions: function*;
function: FUNC IDENTIFIER parameters returnType? LBRACE declarations statements RBRACE;
parameters: LPAREN (decl (COMMA decl)*)? RPAREN;
returnType: type;
statements: statement*;
statement: assignment | print | read | delete | conditional | loop | return | invocation;
block: LBRACE statements RBRACE;
delete: DELETE expression SEMICOLON;
read: SCAN lValue SEMICOLON;
assignment: lValue EQUAL expression SEMICOLON;
print: PRINTF LPAREN STRING (COMMA expression)* RPAREN SEMICOLON;
conditional: IF LPAREN expression RPAREN block (ELSE block)?;
loop: FOR LPAREN expression RPAREN block;
return: RETURN expression? SEMICOLON;
invocation: IDENTIFIER arguments SEMICOLON;
arguments: LPAREN (expression (COMMA expression)*)? RPAREN;
lValue: IDENTIFIER (PERIOD IDENTIFIER)*;
expression: boolTerm (OR boolTerm)*;
boolTerm: equalTerm (AND equalTerm)*;
equalTerm: relationTerm ((EQ | NEQ) relationTerm)*;
relationTerm: simpleTerm ((GT | LT | GEQ | LEQ) simpleTerm)*;
simpleTerm: term ((PLUS | MINUS) term)*;
term: unaryTerm ((ASTERISK | FSLASH) unaryTerm)*;
unaryTerm: NOT selectorTerm | MINUS selectorTerm | selectorTerm;
selectorTerm: factor (PERIOD IDENTIFIER)*;
factor: LPAREN expression RPAREN | IDENTIFIER arguments? | NUMBER | NEW IDENTIFIER | TRUE | FALSE | NIL;
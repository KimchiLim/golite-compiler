lexer grammar GoliteLexer;

FUNC: 'func';
TYPE: 'type';
STRUCT: 'struct';
INT: 'int';
BOOL: 'bool';
VAR: 'var';
IF: 'if';
ELSE: 'else';
FOR: 'for';
RETURN: 'return';

LPAREN: '(';
RPAREN: ')';
LBRACE: '{';
RBRACE: '}';
COMMA: ',';
PERIOD: '.';
SEMICOLON: ';';
PLUS: '+';
MINUS: '-';
ASTERISK: '*';
FSLASH: '/';
EQUAL: '=';

OR: '||';
AND: '&&';
EQ: '==';
NEQ: '!=';
GT: '>';
LT: '<';
GEQ: '>=';
LEQ: '<=';
NOT: '!';
TRUE: 'true';
FALSE: 'false';
NIL: 'nil';

NEW: 'new';
DELETE: 'delete';
SCAN: 'scan';
PRINTF: 'printf';

IDENTIFIER: [a-zA-Z] [a-zA-Z0-9]*;
NUMBER: '0' | [1-9] [0-9]*;
STRING: '"' ~'"'* '"';
COMMENT: '//' ~'\n'* -> skip;
WHITESPACE: [ \r\n\t]+ -> skip;
package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/types")

type Operator int

const (ADD Operator = iota ; SUBTRACT ; MULTIPLY ; DIVIDE ; OR ; AND ; EQ ; NEQ ; LT ; GT ; LEQ ; GEQ ; NOT)

func StringToOperator(op string) Operator {
	switch op {
		case "+": return ADD
		case "-": return SUBTRACT
		case "*": return MULTIPLY
		case "/": return DIVIDE
		case "||": return OR
		case "&&": return AND
		case "==": return EQ
		case "!=": return NEQ
		case "<": return LT
		case ">": return GT
		case "<=": return LEQ
		case ">=": return GEQ
		case "!": return NOT
	}
	panic("Not found operator")
}

func OperatorString(op Operator) string {
	switch op {
		case ADD: return "+"
		case SUBTRACT: return "-"
		case MULTIPLY: return "*"
		case DIVIDE: return "/"
		case OR: return "||"
		case AND: return "&&"
		case EQ: return "=="
		case NEQ: return "!="
		case LT: return "<"
		case GT: return ">"
		case LEQ: return "<="
		case GEQ: return ">="
		case NOT: return "!"
	}
	panic("Not found operator")
}

type Expression interface {
	ExpressionNode()
	String() string
	TypeCheck([]*context.CompilerError, *st.SymbolTables, *st.FuncEntry, *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block)
	TranslateToLLVMStack(*st.FuncEntry, *llvm.Block, *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand)
}

func (n *NewExpr) ExpressionNode() {}
func (c *CallExpr) ExpressionNode() {}
func (b *BinaryExpr) ExpressionNode() {}
func (u *UnaryExpr) ExpressionNode() {}
func (f *FieldExpr) ExpressionNode() {}
func (l *LValue) ExpressionNode() {}
func (v *Variable) ExpressionNode() {}
func (i *IntLiteral) ExpressionNode() {}
func (b *BoolLiteral) ExpressionNode() {}
func (n *NilLiteral) ExpressionNode() {}

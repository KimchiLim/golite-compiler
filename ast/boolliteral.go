package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "fmt")

type BoolLiteral struct {
	*token.Token
	value bool
}

func NewBoolLiteral(value bool, token *token.Token) *BoolLiteral {
	return &BoolLiteral{token, value}
}

func (b *BoolLiteral) String() string {
	return fmt.Sprintf("%v", b.value)
}

func (b *BoolLiteral) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	return errors, types.BoolTySig, block
}

func (b *BoolLiteral) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand *llvm.LLVMImmediate
	var immediateValue int

	if b.value { immediateValue = 1 } else { immediateValue = 0 }
	operand = llvm.NewImmediate(immediateValue, types.BoolTySig)
    return instrs, operand
}

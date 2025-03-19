package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "fmt")

type IntLiteral struct {
	*token.Token
	value int64
}

func NewIntLiteral(value int64, token *token.Token) *IntLiteral {
	return &IntLiteral{token, value}
}

func (i *IntLiteral) String() string {
	return fmt.Sprintf("%d", i.value)
}

func (i *IntLiteral) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	return errors, types.IntTySig, block
}

func (i *IntLiteral) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand *llvm.LLVMImmediate
	var immediateValue int

	immediateValue = int(i.value)
	operand = llvm.NewImmediate(immediateValue, types.IntTySig)
    return instrs, operand
}

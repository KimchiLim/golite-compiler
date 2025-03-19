package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types")

type NilLiteral struct {
	*token.Token
}

func NewNilLiteral(token *token.Token) *NilLiteral {
	return &NilLiteral{token}
}

func (n *NilLiteral) String() string {
	return "nil"
}

func (n *NilLiteral) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	return errors, types.NilTySig, block
}

func (n *NilLiteral) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand *llvm.LLVMImmediate

	operand = llvm.NewImmediate(0, types.NilTySig)
    return instrs, operand
}

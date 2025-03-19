package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Read struct {
	*token.Token
	lvalue *LValue
}

func NewRead(lvalue *LValue, token *token.Token) *Read {
	return &Read{token, lvalue}
}

func (r *Read) String() string {
	var out bytes.Buffer
	out.WriteString("scan")
	out.WriteString(" ")
	out.WriteString(r.lvalue.String())
	out.WriteString(";")
	return out.String()
}

func (r *Read) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	errors, exprType, block = r.lvalue.TypeCheck(errors, tables, function, block)
	if exprType == types.IntTySig { return errors, types.VoidTySig, block }

	errors = st.SemanticError(errors, r.Token.Line, r.Token.Column, fmt.Sprintf("invalid read: expected int but got %s", exprType))
	return errors, types.UnknownTySig, block
}

func (r *Read) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand

	instrs, operand = r.lvalue.TranslateToLLVMStack(function, block, program)
	instrs = append(instrs, llvm.NewScan(operand))
	block.StackInstrs = append(block.StackInstrs, instrs...)
	return block
}

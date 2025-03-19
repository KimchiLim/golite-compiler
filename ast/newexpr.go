package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type NewExpr struct {
	*token.Token
	identifier *Variable
}

func NewNewExpr(identifier *Variable, token *token.Token) *NewExpr {
	return &NewExpr{token, identifier}
}

func (n *NewExpr) String() string {
	var out bytes.Buffer
    out.WriteString("new")
	out.WriteString(" ")
	out.WriteString(n.identifier.String())
	out.WriteString(";")
	return out.String()
}

func (n *NewExpr) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	structName := n.identifier.String()

	if structEntry := tables.Structs.Contains(structName); structEntry != nil {
		return errors, &types.PointerTy{structEntry.Ty}, block
	}
	errors = st.SemanticError(errors, n.Token.Line, n.Token.Column, fmt.Sprintf("invalid new expression of type *%s", structName))
	return errors, types.UnknownTySig, block
}

func (n *NewExpr) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	structName := n.identifier.String()
	exprType = types.UnknownTySig

	if structEntry := program.Tables.Structs.Contains(structName); structEntry != nil {
		exprType = &types.PointerTy{structEntry.Ty}
		entry = st.NewVarEntry(llvm.NewRegisterLabel(), types.UnknownTySig, st.LOCAL, n.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)
		instrs = append(instrs, llvm.NewMalloc(funcRegisters[entry], structEntry))
		operand = funcRegisters[entry]	

		entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, n.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)
		instrs = append(instrs, llvm.NewBitcast(funcRegisters[entry], types.UnknownTySig, operand, exprType))
		operand = funcRegisters[entry]	
	}

	return instrs, operand
}

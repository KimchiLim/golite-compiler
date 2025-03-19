package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Delete struct {
	*token.Token
	expression Expression
}

func NewDelete(expression Expression, token *token.Token) *Delete {
	return &Delete{token, expression}
}

func (d *Delete) String() string {
	var out bytes.Buffer
	out.WriteString("delete")
	out.WriteString(" ")
	out.WriteString(d.expression.String())
	out.WriteString(";")
	return out.String()
}

func (d *Delete) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	errors, exprType, block = d.expression.TypeCheck(errors, tables, function, block)

	if pointerType, ok := exprType.(*types.PointerTy); ok {
		if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
			if structEntry := tables.Structs.Contains(structType.Name); structEntry != nil {
				return errors, types.VoidTySig, block
			} else {
				errors = st.SemanticError(errors, d.Token.Line, d.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
				return errors, types.UnknownTySig, block
			}
		}
	}

	errors = st.SemanticError(errors, d.Token.Line, d.Token.Column, fmt.Sprintf("invalid delete: expected struct but got %s", exprType))
	return errors, types.UnknownTySig, block
}

func (d *Delete) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	instrs, operand = d.expression.TranslateToLLVMStack(function, block, program)
	exprType = operand.GetType()

	entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, d.Token)
	funcRegisters[entry] = llvm.NewRegister(entry)
	instrs = append(instrs, llvm.NewBitcast(funcRegisters[entry], exprType, operand, types.UnknownTySig))
	operand = funcRegisters[entry]	

	instrs = append(instrs, llvm.NewFree(operand))
	block.StackInstrs = append(block.StackInstrs, instrs...)	
	return block
}

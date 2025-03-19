package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes")

type Return struct {
	*token.Token
	expression Expression
}

func NewReturn(expression Expression, token *token.Token) *Return {
	return &Return{token, expression}
}

func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString("return")

	if (r.expression != nil) {
		out.WriteString(" ")
		out.WriteString(r.expression.String())
	}

	out.WriteString(";")
	return out.String()
}

func (r *Return) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	exprType = types.VoidTySig

	if r.expression != nil { errors, exprType, block = r.expression.TypeCheck(errors, tables, function, block) }
	block.AddReturn(exprType)
	return errors, exprType, block
}

func (r *Return) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type

	funcRegisters := program.FuncDefs[function].Registers
	exprType = types.VoidTySig

	if r.expression != nil {
		instrs, operand = r.expression.TranslateToLLVMStack(function, block, program)
		exprType = operand.GetType()

		if varEntry := function.Variables.Contains("_ret_val"); varEntry != nil {
			instrs = append(instrs, llvm.NewStore(exprType, operand, funcRegisters[varEntry]))
			instrs = append(instrs, llvm.NewReturn(exprType, operand))
		}
	} else if function.Name != "main" {
		instrs = append(instrs, llvm.NewReturn(exprType, operand))
	} else { instrs = append(instrs, llvm.NewReturn(types.IntTySig, llvm.NewImmediate(0, types.IntTySig))) }

	block.AddReturn(exprType)
	block.StackInstrs = append(block.StackInstrs, instrs...)
	return block
}

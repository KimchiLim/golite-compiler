package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type UnaryExpr struct {
	*token.Token
	operator Operator
	operand  Expression
}

func NewUnaryExpr(operator Operator, operand Expression, token *token.Token) *UnaryExpr {
	return &UnaryExpr{token, operator, operand}
}

func (u *UnaryExpr) String() string {
	var out bytes.Buffer
	out.WriteString(OperatorString(u.operator))
	out.WriteString(u.operand.String())	
	return out.String()
}

func (u *UnaryExpr) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	errors, exprType, block = u.operand.TypeCheck(errors, tables, function, block)

	if u.operator == NOT && exprType == types.BoolTySig {
		return errors, types.BoolTySig, block

	} else if u.operator == SUBTRACT && exprType == types.IntTySig {
		return errors, types.IntTySig, block
	}

	errors = st.SemanticError(errors, u.Token.Line, u.Token.Column, fmt.Sprintf("invalid unary expression: %s", u))
	return errors, types.UnknownTySig, block
}

func (u *UnaryExpr) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	instrs, operand = u.operand.TranslateToLLVMStack(function, block, program)
	exprType = operand.GetType()

	entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, u.Token)
	funcRegisters[entry] = llvm.NewRegister(entry)

	if u.operator == NOT {
		instrs = append(instrs, llvm.NewBinaryOp(funcRegisters[entry], llvm.XOR, exprType, operand, llvm.NewImmediate(1, types.BoolTySig)))
	} else if u.operator == SUBTRACT {
		instrs = append(instrs, llvm.NewBinaryOp(funcRegisters[entry], llvm.SUB, exprType, llvm.NewImmediate(0, types.IntTySig), operand))
	}

	operand = funcRegisters[entry]
	return instrs, operand
}

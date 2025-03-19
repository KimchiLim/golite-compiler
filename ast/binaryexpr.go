package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type BinaryExpr struct {
	*token.Token
	operator     Operator
	leftOperand  Expression
	rightOperand Expression
}

func NewBinaryExpr(operator Operator, leftOperand Expression, rightOperand Expression, token *token.Token) *BinaryExpr {
	return &BinaryExpr{token, operator, leftOperand, rightOperand}
}

func (b *BinaryExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(b.leftOperand.String())

	if b.rightOperand != nil {
		out.WriteString(" ")
		out.WriteString(OperatorString(b.operator))
		out.WriteString(" ")
		out.WriteString(b.rightOperand.String())
	}

	out.WriteString(")")
	return out.String()
}

func (b *BinaryExpr) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var leftType, rightType types.Type
	errors, leftType, block = b.leftOperand.TypeCheck(errors, tables, function, block)

	if b.rightOperand == nil { return errors, leftType, block }
	errors, rightType, block = b.rightOperand.TypeCheck(errors, tables, function, block)
	if leftType == types.UnknownTySig || rightType == types.UnknownTySig { return errors, types.UnknownTySig, block }

	if leftType == rightType {
		if (b.operator == ADD || b.operator == SUBTRACT || b.operator == MULTIPLY || b.operator == DIVIDE) && leftType == types.IntTySig {
			return errors, types.IntTySig, block

		} else if (b.operator == LT || b.operator == GT || b.operator == LEQ || b.operator == GEQ) && leftType == types.IntTySig {
			return errors, types.BoolTySig, block

		} else if _, ok := leftType.(*types.PointerTy); (b.operator == EQ || b.operator == NEQ) && (leftType == types.IntTySig || ok) {
			return errors, types.BoolTySig, block

		} else if (b.operator == OR || b.operator == AND) && leftType == types.BoolTySig {
			return errors, types.BoolTySig, block
		}
	} else if _, ok := leftType.(*types.PointerTy); rightType == types.NilTySig && ok {
		return errors, types.BoolTySig, block

	} else if _, ok := rightType.(*types.PointerTy); leftType == types.NilTySig && ok {
		return errors, types.BoolTySig, block
	} else if leftType.String() == rightType.String() {
		return errors, types.BoolTySig, block
	}

	errors = st.SemanticError(errors, b.Token.Line, b.Token.Column, fmt.Sprintf("invalid binary expression: %s", b))
	return errors, types.UnknownTySig, block
}

func (b *BinaryExpr) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs, instrss []llvm.LLVMInstruction
	var leftOperand, rightOperand, operand llvm.LLVMOperand
	var leftType, rightType, exprType types.Type
	var operator llvm.LLVMOperator
	var entry *st.VarEntry

	funcRegisters := program.FuncDefs[function].Registers
	exprType = types.Int1TySig

	instrs, leftOperand = b.leftOperand.TranslateToLLVMStack(function, block, program)
	leftType = leftOperand.GetType()
	if b.rightOperand == nil { return instrs, leftOperand }

	instrss, rightOperand = b.rightOperand.TranslateToLLVMStack(function, block, program)
	rightType = rightOperand.GetType()
	instrs = append(instrs, instrss...)

	if b.operator == ADD || b.operator == SUBTRACT || b.operator == MULTIPLY || b.operator == DIVIDE { exprType = types.IntTySig }
	entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, b.Token)
	funcRegisters[entry] = llvm.NewRegister(entry)

	switch b.operator {
		case ADD: operator = llvm.ADD
		case SUBTRACT: operator = llvm.SUB
		case MULTIPLY: operator = llvm.MUL
		case DIVIDE: operator = llvm.SDIV
		case OR: operator = llvm.OR
		case AND: operator = llvm.AND
		case EQ: operator = llvm.EQ
		case NEQ: operator = llvm.NE
		case LT: operator = llvm.SLT
		case GT: operator = llvm.SGT
		case LEQ: operator = llvm.SLE
		case GEQ: operator = llvm.SGE
	}

	if b.operator == ADD || b.operator == SUBTRACT || b.operator == MULTIPLY || b.operator == DIVIDE || b.operator == OR || b.operator == AND {
		instrs = append(instrs, llvm.NewBinaryOp(funcRegisters[entry], operator, leftType, leftOperand, rightOperand))
	} else if _, ok := leftType.(*types.PointerTy); rightType == types.NilTySig && ok {
		instrs = append(instrs, llvm.NewIcmp(funcRegisters[entry], operator, leftType, leftOperand, rightOperand))
	} else { instrs = append(instrs, llvm.NewIcmp(funcRegisters[entry], operator, rightType, leftOperand, rightOperand)) }

	operand = funcRegisters[entry]
	return instrs, operand
}

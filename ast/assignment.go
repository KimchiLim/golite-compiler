package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Assignment struct {
	*token.Token
	lvalue     *LValue
	expression Expression
}

func NewAssignment(lvalue *LValue, expression Expression, token *token.Token) *Assignment {
	return &Assignment{token, lvalue, expression}
}

func (a *Assignment) String() string {
	var out bytes.Buffer
	out.WriteString(a.lvalue.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(a.expression.String())
	out.WriteString(";")
	return out.String()
}

func (a *Assignment) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var leftType, rightType types.Type
	errors, leftType, block = a.lvalue.TypeCheck(errors, tables, function, block)
	errors, rightType, block = a.expression.TypeCheck(errors, tables, function, block)

	if pointerType, ok := leftType.(*types.PointerTy); ok {
		if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
			if structEntry := tables.Structs.Contains(structType.Name); structEntry != nil {
				if rightType == types.NilTySig {
					return errors, types.VoidTySig, block
				} else if exprType, ok := rightType.(*types.PointerTy); ok {
					if structEntry.Ty.String() == exprType.BaseType.String() {
						return errors, types.VoidTySig, block
					}
				}
			} else {
				errors = st.SemanticError(errors, a.Token.Line, a.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
				return errors, types.UnknownTySig, block
	    	}
		}
	} else if leftType == rightType { return errors, types.VoidTySig, block }
  
	errors = st.SemanticError(errors, a.Token.Line, a.Token.Column, fmt.Sprintf("invalid assignment: expected %s but got %s", leftType, rightType))
	return errors, types.UnknownTySig, block
}

func (a *Assignment) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs, instrss []llvm.LLVMInstruction
	var leftOperand, rightOperand llvm.LLVMOperand
	var exprType types.Type

	instrs, rightOperand = a.expression.TranslateToLLVMStack(function, block, program)
	instrss, leftOperand = a.lvalue.TranslateToLLVMStack(function, block, program)
	instrs = append(instrs, instrss...)
	exprType = leftOperand.GetType()

	if rightOperand.GetType() == types.NilTySig {
		instrs = append(instrs, llvm.NewStore(exprType, nil, leftOperand))
	} else {
		instrs = append(instrs, llvm.NewStore(exprType, rightOperand, leftOperand))
	}

	block.StackInstrs = append(block.StackInstrs, instrs...)
	return block
}

package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Invocation struct {
	*token.Token
	identifier *Variable
	arguments  []Expression
}

func NewInvocation(identifier *Variable, arguments []Expression, token *token.Token) *Invocation {
	return &Invocation{token, identifier, arguments}
}

func (i *Invocation) String() string {
	var out bytes.Buffer
	out.WriteString(i.identifier.String())
	out.WriteString("(")

	for k, expression := range i.arguments {
		if (k > 0) {out.WriteString(", ")}
		out.WriteString(expression.String())
	}

	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}

func (i *Invocation) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType, argType types.Type
	exprType = types.UnknownTySig
	funcName := i.identifier.String()

	if funcEntry := tables.Funcs.Contains(funcName); funcEntry != nil {
		exprType = types.VoidTySig

		if len(i.arguments) != len(funcEntry.Parameters) {
			errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("invalid invocation: expected %v arguments but got %v arguments", len(funcEntry.Parameters), len(i.arguments)))
			return errors, types.UnknownTySig, block
		}

		for k, argument := range i.arguments {
			errors, argType, block = argument.TypeCheck(errors, tables, function, block)

			if pointerType, ok := argType.(*types.PointerTy); ok {
				if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
					if structEntry := tables.Structs.Contains(structType.Name); structEntry != nil {
						if paramType, ok := funcEntry.Parameters[k].Ty.(*types.PointerTy); ok {
							if structEntry.Ty != paramType.BaseType {
								errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[k].Ty, argType))
								return errors, types.UnknownTySig, block
							}
						} else {
							errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[k].Ty, argType))
							return errors, types.UnknownTySig, block
						}
					} else {
						errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
						return errors, types.UnknownTySig, block		
					}
				}
			} else if argType != funcEntry.Parameters[k].Ty {
				errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[k].Ty, argType))
				return errors, types.UnknownTySig, block
			}
		}
	} else { errors = st.SemanticError(errors, i.Token.Line, i.Token.Column, fmt.Sprintf("undeclared funcEntry: %s", funcName)) }

	return errors, exprType, block
}

func (i *Invocation) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs, instrss []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var operands []llvm.LLVMOperand

	funcName := i.identifier.String()
	exprType = types.UnknownTySig

	if funcEntry := program.Tables.Funcs.Contains(funcName); funcEntry != nil {
		exprType = funcEntry.ReturnTy

		for _, argument := range i.arguments {
			instrss, operand = argument.TranslateToLLVMStack(function, block, program)
			instrs = append(instrs, instrss...)
			operands = append(operands, operand)
		}

		instrs = append(instrs, llvm.NewCall(nil, exprType, funcEntry, operands))
		block.StackInstrs = append(block.StackInstrs, instrs...)	
	}

	return block
}

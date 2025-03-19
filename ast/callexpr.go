package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type CallExpr struct {
	*token.Token
	identifier *Variable
	arguments  []Expression
}

func NewCallExpr(identifier *Variable, arguments []Expression, token *token.Token) *CallExpr {
	return &CallExpr{token, identifier, arguments}
}

func (c *CallExpr) String() string {
	var out bytes.Buffer
	out.WriteString(c.identifier.String())
	out.WriteString("(")

	for k, expression := range c.arguments {
		if (k > 0) {out.WriteString(", ")}
		out.WriteString(expression.String())
	}

	out.WriteString(")")
	return out.String()
}

func (c *CallExpr) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType, argType types.Type
	exprType = types.UnknownTySig
	funcName := c.identifier.String()

	if funcEntry := tables.Funcs.Contains(funcName); funcEntry != nil {
		exprType = funcEntry.ReturnTy

		if len(c.arguments) != len(funcEntry.Parameters) {
			errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("invalid call expression: expected %v arguments but got %v arguments", len(funcEntry.Parameters), len(c.arguments)))
			return errors, types.UnknownTySig, block
		}

		for i, argument := range c.arguments {
			errors, argType, block = argument.TypeCheck(errors, tables, function, block)

			if pointerType, ok := argType.(*types.PointerTy); ok {
				if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
					if structEntry := tables.Structs.Contains(structType.Name); structEntry != nil {
						if paramType, ok := funcEntry.Parameters[i].Ty.(*types.PointerTy); ok {
							if structEntry.Ty != paramType.BaseType {
								errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[i].Ty, argType))
								return errors, types.UnknownTySig, block
							}
						} else {
							errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[i].Ty, argType))
							return errors, types.UnknownTySig, block
						}
					} else {
						errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
						return errors, types.UnknownTySig, block	
					}
				}
			} else if argType != funcEntry.Parameters[i].Ty {
				errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("type mismatch: expected %s but got %s", funcEntry.Parameters[i].Ty, argType))
				return errors, types.UnknownTySig, block
			}
		}
	} else { errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("undeclared funcEntry: %s", funcName)) }

	return errors, exprType, block
}

func (c *CallExpr) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs, instrss []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	var operands []llvm.LLVMOperand
	funcRegisters := program.FuncDefs[function].Registers

	funcName := c.identifier.String()
	exprType = types.UnknownTySig

	if funcEntry := program.Tables.Funcs.Contains(funcName); funcEntry != nil {
		exprType = funcEntry.ReturnTy
		entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, c.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)

		for _, argument := range c.arguments {
			instrss, operand = argument.TranslateToLLVMStack(function, block, program)
			instrs = append(instrs, instrss...)
			operands = append(operands, operand)
		}

		instrs = append(instrs, llvm.NewCall(funcRegisters[entry], exprType, funcEntry, operands))
		operand = funcRegisters[entry]
	}

	return instrs, operand
}

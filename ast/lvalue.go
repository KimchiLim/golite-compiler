package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type LValue struct {
	*token.Token
	identifier *Variable
	fields     []*Variable
}

func NewLValue(identifier *Variable, fields []*Variable, token *token.Token) *LValue {
	return &LValue{token, identifier, fields}
}

func (l *LValue) String() string {
	var out bytes.Buffer
	out.WriteString(l.identifier.String())

	for _, field := range l.fields {
		out.WriteString(".")
		out.WriteString(field.String())
	}

	return out.String()
}

func (l *LValue) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	errors, exprType, block = l.identifier.TypeCheck(errors, tables, function, block)
	if len(l.fields) == 0 { return errors, exprType, block }

	for _, field := range l.fields {
		if pointerType, ok := exprType.(*types.PointerTy); ok { exprType = pointerType.BaseType }
		var fieldEntry *st.VarEntry

		if structType, ok := exprType.(*types.StructTy); !ok {
			errors = st.SemanticError(errors, l.Token.Line, l.Token.Column, fmt.Sprintf("invalid l-value: %s", l))
			return errors, types.UnknownTySig, block
		} else {
			if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
				errors = st.SemanticError(errors, l.Token.Line, l.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
				return errors, types.UnknownTySig, block

			} else {
				for _, structField := range structEntry.Fields {
					if structField.Name == field.String() { fieldEntry = structField }
				}

				if fieldEntry == nil {
					errors = st.SemanticError(errors, l.Token.Line, l.Token.Column, fmt.Sprintf("invalid field of %s: %s", structType.Name, field))
					return errors, types.UnknownTySig, block
				}
				
				if fieldEntry.Ty == types.IntTySig || fieldEntry.Ty == types.BoolTySig { exprType = fieldEntry.Ty }
				if _, ok := exprType.(*types.StructTy); ok { exprType = &types.PointerTy{fieldEntry.Ty} }
			}
		}
	}

	return errors, exprType, block
}

func (l *LValue) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	if len(l.fields) == 0 {
		instrs, operand = l.identifier.TranslateToLLVMStack(function, block, program)
		exprType = operand.GetType()
	} else {
		if varEntry := function.Variables.Contains(l.identifier.String()); varEntry != nil {
			exprType = varEntry.Ty
			entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, l.Token)
			funcRegisters[entry] = llvm.NewRegister(entry)
			instrs = append(instrs, llvm.NewLoad(funcRegisters[entry], exprType, funcRegisters[varEntry]))
			operand = funcRegisters[entry]
		}
	}

	for i, field := range l.fields {
		if pointerType, ok := exprType.(*types.PointerTy); ok { exprType = pointerType.BaseType }
		var fieldEntry *st.VarEntry
		var index int

		if structType, ok := exprType.(*types.StructTy); ok {
			if structEntry := program.Tables.Structs.Contains(structType.Name); structEntry != nil {
				for _, structField := range structEntry.Fields {
					if structField.Name == field.String() { fieldEntry = structField ; break } else { index++ }
				}

				entry = st.NewVarEntry(llvm.NewRegisterLabel(), fieldEntry.Ty, st.LOCAL, l.Token)
				funcRegisters[entry] = llvm.NewRegister(entry)
				instrs = append(instrs, llvm.NewGetElementPtr(funcRegisters[entry], structType, operand, index))
				operand = funcRegisters[entry]

				if fieldEntry.Ty == types.IntTySig || fieldEntry.Ty == types.BoolTySig { exprType = fieldEntry.Ty }
				if _, ok := exprType.(*types.StructTy); ok { exprType = &types.PointerTy{fieldEntry.Ty} }

				if i < len(l.fields) - 1 {
					entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, l.Token)
					funcRegisters[entry] = llvm.NewRegister(entry)
					instrs = append(instrs, llvm.NewLoad(funcRegisters[entry], exprType, operand))
					operand = funcRegisters[entry]	
				}
			}
		}
	}

	return instrs, operand
}

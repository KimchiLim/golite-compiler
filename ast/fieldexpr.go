package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type FieldExpr struct {
	*token.Token
	factor Expression
	fields []*Variable
}

func NewFieldExpr(factor Expression, fields [](*Variable), token *token.Token) *FieldExpr {
	return &FieldExpr{token, factor, fields}
}

func (f *FieldExpr) String() string {
	var out bytes.Buffer
	out.WriteString(f.factor.String())
    
	for _, field := range f.fields {
		out.WriteString(".")
		out.WriteString(field.String())
	}

	return out.String()
}

func (f *FieldExpr) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	errors, exprType, block = f.factor.TypeCheck(errors, tables, function, block)
	if len(f.fields) == 0 { return errors, exprType, block }

	for _, field := range f.fields {
		if pointerType, ok := exprType.(*types.PointerTy); ok { exprType = pointerType.BaseType }
		var fieldEntry *st.VarEntry

		if structType, ok := exprType.(*types.StructTy); !ok {
			errors = st.SemanticError(errors, f.Token.Line, f.Token.Column, fmt.Sprintf("invalid field expression: %s", f))
			return errors, types.UnknownTySig, block
		} else {
			if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
				errors = st.SemanticError(errors, f.Token.Line, f.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
				return errors, types.UnknownTySig, block

			} else {
				for _, structField := range structEntry.Fields {
					if structField.Name == field.String() { fieldEntry = structField }
				}

				if fieldEntry == nil {
					errors = st.SemanticError(errors, f.Token.Line, f.Token.Column, fmt.Sprintf("invalid field of %s: %s", structType.Name, field))
					return errors, types.UnknownTySig, block
				}
				
				if fieldEntry.Ty == types.IntTySig || fieldEntry.Ty == types.BoolTySig { exprType = fieldEntry.Ty }
				if _, ok := exprType.(*types.StructTy); ok { exprType = &types.PointerTy{fieldEntry.Ty} }
			}
		}
	}

	return errors, exprType, block
}


func (f *FieldExpr) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	if varEntry := function.Variables.Contains(f.factor.String()); varEntry != nil {
		exprType = varEntry.Ty
		entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, f.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)
		instrs = append(instrs, llvm.NewLoad(funcRegisters[entry], exprType, funcRegisters[varEntry]))
		operand = funcRegisters[entry]
	} else {
		instrs, operand = f.factor.TranslateToLLVMStack(function, block, program)
		exprType = operand.GetType()
	}

	for _, field := range f.fields {
		if pointerType, ok := exprType.(*types.PointerTy); ok { exprType = pointerType.BaseType }
		var fieldEntry *st.VarEntry
		var index int

		if structType, ok := exprType.(*types.StructTy); ok {
			if structEntry := program.Tables.Structs.Contains(structType.Name); structEntry != nil {
				for _, structField := range structEntry.Fields {
					if structField.Name == field.String() { fieldEntry = structField ; break } else { index++ }
				}

				entry = st.NewVarEntry(llvm.NewRegisterLabel(), fieldEntry.Ty, st.LOCAL, f.Token)
				funcRegisters[entry] = llvm.NewRegister(entry)
				instrs = append(instrs, llvm.NewGetElementPtr(funcRegisters[entry], structType, operand, index))
				operand = funcRegisters[entry]

				if fieldEntry.Ty == types.IntTySig || fieldEntry.Ty == types.BoolTySig { exprType = fieldEntry.Ty }
				if _, ok := exprType.(*types.StructTy); ok { exprType = &types.PointerTy{fieldEntry.Ty} }

				entry = st.NewVarEntry(llvm.NewRegisterLabel(), exprType, st.LOCAL, f.Token)
				funcRegisters[entry] = llvm.NewRegister(entry)
				instrs = append(instrs, llvm.NewLoad(funcRegisters[entry], exprType, operand))
				operand = funcRegisters[entry]
			}
		}
	}

	return instrs, operand
}

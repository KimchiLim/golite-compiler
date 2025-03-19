package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "fmt")

type Variable struct {
	*token.Token
	identifier string
}

func NewVariable(identifier string, token *token.Token) *Variable {
	return &Variable{token, identifier}
}

func (v *Variable) String() string {
	return v.identifier
}

func (v *Variable) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	exprType = types.UnknownTySig
	if varEntry := tables.Globals.Contains(v.identifier); varEntry != nil { exprType = varEntry.Ty }

	if function != nil {
		for _, paramEntry := range function.Parameters {
			if paramEntry.Name == v.identifier { return errors, paramEntry.Ty, block }
		}
		if varEntry := function.Variables.Contains(v.identifier); varEntry != nil { return errors, varEntry.Ty, block }
	}

	if exprType == types.IntTySig || exprType == types.BoolTySig { return errors, exprType, block }

	if pointerType, ok := exprType.(*types.PointerTy); ok {
		if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
			if structEntry := tables.Structs.Contains(structType.Name); structEntry != nil {
				return errors, structEntry.Ty, block
			} else {
				errors = st.SemanticError(errors, v.Token.Line, v.Token.Column, fmt.Sprintf("undeclared type: %s", structType.Name))
				return errors, types.UnknownTySig, block
			}
		}
	}

	errors = st.SemanticError(errors, v.Token.Line, v.Token.Column, fmt.Sprintf("undeclared variable: %s", v.identifier))
	return errors, types.UnknownTySig, block
}

func (v *Variable) TranslateToLLVMStack(function *st.FuncEntry, block *llvm.Block, program *llvm.LLVMProgram) ([]llvm.LLVMInstruction, llvm.LLVMOperand) {
	var instrs []llvm.LLVMInstruction
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers
	if varEntry := program.Tables.Globals.Contains(v.identifier); varEntry != nil { entry = varEntry }

	if function != nil {
		for _, paramEntry := range function.Parameters {
			if paramEntry.Name == v.identifier { return instrs, funcRegisters[paramEntry] }
		}
		if varEntry := function.Variables.Contains(v.identifier); varEntry != nil { return instrs, funcRegisters[varEntry] }
	}

	return instrs, funcRegisters[entry]
}

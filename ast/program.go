package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Program struct {
	*token.Token
	Types        []*TypeDeclaration
	Declarations []*Declaration
	Functions    []*Function
}

func NewProgram(types []*TypeDeclaration, declarations []*Declaration, functions []*Function, token *token.Token) *Program {
	return &Program{token, types, declarations, functions}
}

func (p *Program) String() string {
	var out bytes.Buffer
	out.WriteString("\n")

	for i, typeDeclaration := range p.Types {
		out.WriteString(typeDeclaration.String())
		if i < len(p.Types) - 1 {out.WriteString("\n")}
	}

	if len(p.Types) > 0 && len(p.Declarations) > 0 { out.WriteString("\n") }

	for _, declaration := range p.Declarations { out.WriteString(declaration.String()) }
	
	if len(p.Types) > 0 && len(p.Declarations) == 0 && len(p.Functions) > 0 { out.WriteString("\n") }
	if len(p.Declarations) > 0 && len(p.Functions) > 0 { out.WriteString("\n") }

	for i, function := range p.Functions {
		out.WriteString(function.String())
		if i < len(p.Functions) - 1 {out.WriteString("\n")}
	}

	return out.String()
}

func (p *Program) BuildSymbolTable(errors []*context.CompilerError, tables *st.SymbolTables) []*context.CompilerError {

	for _, typeDeclaration := range p.Types {
		structName := typeDeclaration.identifier.String()
		structType := &types.StructTy{structName}
		var fields []*st.VarEntry

		if tables.Structs.Contains(structName) != nil {
			errors = st.SemanticError(errors, typeDeclaration.Token.Line, typeDeclaration.Token.Column,
				fmt.Sprintf("two struct definitions with same name cannot be defined: %s", structName))
		} else {
			for _, decl := range typeDeclaration.fields {
				fieldName := decl.identifier.String()
				var isDefined bool

				for _, field := range fields {
					if field.Name == fieldName {
						isDefined = true
						errors = st.SemanticError(errors, decl.Token.Line, decl.Token.Column,
							fmt.Sprintf("two fields within a struct cannot be defined: %s", fieldName))
					}
				}

				if !isDefined {
					if pointerType, ok := decl.ty.(*types.PointerTy); ok {
						if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
							if structType.Name == structName {
								decl.ty = structType
							} else {
								if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
									errors = st.SemanticError(errors, decl.Token.Line, decl.Token.Column,
										fmt.Sprintf("undeclared type: %s", structType.Name))
								} else { decl.ty = structEntry.Ty }
							}
						}
					}
					fields = append(fields, st.NewVarEntry(fieldName, decl.ty, st.GLOBAL, decl.Token))
				}
			}
			tables.Structs.Insert(structName, st.NewStructEntry(structName, structType, fields, typeDeclaration.Token))
		}
	}

	for _, declaration := range p.Declarations {
		for _, identifier := range declaration.identifiers {
			varName := identifier.String()
	
			if tables.Globals.Contains(varName) != nil {
				errors = st.SemanticError(errors, identifier.Token.Line, identifier.Token.Column,
					fmt.Sprintf("global variables with the same name cannot be defined: %s", varName))
			} else {
				if pointerType, ok := declaration.ty.(*types.PointerTy); ok {
					if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
						if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
							errors = st.SemanticError(errors, declaration.Token.Line, declaration.Token.Column,
								fmt.Sprintf("undeclared type: %s", structType.Name))
						} else { declaration.ty = &types.PointerTy{structEntry.Ty} }
					}
				}
				tables.Globals.Insert(varName, st.NewVarEntry(varName, declaration.ty, st.GLOBAL, declaration.Token))
			}
		}
	}

	for _, function := range p.Functions {
		funcName := function.identifier.String()
		var parameters []*st.VarEntry
		variables := st.NewSymbolTable[*st.VarEntry](nil)

		if tables.Funcs.Contains(funcName) != nil {
			errors = st.SemanticError(errors, function.Token.Line, function.Token.Column,
				fmt.Sprintf("two functions with the same name cannot be defined: %s", funcName))
		} else {
			for _, decl := range function.parameters {
				paramName := decl.identifier.String()
				var isDefined bool

				for _, parameter := range parameters {
					if parameter.Name == paramName {
						isDefined = true
						errors = st.SemanticError(errors, decl.Token.Line, decl.Token.Column,
							fmt.Sprintf("two parameters within a function cannot be defined: %s", paramName))
					}
				}

				if !isDefined {
					if pointerType, ok := decl.ty.(*types.PointerTy); ok {
						if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
							if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
								errors = st.SemanticError(errors, decl.Token.Line, decl.Token.Column,
									fmt.Sprintf("undeclared type: %s", structType.Name))
							} else { decl.ty = &types.PointerTy{structEntry.Ty} }
						}
					}
					parameters = append(parameters, st.NewVarEntry(paramName, decl.ty, st.LOCAL, decl.Token))
				}
			}

			for _, declaration := range function.declarations {
				for _, identifier := range declaration.identifiers {
					varName := identifier.String()
					var isDefined bool

					for _, parameter := range parameters {
						if parameter.Name == varName {
							isDefined = true
							errors = st.SemanticError(errors, identifier.Token.Line, identifier.Token.Column,
								fmt.Sprintf("a local declaration may not redefine a parameter: %s", varName))
						}
					}

					if !isDefined {
						if variables.Contains(varName) != nil {
							errors = st.SemanticError(errors, identifier.Token.Line, identifier.Token.Column,
								fmt.Sprintf("local variables of the same function cannot be defined: %s", varName))
						} else {
							if pointerType, ok := declaration.ty.(*types.PointerTy); ok {
								if structType, ok := pointerType.BaseType.(*types.StructTy); ok {
									if structEntry := tables.Structs.Contains(structType.Name); structEntry == nil {
										errors = st.SemanticError(errors, declaration.Token.Line, declaration.Token.Column,
											fmt.Sprintf("undeclared type: %s", structType.Name))
									} else { declaration.ty = &types.PointerTy{structEntry.Ty} }
								}
							}
							variables.Insert(varName, st.NewVarEntry(varName, declaration.ty, st.LOCAL, declaration.Token))
						}
					}
				}
			}

			variables.Parent = tables.Globals
			if function.returnType == nil { function.returnType = types.VoidTySig }
			tables.Funcs.Insert(funcName, st.NewFuncEntry(funcName, function.returnType, parameters, variables, function.Token))
		}
	}

	if funcEntry := tables.Funcs.Contains("main"); funcEntry != nil {
		if len(funcEntry.Parameters) > 0 {
			errors = st.SemanticError(errors, 0, 0, fmt.Sprintf("main function: expected 0 arguments but got %v arguments", len(funcEntry.Parameters)))
		}
		if funcEntry.ReturnTy != types.VoidTySig {
			errors = st.SemanticError(errors, 0, 0, fmt.Sprintf("main function: expected void but got %s return type", funcEntry.ReturnTy))
		}
	} else if funcEntry == nil {
		errors = st.SemanticError(errors, 0, 0, fmt.Sprintf("invalid program: a main function must be defined"))
	}

	return errors
}

func (p *Program) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, blocks []*cfg.Block) []*context.CompilerError {
	var root, curr *cfg.Block
	var returnType types.Type
	var valid bool

	for _, function := range p.Functions {
		funcName := function.identifier.String()

		root = cfg.NewBlock(function.Token)
		blocks = append(blocks, root)
		curr = root

		if funcEntry := tables.Funcs.Contains(funcName); funcEntry != nil {
			for _, statement := range function.statements { errors, _, curr = statement.TypeCheck(errors, tables, funcEntry, curr) }

			returnType = funcEntry.ReturnTy
			visited := make(map[string]bool)
			errors, valid = root.Validate(errors, returnType, visited)

			if !valid && returnType != types.VoidTySig {
				errors = st.SemanticError(errors, function.Token.Line, function.Token.Column, fmt.Sprintf("not all control paths have a return statement: %s", funcName))
			}
		} else { errors = st.SemanticError(errors, function.Token.Line, function.Token.Column, fmt.Sprintf("undeclared function: %s", funcName)) }	
	}

	return errors
}

func (p *Program) TranslateToLLVMStack(source string, target string, tables *st.SymbolTables) *llvm.LLVMProgram {
	var llvmProgram *llvm.LLVMProgram
	var root, curr *llvm.Block
	llvmProgram = llvm.NewLLVMProgram(source, target, tables)

	for _, function := range p.Functions {
		funcName := function.identifier.String()
		root = llvm.NewBlock()
		curr = root

		if funcEntry := tables.Funcs.Contains(funcName); funcEntry != nil {
			llvmProgram.FuncDefs[funcEntry] = llvm.NewLLVMFunction(funcEntry, root)
			funcRegisters := llvmProgram.FuncDefs[funcEntry].Registers
			for _, varEntry := range llvmProgram.Tables.Globals.Table { funcRegisters[varEntry] = llvm.NewRegister(varEntry) }

			if funcName == "main" {
				returnEntry := st.NewVarEntry("_ret_val", types.IntTySig, st.LOCAL, funcEntry.Token)
				funcEntry.Variables.Insert(returnEntry.Name, returnEntry)
				funcRegisters[returnEntry] = llvm.NewRegister(returnEntry)
				root.StackInstrs = append(root.StackInstrs, llvm.NewAlloca(funcRegisters[returnEntry], types.IntTySig))

			} else if funcEntry.ReturnTy != types.VoidTySig {
				returnEntry := st.NewVarEntry("_ret_val", funcEntry.ReturnTy, st.LOCAL, funcEntry.Token)
				funcEntry.Variables.Insert(returnEntry.Name, returnEntry)
				funcRegisters[returnEntry] = llvm.NewRegister(returnEntry)
				root.StackInstrs = append(root.StackInstrs, llvm.NewAlloca(funcRegisters[returnEntry], funcEntry.ReturnTy))
			}

			for _, decl := range function.parameters {
				varEntry := st.NewVarEntry("_P_" + decl.identifier.String(), decl.ty, st.LOCAL, decl.Token)
				var paramEntry *st.VarEntry

				for _, parameter := range funcEntry.Parameters {
					if parameter.Name == decl.identifier.String() { paramEntry = parameter }
				}

				funcEntry.Variables.Insert(varEntry.Name, varEntry)
				funcEntry.Variables.Insert(paramEntry.Name, paramEntry)
				funcRegisters[varEntry] = llvm.NewRegister(varEntry)
				funcRegisters[paramEntry] = funcRegisters[varEntry]
				root.StackInstrs = append(root.StackInstrs, llvm.NewAlloca(funcRegisters[varEntry], varEntry.Ty))
				root.StackInstrs = append(root.StackInstrs, llvm.NewStore(varEntry.Ty, llvm.NewRegister(paramEntry), funcRegisters[varEntry]))
			}

			for _, declaration := range function.declarations {
				for _, identifier := range declaration.identifiers {
					if varEntry := funcEntry.Variables.Contains(identifier.String()); varEntry != nil {
						funcRegisters[varEntry] = llvm.NewRegister(varEntry)
						root.StackInstrs = append(root.StackInstrs, llvm.NewAlloca(funcRegisters[varEntry], varEntry.Ty))
					}
				}
			}

			for _, statement := range function.statements { curr = statement.TranslateToLLVMStack(curr, funcEntry, llvmProgram) }

			if funcEntry.Name == "main" && !curr.HasReturn {
				curr.StackInstrs = append(curr.StackInstrs, llvm.NewReturn(types.IntTySig, llvm.NewImmediate(0, types.IntTySig)))
			} else if !curr.HasReturn && funcEntry.ReturnTy != types.VoidTySig {
				if varEntry := funcEntry.Variables.Contains("_ret_val"); varEntry != nil {
					entry := st.NewVarEntry(llvm.NewRegisterLabel(), varEntry.Ty, st.LOCAL, function.Token)
					funcRegisters[entry] = llvm.NewRegister(entry)
					curr.StackInstrs = append(curr.StackInstrs, llvm.NewLoad(funcRegisters[entry], varEntry.Ty, funcRegisters[varEntry]))
					curr.StackInstrs = append(curr.StackInstrs, llvm.NewReturn(varEntry.Ty, funcRegisters[entry]))
				}
			} else if !curr.HasReturn && funcEntry.ReturnTy == types.VoidTySig { curr.StackInstrs = append(curr.StackInstrs, llvm.NewReturn(types.VoidTySig, nil)) }
		}
	}

	return llvmProgram
}

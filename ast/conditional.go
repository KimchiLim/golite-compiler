package ast

import ("golite/cfg" ; "golite/context" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt" ; "golite/llvm")

type Conditional struct {
	*token.Token
	expression Expression
	ifBlock    []Statement
	elseBlock  []Statement
}

func NewConditional(expression Expression, ifBlock []Statement, elseBlock []Statement, token *token.Token) *Conditional {
	return &Conditional{token, expression, ifBlock, elseBlock}
}

func (c *Conditional) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(" ")
	out.WriteString(c.expression.String())
	out.WriteString(" ")
	out.WriteString("{")
	out.WriteString("\n")

	for _, statement := range c.ifBlock {
		out.WriteString("\t")
		out.WriteString(statement.String())
		out.WriteString("\n")
	}

	if (len(c.elseBlock) > 0) {
		out.WriteString("    ")
		out.WriteString("}")
		out.WriteString(" else ")
		out.WriteString("{")
		out.WriteString("\n")

		for _, statement := range c.elseBlock {
			out.WriteString("\t")
			out.WriteString(statement.String())
			out.WriteString("\n")
		}
	}
    
	out.WriteString("    ")
	out.WriteString("}")
	return out.String()
}

func (c *Conditional) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	var ifBlock, elseBlock, exitBlock *cfg.Block
	errors, exprType, block = c.expression.TypeCheck(errors, tables, function, block)

	if exprType != types.BoolTySig {
		errors = st.SemanticError(errors, c.Token.Line, c.Token.Column, fmt.Sprintf("invalid conditional expression: %s", c.expression))
		return errors, types.UnknownTySig, block
	}

	ifBlock = cfg.NewBlock(c.Token)
	elseBlock = cfg.NewBlock(c.Token)
	block.AddNext(ifBlock)
	block.AddNext(elseBlock)

	for _, statement := range c.ifBlock { errors, _, ifBlock = statement.TypeCheck(errors, tables, function, ifBlock) }
	for _, statement := range c.elseBlock { errors, _, elseBlock = statement.TypeCheck(errors, tables, function, elseBlock) }

	exitBlock = cfg.NewBlock(c.Token)
	ifBlock.AddNext(exitBlock)
	elseBlock.AddNext(exitBlock)
	return errors, types.VoidTySig, exitBlock
}

func (c *Conditional) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	ifBlock := llvm.NewBlock()
	elseBlock := llvm.NewBlock()
	exitBlock := llvm.NewBlock()
	block.AddNext(ifBlock)
	block.AddNext(elseBlock)

	instrs, operand = c.expression.TranslateToLLVMStack(function, block, program)
	exprType = operand.GetType()

	if exprType != types.Int1TySig {
		entry = st.NewVarEntry(llvm.NewRegisterLabel(), types.Int1TySig, st.LOCAL, c.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)
		instrs = append(instrs, llvm.NewIcmp(funcRegisters[entry], llvm.EQ, types.BoolTySig, operand, llvm.NewImmediate(1, types.BoolTySig)))
		operand = funcRegisters[entry]
	}

	instrs = append(instrs, llvm.NewBranch(nil, operand, ifBlock, elseBlock))
	block.StackInstrs = append(block.StackInstrs, instrs...)

	for _, statement := range c.ifBlock { ifBlock = statement.TranslateToLLVMStack(ifBlock, function, program) }
	for _, statement := range c.elseBlock { elseBlock = statement.TranslateToLLVMStack(elseBlock, function, program) }

	if !ifBlock.HasReturn {
		ifBlock.StackInstrs = append(ifBlock.StackInstrs, llvm.NewBranch(exitBlock, nil, nil, nil))
		ifBlock.AddNext(exitBlock)
	}

	if !elseBlock.HasReturn {
		elseBlock.StackInstrs = append(elseBlock.StackInstrs, llvm.NewBranch(exitBlock, nil, nil, nil))
		elseBlock.AddNext(exitBlock)
	}

	return exitBlock
}

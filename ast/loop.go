package ast

import ("golite/cfg" ; "golite/context" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt" ; "golite/llvm")

type Loop struct {
	*token.Token
	expression Expression
	forBlock   []Statement
}

func NewLoop(expression Expression, forBlock []Statement, token *token.Token) *Loop {
	return &Loop{token, expression, forBlock}
}

func (l *Loop) String() string {
	var out bytes.Buffer
	out.WriteString("for")
	out.WriteString(" ")
	out.WriteString(l.expression.String())
	out.WriteString(" ")
	out.WriteString("{")
	out.WriteString("\n")

	for _, statement := range l.forBlock {
		out.WriteString("\t")
		out.WriteString(statement.String())
		out.WriteString("\n")
    }

	out.WriteString("    ")
	out.WriteString("}")
	return out.String()
}

func (l *Loop) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	var bodyBlock, exitBlock *cfg.Block
	errors, exprType, block = l.expression.TypeCheck(errors, tables, function, block)

	if exprType != types.BoolTySig {
		errors = st.SemanticError(errors, l.Token.Line, l.Token.Column, fmt.Sprintf("invalid loop expression: %s", l.expression))
		return errors, types.UnknownTySig, block
	}

	bodyBlock = cfg.NewBlock(l.Token)
	block.AddNext(bodyBlock)

	for _, statement := range l.forBlock { errors, _, bodyBlock = statement.TypeCheck(errors, tables, function, bodyBlock) }
	
	bodyBlock.AddNext(block)
	exitBlock = cfg.NewBlock(l.Token)
	block.AddNext(exitBlock)
	return errors, types.VoidTySig, exitBlock
}

func (l *Loop) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var exprType types.Type
	var entry *st.VarEntry
	funcRegisters := program.FuncDefs[function].Registers

	condBlock := llvm.NewBlock()
	bodyBlock := llvm.NewBlock()
	exitBlock := llvm.NewBlock()
	block.AddNext(condBlock)
	condBlock.AddNext(bodyBlock)
	condBlock.AddNext(exitBlock)

	instrs, operand = l.expression.TranslateToLLVMStack(function, block, program)
	exprType = operand.GetType()

	if exprType != types.Int1TySig {
		entry = st.NewVarEntry(llvm.NewRegisterLabel(), types.Int1TySig, st.LOCAL, l.Token)
		funcRegisters[entry] = llvm.NewRegister(entry)
		instrs = append(instrs, llvm.NewIcmp(funcRegisters[entry], llvm.EQ, types.BoolTySig, operand, llvm.NewImmediate(1, types.BoolTySig)))
		operand = funcRegisters[entry]
	}

	block.StackInstrs = append(block.StackInstrs, llvm.NewBranch(condBlock, nil, nil, nil))
	condBlock.StackInstrs = append(instrs, llvm.NewBranch(nil, operand, bodyBlock, exitBlock))

	for _, statement := range l.forBlock { bodyBlock = statement.TranslateToLLVMStack(bodyBlock, function, program) }

	if !bodyBlock.HasReturn {
		bodyBlock.StackInstrs = append(bodyBlock.StackInstrs, llvm.NewBranch(condBlock, nil, nil, nil))
		bodyBlock.AddNext(condBlock)
	}
	return exitBlock
}

package ast

import ("golite/cfg" ; "golite/context" ; "golite/llvm" ; "golite/st" ; "golite/token" ; "golite/types" ; "bytes" ; "fmt")

type Print struct {
	*token.Token
	strLiteral  string
	expressions []Expression
}

func NewPrint(strLiteral string, expressions []Expression, token *token.Token) *Print {
	return &Print{token, strLiteral, expressions}
}

func (p *Print) String() string {
	var out bytes.Buffer
	out.WriteString("printf")
	out.WriteString("(")
	out.WriteString(p.strLiteral)

	for _, expression := range p.expressions {
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString(expression.String())
	}

	out.WriteString(")")
	out.WriteString(";")
	return out.String()
}

func (p *Print) TypeCheck(errors []*context.CompilerError, tables *st.SymbolTables, function *st.FuncEntry, block *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block) {
	var exprType types.Type
	exprType = types.UnknownTySig
	placeholders := 0

	for i := 0; i < len(p.strLiteral); i++ {
		if p.strLiteral[i] == '%' {
			if i + 1 < len(p.strLiteral) && p.strLiteral[i + 1] == 'd' { placeholders++ }
		}
	}

	if len(p.expressions) != placeholders {
		errors = st.SemanticError(errors, p.Token.Line, p.Token.Column, fmt.Sprintf("invalid print: expected %v arguments but got %v arguments", placeholders, len(p.expressions)))
		return errors, types.UnknownTySig, block
	}

	for _, expression := range p.expressions {
		errors, exprType, block = expression.TypeCheck(errors, tables, function, block)

		if exprType != types.IntTySig {
			errors = st.SemanticError(errors, p.Token.Line, p.Token.Column, fmt.Sprintf("type mismatch: expected int but got %s", exprType))
			return errors, types.UnknownTySig, block
		}
	}

	return errors, types.VoidTySig, block
}

func (p *Print) TranslateToLLVMStack(block *llvm.Block, function *st.FuncEntry, program *llvm.LLVMProgram) *llvm.Block {
	var instrs, instrss []llvm.LLVMInstruction
	var operand llvm.LLVMOperand
	var fstring *llvm.FString
	var operands []llvm.LLVMOperand

	fstring = llvm.NewFString(p.strLiteral[1 : len(p.strLiteral) - 1])
	program.FStrings = append(program.FStrings, fstring)

	for _, expression := range p.expressions {
		instrss, operand = expression.TranslateToLLVMStack(function, block, program)
		instrs = append(instrs, instrss...)
		operands = append(operands, operand)
	}

	instrs = append(instrs, llvm.NewPrint(fstring, operands))
	block.StackInstrs = append(block.StackInstrs, instrs...)
	return block
}

package llvm

import ("golite/codegen" ; "bytes")

type Move struct {
	op   LLVMOperand
	dest LLVMOperand
}

func NewMove(op LLVMOperand, dest LLVMOperand) *Move {
	return &Move{op, dest}
}

func (m *Move) String() string {
	var out bytes.Buffer
	out.WriteString("mov")
	out.WriteString(" ")
	out.WriteString(m.op.String())
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(m.dest.String())
	out.WriteString("\n")
	return out.String()
}

func (m *Move) MemtoReg(block *Block) {}

func (m *Move) ComputeLiveRange(function *LLVMFunction, position int) {}

func (m *Move) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {}

func (m *Move) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instr codegen.Instruction
	instr = codegen.NewMove(m.dest.GetValue(), m.op.GetValue())
	block.AssemblyInstrs = append(block.AssemblyInstrs, instr)
}

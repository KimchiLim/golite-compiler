package llvm

import ("golite/codegen" ; "golite/types" ; "bytes" ; "fmt")

type Branch struct {
	dest    *Block
	cond    LLVMOperand
	ifTrue  *Block
	ifFalse *Block
}

func NewBranch(dest *Block, cond LLVMOperand, ifTrue *Block, ifFalse *Block) *Branch {
	return &Branch{dest, cond, ifTrue, ifFalse}
}

func (b *Branch) String() string {
	var out bytes.Buffer
	out.WriteString("br")
	out.WriteString(" ")

	if b.dest != nil {
		out.WriteString("label")
		out.WriteString(" ")
		out.WriteString(b.dest.String())
	} else {
		out.WriteString("i1")
		out.WriteString(" ")
		out.WriteString(b.cond.String())
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString("label")
		out.WriteString(" ")
		out.WriteString(b.ifTrue.String())
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString("label")
		out.WriteString(" ")
		out.WriteString(b.ifFalse.String())
	}

	out.WriteString("\n")
	return out.String()
}

func (b *Branch) MemtoReg(block *Block) {
	if b.cond != nil { if value, exists := block.DefsMap[b.cond.String()]; exists { b.cond = value } }
	block.RegisterInstrs = append(block.RegisterInstrs, b)
}

func (b *Branch) ComputeLiveRange(function *LLVMFunction, position int) {
	if b.cond != nil {
		if _, ok := b.cond.(*LLVMRegister); ok {
			if value, exists := function.Allocation[b.cond.String()]; !exists {
				function.Allocation[b.cond.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}
}

func (b *Branch) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if b.cond != nil {
		if value, exists := function.Allocation[b.cond.String()]; exists {
			if value.Spilled {
				b.cond = NewPhysicalRegister("x9")
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(b.cond, types.Int1TySig, spillRegister))
			} else { b.cond = NewPhysicalRegister(value.Register) }
		}
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, b)
}

func (b *Branch) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	if b.cond != nil {
		instrs = append(instrs, codegen.NewBranch(b.cond.GetValue(), fmt.Sprintf(".%s", b.ifFalse.Label)))
		instrs = append(instrs, codegen.NewBranch("", fmt.Sprintf(".%s", b.ifTrue.Label)))
	} else {
		instrs = append(instrs, codegen.NewBranch("", fmt.Sprintf(".%s", b.dest.Label)))
	}

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

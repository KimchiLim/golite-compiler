package llvm

import ("golite/codegen" ; "golite/types" ; "bytes" ; "fmt")

type Return struct {
	ty types.Type
	op LLVMOperand
}

func NewReturn(ty types.Type, op LLVMOperand) *Return {
	return &Return{ty, op}
}

func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString("ret")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(r.ty))

	if r.ty != types.VoidTySig {
		out.WriteString(" ")
		out.WriteString(r.op.String())
	}

	out.WriteString("\n")
	return out.String()
}

func (r *Return) MemtoReg(block *Block) {
	if r.ty != types.VoidTySig {
		if value, exists := block.DefsMap[r.op.String()]; exists { r.op = value }
	}

	block.RegisterInstrs = append(block.RegisterInstrs, r)
}

func (r *Return) ComputeLiveRange(function *LLVMFunction, position int) {
	if r.ty != types.VoidTySig {
		if _, ok := r.op.(*LLVMRegister); ok {
			if value, exists := function.Allocation[r.op.String()]; !exists {
				function.Allocation[r.op.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}
}

func (r *Return) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if r.ty != types.VoidTySig {
		if value, exists := function.Allocation[r.op.String()]; exists {
			if value.Spilled {
				r.op = NewPhysicalRegister("x9")
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(r.op, r.ty, spillRegister))
			} else { r.op = NewPhysicalRegister(value.Register) }
		}
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, r)
}

func (r *Return) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction
	var offset int

	offset = len(function.Allocation) * 8 + 8 + 8
	offset = ((offset + 15) / 16) * 16

	if r.ty != types.VoidTySig { instrs = append(instrs, codegen.NewMove("x0", r.op.GetValue())) }

	instrs = append(instrs, codegen.NewLdp("x29", "x30", "sp"))
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", fmt.Sprintf("#%d", offset)))
	instrs = append(instrs, codegen.NewReturn(""))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

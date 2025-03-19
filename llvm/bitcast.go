package llvm

import ("golite/codegen" ; "golite/types" ; "bytes")

type Bitcast struct {
	result LLVMOperand
	ty     types.Type
	value  LLVMOperand
	ty2    types.Type
}

func NewBitcast(result LLVMOperand, ty types.Type, value LLVMOperand, ty2 types.Type) *Bitcast {
	return &Bitcast{result, ty, value, ty2}
}

func (b *Bitcast) String() string {
	var out bytes.Buffer
	out.WriteString(b.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("bitcast")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(b.ty))
	out.WriteString(" ")
	out.WriteString(b.value.String())
	out.WriteString(" ")
	out.WriteString("to")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(b.ty2))
	out.WriteString("\n")
	return out.String()
}

func (b *Bitcast) MemtoReg(block *Block) {
	if value, exists := block.DefsMap[b.value.String()]; exists { b.value = value }

	if _, exists := block.DefsMap[b.result.String()]; exists {
		block.DefsMap[b.result.String()] = b.result
	}

	block.RegisterInstrs = append(block.RegisterInstrs, b)
}

func (b *Bitcast) ComputeLiveRange(function *LLVMFunction, position int) {
	if value, exists := function.Allocation[b.value.String()]; !exists {
		function.Allocation[b.value.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }

	if value, exists := function.Allocation[b.result.String()]; !exists {
		function.Allocation[b.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (b *Bitcast) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[b.value.String()]; exists {
		if value.Spilled {
			b.value = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(b.value, b.ty, spillRegister))
		} else { b.value = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[b.result.String()]; exists {
		if value.Spilled {
			b.result = NewPhysicalRegister("x9")
		} else { b.result = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, b)

	if value, exists := function.Allocation[b.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(b.ty2, b.result, spillRegister))
		}
	}
}

func (b *Bitcast) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instr codegen.Instruction
	instr = codegen.NewMove(b.result.GetValue(), b.value.GetValue())
	block.AssemblyInstrs = append(block.AssemblyInstrs, instr)
}

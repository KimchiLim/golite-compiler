package llvm

import ("golite/codegen" ; "golite/types" ; "bytes" ; "fmt" )

type Load struct {
	result  LLVMOperand
	ty      types.Type
	pointer LLVMOperand
}

func NewLoad(result LLVMOperand, ty types.Type, pointer LLVMOperand) *Load {
	return &Load{result, ty, pointer}
}

func (l *Load) String() string {
	var out bytes.Buffer
	out.WriteString(l.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("load")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(l.ty))
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(l.ty) + "*")
	out.WriteString(" ")
	out.WriteString(l.pointer.String())
	out.WriteString("\n")
	return out.String()
}

func (l *Load) MemtoReg(block *Block) {
	if value, exists := block.DefsMap[l.pointer.String()]; exists {
		if value != nil { l.pointer = value } else { l.pointer = NewImmediate(0, l.ty) }
	}

	if l.pointer.PointsToStack() {
		block.DefsMap[l.result.String()] = l.pointer
	} else {
		block.DefsMap[l.result.String()] = l.result
		block.RegisterInstrs = append(block.RegisterInstrs, l)
	}
}

func (l *Load) ComputeLiveRange(function *LLVMFunction, position int) {
	if _, ok := l.pointer.(*LLVMRegister); ok {
		if l.pointer.String()[0] != '@' {
			if value, exists := function.Allocation[l.pointer.String()]; !exists {
				function.Allocation[l.pointer.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}

	if value, exists := function.Allocation[l.result.String()]; !exists {
		function.Allocation[l.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (l *Load) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[l.pointer.String()]; exists {
		if value.Spilled {
			l.pointer = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(l.pointer, l.ty, spillRegister))
		} else { l.pointer = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[l.result.String()]; exists {
		if value.Spilled {
			l.result = NewPhysicalRegister("x9")
		} else { l.result = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, l)

	if value, exists := function.Allocation[l.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(l.ty, l.result, spillRegister))
		}
	}
}

func (l *Load) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	if l.pointer.String()[0] == '@' {
		instrs = append(instrs, codegen.NewAdrp("x10", l.pointer.GetValue()))
		instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "x10", "x10", fmt.Sprintf(":lo12:%s", l.pointer.GetValue())))
		instrs = append(instrs, codegen.NewLoad(l.result.GetValue(), "x10"))
	} else {
		instrs = append(instrs, codegen.NewLoad(l.result.GetValue(), l.pointer.GetValue()))
	}

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

package llvm

import ("golite/codegen" ; "bytes")

type Free struct {
	pointer LLVMOperand
}

func NewFree(pointer LLVMOperand) *Free {
	return &Free{pointer}
}

func (f *Free) String() string {
	var out bytes.Buffer
	out.WriteString("call")
	out.WriteString(" ")
	out.WriteString("void")
	out.WriteString(" ")
	out.WriteString("@free")
	out.WriteString("(")
	out.WriteString("i8*")
	out.WriteString(" ")
	out.WriteString(f.pointer.String())
	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}

func (f *Free) MemtoReg(block *Block) {
	delete(block.DefsMap, f.pointer.String())
	block.RegisterInstrs = append(block.RegisterInstrs, f)
}

func (f *Free) ComputeLiveRange(function *LLVMFunction, position int) {
	if value, exists := function.Allocation[f.pointer.String()]; !exists {
		function.Allocation[f.pointer.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (f *Free) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	ty := f.pointer.GetType()

	if value, exists := function.Allocation[f.pointer.String()]; exists {
		if value.Spilled {
			f.pointer = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(f.pointer, ty, spillRegister))
		} else { f.pointer = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, f)
}

func (f *Free) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", "#16"))
	instrs = append(instrs, codegen.NewStore("x0", "sp, #0"))

	instrs = append(instrs, codegen.NewMove("x0", f.pointer.GetValue()))
	instrs = append(instrs, codegen.NewCall("free"))

	instrs = append(instrs, codegen.NewLoad("x0", "sp, #0"))
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", "#16"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}
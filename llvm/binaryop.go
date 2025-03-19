package llvm

import ("golite/codegen" ; "golite/types" ; "bytes")

type BinaryOp struct {
	result   LLVMOperand
	operator LLVMOperator
	ty       types.Type
	op1      LLVMOperand
	op2      LLVMOperand
}

func NewBinaryOp(result LLVMOperand, operator LLVMOperator, ty types.Type, op1 LLVMOperand, op2 LLVMOperand) *BinaryOp {
	return &BinaryOp{result, operator, ty, op1, op2}
}

func (b *BinaryOp) String() string {
	var out bytes.Buffer
	out.WriteString(b.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString(b.operator.String())
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(b.ty))
	out.WriteString(" ")
	out.WriteString(b.op1.String())
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(b.op2.String())
	out.WriteString("\n")
	return out.String()
}

func (b *BinaryOp) MemtoReg(block *Block) {
	if value, exists := block.DefsMap[b.op1.String()]; exists { b.op1 = value }
	if value, exists := block.DefsMap[b.op2.String()]; exists { b.op2 = value }

	if _, exists := block.DefsMap[b.result.String()]; exists {
		block.DefsMap[b.result.String()] = b.result
	}

	block.RegisterInstrs = append(block.RegisterInstrs, b)
}

func (b *BinaryOp) ComputeLiveRange(function *LLVMFunction, position int) {
	if _, ok := b.op1.(*LLVMRegister); ok {
		if value, exists := function.Allocation[b.op1.String()]; !exists {
			function.Allocation[b.op1.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}

	if _, ok := b.op2.(*LLVMRegister); ok {
		if value, exists := function.Allocation[b.op2.String()]; !exists {
			function.Allocation[b.op2.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}

	if value, exists := function.Allocation[b.result.String()]; !exists {
		function.Allocation[b.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (b *BinaryOp) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[b.op1.String()]; exists {
		if value.Spilled {
			b.op1 = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(b.op1, b.ty, spillRegister))
		} else { b.op1 = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[b.op2.String()]; exists {
		if value.Spilled {
			b.op2 = NewPhysicalRegister("x10")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(b.op2, b.ty, spillRegister))
		} else { b.op2 = NewPhysicalRegister(value.Register) }
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
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(b.ty, b.result, spillRegister))
		}
	}
}

func (b *BinaryOp) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	instrs = append(instrs, codegen.NewMove("x9", b.op1.GetValue()))
	instrs = append(instrs, codegen.NewMove("x10", b.op2.GetValue()))
	instrs = append(instrs, codegen.NewBinaryOp(b.operator.Translate(), b.result.GetValue(), "x9", "x10"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

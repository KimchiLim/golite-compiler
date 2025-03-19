package llvm

import ("golite/codegen" ; "golite/types" ; "bytes")

type Icmp struct {
	result LLVMOperand
	cond   LLVMOperator
	ty     types.Type
	op1    LLVMOperand
	op2    LLVMOperand
}

func NewIcmp(result LLVMOperand, cond LLVMOperator, ty types.Type, op1 LLVMOperand, op2 LLVMOperand) *Icmp {
	return &Icmp{result, cond, ty, op1, op2}
}

func (i *Icmp) String() string {
	var out bytes.Buffer
	out.WriteString(i.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("icmp")
	out.WriteString(" ")
	out.WriteString(i.cond.String())
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(i.ty))
	out.WriteString(" ")
	out.WriteString(i.op1.String())
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(i.op2.String())
	out.WriteString("\n")
	return out.String()
}

func (i *Icmp) MemtoReg(block *Block) {
	if value, exists := block.DefsMap[i.op1.String()]; exists { i.op1 = value }
	if value, exists := block.DefsMap[i.op2.String()]; exists { i.op2 = value }

	if _, exists := block.DefsMap[i.result.String()]; exists {
		block.DefsMap[i.result.String()] = i.result
	}

	block.RegisterInstrs = append(block.RegisterInstrs, i)
}

func (i *Icmp) ComputeLiveRange(function *LLVMFunction, position int) {
	if _, ok := i.op1.(*LLVMRegister); ok {
		if value, exists := function.Allocation[i.op1.String()]; !exists {
			function.Allocation[i.op1.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}

	if _, ok := i.op2.(*LLVMRegister); ok {
		if value, exists := function.Allocation[i.op2.String()]; !exists {
			function.Allocation[i.op2.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}

	if value, exists := function.Allocation[i.result.String()]; !exists {
		function.Allocation[i.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (i *Icmp) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[i.op1.String()]; exists {
		if value.Spilled {
			i.op1 = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(i.op1, i.ty, spillRegister))
		} else { i.op1 = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[i.op2.String()]; exists {
		if value.Spilled {
			i.op2 = NewPhysicalRegister("x10")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(i.op2, i.ty, spillRegister))
		} else { i.op2 = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[i.result.String()]; exists {
		if value.Spilled {
			i.result = NewPhysicalRegister("x9")
		} else { i.result = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, i)

	if value, exists := function.Allocation[i.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(i.ty, i.result, spillRegister))
		}
	}
}

func (i *Icmp) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	instrs = append(instrs, codegen.NewMove("x9", i.op1.GetValue()))
	instrs = append(instrs, codegen.NewMove("x10", i.op2.GetValue()))
	instrs = append(instrs, codegen.NewCompare(i.cond.Translate(), i.result.GetValue(), "x9", "x10"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

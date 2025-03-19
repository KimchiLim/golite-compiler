package llvm

import ("golite/codegen" ; "golite/types" ; "bytes" ; "fmt")

type GetElementPtr struct {
	result LLVMOperand
	ty     types.Type
	ptrval LLVMOperand
	index  int
}

func NewGetElementPtr(result LLVMOperand, ty types.Type, ptrval LLVMOperand, index int) *GetElementPtr {
	if value, ok := result.(*LLVMRegister); ok { value.pointsToStack = false }
	return &GetElementPtr{result, ty, ptrval, index}
}

func (g *GetElementPtr) String() string {
	var out bytes.Buffer
	out.WriteString(g.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("getelementptr")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(g.ty))
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(g.ty) + "*")
	out.WriteString(" ")
	out.WriteString(g.ptrval.String())
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32 0")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32")
	out.WriteString(" ")
	out.WriteString(fmt.Sprintf("%d", g.index))
	out.WriteString("\n")
	return out.String()
}

func (g *GetElementPtr) MemtoReg(block *Block) {
	if value, exists := block.DefsMap[g.ptrval.String()]; exists { g.ptrval = value }

	if _, exists := block.DefsMap[g.result.String()]; exists {
		block.DefsMap[g.result.String()] = g.result
	}

	block.RegisterInstrs = append(block.RegisterInstrs, g)
}

func (g *GetElementPtr) ComputeLiveRange(function *LLVMFunction, position int) {
	if value, exists := function.Allocation[g.ptrval.String()]; !exists {
		function.Allocation[g.ptrval.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }

	if value, exists := function.Allocation[g.result.String()]; !exists {
		function.Allocation[g.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (g *GetElementPtr) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[g.ptrval.String()]; exists {
		if value.Spilled {
			g.ptrval = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(g.ptrval, g.ty, spillRegister))
		} else { g.ptrval = NewPhysicalRegister(value.Register) }
	}

	if value, exists := function.Allocation[g.result.String()]; exists {
		if value.Spilled {
			g.result = NewPhysicalRegister("x9")
		} else { g.result = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, g)

	if value, exists := function.Allocation[g.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(g.ty, g.result, spillRegister))
		}
	}
}

func (g *GetElementPtr) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instr codegen.Instruction
	offset := fmt.Sprintf("#%d", g.index * 8)
	instr = codegen.NewBinaryOp(codegen.ADD, g.result.GetValue(), g.ptrval.GetValue(), offset)
	block.AssemblyInstrs = append(block.AssemblyInstrs, instr)
}

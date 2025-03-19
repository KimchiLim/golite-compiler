package llvm

import ("golite/codegen" ; "golite/st" ; "bytes" ; "fmt")

type Malloc struct {
	result LLVMOperand
	stval  *st.StructEntry
}

func NewMalloc(result LLVMOperand, stval *st.StructEntry) *Malloc {
	if value, ok := result.(*LLVMRegister); ok { value.pointsToStack = false }
	return &Malloc{result, stval}
}

func (m *Malloc) String() string {
	var out bytes.Buffer
	out.WriteString(m.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("call")
	out.WriteString(" ")
	out.WriteString("i8*")
	out.WriteString(" ")
	out.WriteString("@malloc")
	out.WriteString("(")
	out.WriteString("i32")
	out.WriteString(" ")
	out.WriteString(fmt.Sprintf("%d", len(m.stval.Fields) * 8))
	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}

func (m *Malloc) MemtoReg(block *Block) {
	block.DefsMap[m.result.String()] = m.result
	block.RegisterInstrs = append(block.RegisterInstrs, m)
}

func (m *Malloc) ComputeLiveRange(function *LLVMFunction, position int) {
	if value, exists := function.Allocation[m.result.String()]; !exists {
		function.Allocation[m.result.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (m *Malloc) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	ty := m.result.GetType()

	if value, exists := function.Allocation[m.result.String()]; exists {
		if value.Spilled {
			m.result = NewPhysicalRegister("x9")
		} else { m.result = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, m)

	if value, exists := function.Allocation[m.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(ty, m.result, spillRegister))
		}
	}
}

func (m *Malloc) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	instrs = append(instrs, codegen.NewMove("x0", fmt.Sprintf("#%d", len(m.stval.Fields) * 8))) 
	instrs = append(instrs, codegen.NewCall("malloc"))
	instrs = append(instrs, codegen.NewMove(m.result.GetValue(), "x0"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

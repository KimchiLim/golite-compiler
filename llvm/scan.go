package llvm

import ("golite/codegen" ; "golite/st" ; "golite/types" ; "bytes" ; "fmt")

type Scan struct {
	arg    LLVMOperand
	target LLVMOperand
}

func NewScan(arg LLVMOperand) *Scan {
	return &Scan{arg, nil}
}

func (s *Scan) String() string {
	var out bytes.Buffer
	out.WriteString("call")
	out.WriteString(" ")
	out.WriteString("i32")
	out.WriteString(" ")
	out.WriteString("(i8*, ...)")
	out.WriteString(" ")
	out.WriteString("@scanf")
	out.WriteString("(")
	out.WriteString("i8*")
	out.WriteString(" ")
	out.WriteString("getelementptr")
	out.WriteString(" ")
	out.WriteString("inbounds")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString("[4 x i8]")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("[4 x i8]*")
	out.WriteString(" ")
	out.WriteString("@.read")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32 0")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32 0")
	out.WriteString(")")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i64*")
	out.WriteString(" ")
	out.WriteString(s.arg.String())
	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}

func (s *Scan) MemtoReg(block *Block) {
	var scratchEntry, registerEntry *st.VarEntry
	var scratch, register LLVMOperand

	scratchEntry = st.NewVarEntry(".read_scratch", s.arg.GetType(), st.GLOBAL, nil)
	scratch = NewRegister(scratchEntry)
	s.target = s.arg
	s.arg = scratch

	registerEntry = st.NewVarEntry(NewRegisterLabel(), s.target.GetType(), st.LOCAL, nil)
	register = NewRegister(registerEntry)

	block.RegisterInstrs = append(block.RegisterInstrs, s)
	block.RegisterInstrs = append(block.RegisterInstrs, NewLoad(register, s.target.GetType(), scratch))
	block.DefsMap[s.target.String()] = register
}

func (s *Scan) ComputeLiveRange(function *LLVMFunction, position int) {
	if s.target.String()[0] == '@' { return }
	if value, exists := function.Allocation[s.target.String()]; !exists {
		function.Allocation[s.target.String()] = NewRegisterAlloc(position)
	} else { value.End = position - 1 }
}

func (s *Scan) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[s.target.String()]; exists {
		if value.Spilled {
			s.target = NewPhysicalRegister("x9")
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(s.target, types.IntTySig, spillRegister))
		} else { s.target = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, s)

	if value, exists := function.Allocation[s.target.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(types.IntTySig, s.target, spillRegister))
		}
	}
}

func (s *Scan) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", "#128"))
	for i := 0; i < 16; i++ { instrs = append(instrs, codegen.NewStore(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }

	instrs = append(instrs, codegen.NewAdrp("x0", ".READ"))
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "x0", "x0", ":lo12:.READ"))

	instrs = append(instrs, codegen.NewAdrp("x1", ".read_scratch"))
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "x1", "x1", ":lo12:.read_scratch"))

	instrs = append(instrs, codegen.NewCall("scanf"))

	for i := 0; i < 16; i++ { instrs = append(instrs, codegen.NewLoad(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", "#128"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

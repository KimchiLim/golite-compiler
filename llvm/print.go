package llvm

import ("golite/codegen" ; "bytes" ; "fmt" ; "strings")

type Print struct {
	fstring	*FString
	args    []LLVMOperand
}

func NewPrint(fstring *FString, args []LLVMOperand) *Print {
	return &Print{fstring, args}
}

func (p *Print) String() string {
	var out bytes.Buffer
	out.WriteString("call")
	out.WriteString(" ")
	out.WriteString("i32")
	out.WriteString(" ")
	out.WriteString("(i8*, ...)")
	out.WriteString(" ")
	out.WriteString("@printf")
	out.WriteString("(")
	out.WriteString("i8*")
	out.WriteString(" ")
	out.WriteString("getelementptr")
	out.WriteString(" ")
	out.WriteString("inbounds")
	out.WriteString(" ")
	out.WriteString("(")
	out.WriteString("[")
	out.WriteString(fmt.Sprintf("%d", p.fstring.Length))
	out.WriteString(" x i8], [")
	out.WriteString(fmt.Sprintf("%d", p.fstring.Length))
	out.WriteString(" x i8]*")
	out.WriteString(" ")
	out.WriteString("@." + p.fstring.Name)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32 0")
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("i32 0")
	out.WriteString(")")

	for _, pointer := range p.args {
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString("i64")
		out.WriteString(" ")
		out.WriteString(pointer.String())
	}

	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}

func (p *Print) MemtoReg(block *Block) {
	for i, pointer := range p.args {
		if pointer.String()[0] == '@' { continue }
		if value, exists := block.DefsMap[pointer.String()]; exists { p.args[i] = value }
	}

	block.RegisterInstrs = append(block.RegisterInstrs, p)
}

func (p *Print) ComputeLiveRange(function *LLVMFunction, position int) {
	for _, pointer := range p.args {
		if pointer.String()[0] == '@' { continue }
		if _, ok := pointer.(*LLVMRegister); ok {
			if value, exists := function.Allocation[pointer.String()]; !exists {
				function.Allocation[pointer.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}
}

func (p *Print) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	for i, pointer := range p.args {
		ty := pointer.GetType()

		if value, exists := function.Allocation[pointer.String()]; exists {
			if value.Spilled {
				if i % 2 == 0 {
					p.args[i] = NewPhysicalRegister("x9")
				} else {
					p.args[i] = NewPhysicalRegister("x10")
				}
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(p.args[i], ty, spillRegister))
			} else { p.args[i] = NewPhysicalRegister(value.Register) }
		}
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, p)
}

func (p *Print) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction
	var pushed_size, popped_size int
	num_popped := min(7, len(p.args))

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", "#128"))
	for i := 0; i < 16; i++ { instrs = append(instrs, codegen.NewStore(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }

	pushed_size = (len(p.args) + 1) * 8
	if len(p.args) % 2 == 0 { pushed_size += 8 }

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", fmt.Sprintf("#%d", pushed_size)))

	for i := 0; i < len(p.args); i++ {
		if _, ok := p.args[i].(*LLVMImmediate); ok {
			instrs = append(instrs, codegen.NewMove("x9", p.args[i].GetValue()))
			instrs = append(instrs, codegen.NewStore("x9", fmt.Sprintf("sp, #%d", (i + 1) * 8)))
		} else {
			instrs = append(instrs, codegen.NewStore(p.args[i].GetValue(), fmt.Sprintf("sp, #%d", (i + 1) * 8)))
		}
	}

	instrs = append(instrs, codegen.NewAdrp("x10", fmt.Sprintf(".%s", strings.ToUpper(p.fstring.Name))))
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "x10", "x10", fmt.Sprintf(":lo12:.%s", strings.ToUpper(p.fstring.Name))))
	instrs = append(instrs, codegen.NewMove("x0", "x10"))

	popped_size = num_popped * 8
	if num_popped % 2 == 1 { popped_size += 8 }

	for i := 0; i < num_popped; i++ { instrs = append(instrs, codegen.NewLoad(fmt.Sprintf("x%d", i + 1), fmt.Sprintf("sp, #%d", (i + 1) * 8))) }
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", fmt.Sprintf("#%d", popped_size)))

	instrs = append(instrs, codegen.NewCall("printf"))

	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", fmt.Sprintf("#%d", pushed_size - popped_size)))

	for i := 0; i < 16; i++ { instrs = append(instrs, codegen.NewLoad(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", "#128"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

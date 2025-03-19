package llvm

import ("golite/codegen" ; "golite/st" ; "golite/types" ; "bytes" ; "fmt" ; "strconv")

type Call struct {
	result LLVMOperand
	ty     types.Type
	fnval  *st.FuncEntry
	args   []LLVMOperand
}

func NewCall(result LLVMOperand, ty types.Type, fnval *st.FuncEntry, args []LLVMOperand) *Call {
	return &Call{result, ty, fnval, args}
}

func (c *Call) String() string {
	var out bytes.Buffer

	if c.result != nil {
		out.WriteString(c.result.String())
		out.WriteString(" ")
		out.WriteString("=")
		out.WriteString(" ")
	}

	out.WriteString("call")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(c.ty))
	out.WriteString(" ")
	out.WriteString("@" + c.fnval.Name)
	out.WriteString("(")

	for i, op := range c.args {
		if i > 0 {out.WriteString(", ")}
		out.WriteString(TypesToLLVMType(op.GetType()))
		out.WriteString(" ")
		out.WriteString(op.String())
	}

	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}

func (c *Call) MemtoReg(block *Block) {
	for i, op := range c.args {
		if op.String()[0] == '@' { continue }
		if value, exists := block.DefsMap[op.String()]; exists { c.args[i] = value }
	}

	if c.result != nil {
		if _, exists := block.DefsMap[c.result.String()]; exists {
			block.DefsMap[c.result.String()] = c.result
		}
	}

	block.RegisterInstrs = append(block.RegisterInstrs, c)
}

func (c *Call) ComputeLiveRange(function *LLVMFunction, position int) {
	for _, op := range c.args {
		if op.String()[0] == '@' { continue }
		if _, ok := op.(*LLVMRegister); ok {
			if value, exists := function.Allocation[op.String()]; !exists {
				function.Allocation[op.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}

	if c.result != nil {
		if value, exists := function.Allocation[c.result.String()]; !exists {
			function.Allocation[c.result.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}
}

func (c *Call) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	for i, op := range c.args {
		ty := op.GetType()

		if value, exists := function.Allocation[op.String()]; exists {
			if value.Spilled {
				if i % 2 == 0 {
					c.args[i] = NewPhysicalRegister("x9")
				} else {
					c.args[i] = NewPhysicalRegister("x10")
				}
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(c.args[i], ty, spillRegister))
			} else { c.args[i] = NewPhysicalRegister(value.Register) }
		}
	}

	if c.result != nil {
		if value, exists := function.Allocation[c.result.String()]; exists {
			if value.Spilled {
				c.result = NewPhysicalRegister("x9")
			} else { c.result = NewPhysicalRegister(value.Register) }
		}
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, c)

	if c.result != nil {
		if value, exists := function.Allocation[c.result.String()]; exists {
			if value.Spilled {
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(c.ty, c.result, spillRegister))
			}
		}
	}
}

func (c *Call) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction
	var pushed_size, popped_size int
	num_popped := min(8, len(c.args))
	result := -1

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", "#128"))
	for i := 0; i < 16; i++ { instrs = append(instrs, codegen.NewStore(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }

	pushed_size = len(c.args) * 8
	if len(c.args) % 2 == 1 { pushed_size += 8 }

	instrs = append(instrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", fmt.Sprintf("#%d", pushed_size)))

	for i := 0; i < len(c.args); i++ {
		if _, ok := c.args[i].(*LLVMImmediate); ok {
			instrs = append(instrs, codegen.NewMove("x9", c.args[i].GetValue()))
			instrs = append(instrs, codegen.NewStore("x9", fmt.Sprintf("sp, #%d", i * 8)))
		} else {
			instrs = append(instrs, codegen.NewStore(c.args[i].GetValue(), fmt.Sprintf("sp, #%d", i * 8)))
		}
	}

	popped_size = num_popped * 8
	if num_popped % 2 == 1 { popped_size -= 8 }

	for i := 0; i < num_popped; i++ { instrs = append(instrs, codegen.NewLoad(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8))) }
	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", fmt.Sprintf("#%d", popped_size)))

	instrs = append(instrs, codegen.NewCall(c.fnval.Name))

	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", fmt.Sprintf("#%d", pushed_size - popped_size)))

	if c.result != nil {
		instrs = append(instrs, codegen.NewMove(c.result.GetValue(), "x0"))
		result, _ = strconv.Atoi(c.result.GetValue()[1:])
	}

	for i := 0; i < 16; i++ {
		if i == result { continue }
		instrs = append(instrs, codegen.NewLoad(fmt.Sprintf("x%d", i), fmt.Sprintf("sp, #%d", i * 8)))
	}

	instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "sp", "sp", "#128"))

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

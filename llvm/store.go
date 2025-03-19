package llvm

import ("golite/codegen" ; "golite/types" ; "bytes" ; "fmt")

type Store struct {
	ty      types.Type
	value   LLVMOperand
	pointer LLVMOperand
}

func NewStore(ty types.Type, value LLVMOperand, pointer LLVMOperand) *Store {
	return &Store{ty, value, pointer}
}

func (s *Store) String() string {
	var out bytes.Buffer
	out.WriteString("store")
	out.WriteString(" ")

	if s.value != nil {
		out.WriteString(TypesToLLVMType(s.ty))
		out.WriteString(" ")
		out.WriteString(s.value.String())

	} else if s.ty == types.IntTySig || s.ty == types.BoolTySig {
		out.WriteString("i64 0")
	} else { out.WriteString("ptr null") }

	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(s.ty) + "*")
	out.WriteString(" ")
	out.WriteString(s.pointer.String())
	out.WriteString("\n")
	return out.String()
}

func (s *Store) MemtoReg(block *Block) {
	if s.value != nil {
		if value, exists := block.DefsMap[s.value.String()]; exists { s.value = value }
	}

	if s.pointer.PointsToStack() {
		block.DefsMap[s.pointer.String()] = s.value
	} else {
		block.RegisterInstrs = append(block.RegisterInstrs, s)
	}
}

func (s *Store) ComputeLiveRange(function *LLVMFunction, position int) {
	if s.value != nil {
		if _, ok := s.value.(*LLVMRegister); ok {
			if value, exists := function.Allocation[s.value.String()]; !exists {
				function.Allocation[s.value.String()] = NewRegisterAlloc(position)
			} else { value.End = position - 1 }
		}
	}

	if s.pointer.String()[0] != '@' {
		if value, exists := function.Allocation[s.pointer.String()]; !exists {
			function.Allocation[s.pointer.String()] = NewRegisterAlloc(position)
		} else { value.End = position - 1 }
	}
}

func (s *Store) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {
	if s.value != nil {
		if value, exists := function.Allocation[s.value.String()]; exists {
			if value.Spilled {
				s.value = NewPhysicalRegister("x9")
				spillRegister := NewPhysicalRegister(value.Register)
				block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewLoad(s.value, s.ty, spillRegister))
			} else { s.value = NewPhysicalRegister(value.Register) }
		}	
	}

	if value, exists := function.Allocation[s.pointer.String()]; exists {
		if value.Spilled {
			s.pointer = NewPhysicalRegister("x9")
		} else { s.pointer = NewPhysicalRegister(value.Register) }
	}

	block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, s)

	if value, exists := function.Allocation[s.pointer.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(s.ty, s.pointer, spillRegister))
		}
	}
}

func (s *Store) TranslateToAssembly(function *LLVMFunction, block *Block) {
	var instrs []codegen.Instruction

	if s.value != nil {
		if s.pointer.String()[0] == '@' {
			instrs = append(instrs, codegen.NewAdrp("x10", s.pointer.GetValue()))
			instrs = append(instrs, codegen.NewBinaryOp(codegen.ADD, "x10", "x10", fmt.Sprintf(":lo12:%s", s.pointer.GetValue())))

			if _, ok := s.value.(*LLVMImmediate); ok {
				instrs = append(instrs, codegen.NewMove("x9", s.value.GetValue()))
				instrs = append(instrs, codegen.NewStore("x9", "x10"))
			} else {
				instrs = append(instrs, codegen.NewStore(s.value.GetValue(), "x10"))
			}
		} else {
			if _, ok := s.value.(*LLVMImmediate); ok {
				instrs = append(instrs, codegen.NewMove("x9", s.value.GetValue()))
				instrs = append(instrs, codegen.NewStore("x9", s.pointer.GetValue()))
			} else {
				instrs = append(instrs, codegen.NewStore(s.value.GetValue(), s.pointer.GetValue()))
			}
		}
	}

	block.AssemblyInstrs = append(block.AssemblyInstrs, instrs...)
}

package llvm

import ("golite/st" ; "golite/types" ; "bytes" ; "fmt" ; "slices" ; "sort")

type LLVMProgram struct {
	SourceName   string
	TargetTriple string
	Tables       *st.SymbolTables
	FuncDefs     map[*st.FuncEntry]*LLVMFunction
	FStrings	 []*FString
}

type LLVMFunction struct {
	Name       string
	Block      *Block
	Registers  map[*st.VarEntry]*LLVMRegister
	Allocation map[string]*RegisterAlloc
}

type RegisterAlloc struct {
	Start    int
	End      int
	Register string
	Spilled  bool
}

func NewLLVMFunction(function *st.FuncEntry, block *Block) *LLVMFunction {
	return &LLVMFunction{function.Name, block, make(map[*st.VarEntry]*LLVMRegister), make(map[string]*RegisterAlloc)}
}

func NewLLVMProgram(sourceName string, targetTriple string, tables *st.SymbolTables) *LLVMProgram {
	return &LLVMProgram{sourceName, targetTriple, tables, make(map[*st.FuncEntry]*LLVMFunction), nil}
}

func NewRegisterAlloc(position int) *RegisterAlloc {
	return &RegisterAlloc{position, position - 1, "", false}
}

func (program *LLVMProgram) LLVMString(stackBased bool) string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("source_filename = \"%s\"\n", program.SourceName))
	out.WriteString(fmt.Sprintf("target triple = \"%s\"\n", program.TargetTriple))

	for _, structEntry := range program.Tables.Structs.Table {
		out.WriteString(fmt.Sprintf("%%struct.%s = type", structEntry.Name))
		out.WriteString(" ")
		out.WriteString("{")

		for i, fieldEntry := range structEntry.Fields {
			if i > 0 { out.WriteString(", ") }
			out.WriteString(TypesToLLVMType(fieldEntry.Ty))
			if fieldEntry.Ty != types.IntTySig && fieldEntry.Ty != types.BoolTySig { out.WriteString("*") }
		}

		out.WriteString("}")
		out.WriteString("\n")
	}

	for _, varEntry := range program.Tables.Globals.Table {
		var constant string
		if _, ok := varEntry.Ty.(*types.PointerTy); ok { constant = "null"}
		if varEntry.Ty == types.IntTySig || varEntry.Ty == types.BoolTySig { constant = "0" }
		out.WriteString(fmt.Sprintf("@%s = common global %s %s\n", varEntry.Name, TypesToLLVMType(varEntry.Ty), constant))
	}

	out.WriteString("@.read_scratch = common global i64 0\n")
	out.WriteString("\n")

	for funcEntry, llvmFunc := range program.FuncDefs {
		if funcEntry.Name == "main" {
			out.WriteString("define i64 @main")
		} else { out.WriteString(fmt.Sprintf("define %s @%s", TypesToLLVMType(funcEntry.ReturnTy), funcEntry.Name)) }

		out.WriteString("(")

		for i, paramEntry := range funcEntry.Parameters {
			if i > 0 { out.WriteString(", ") }
			out.WriteString(fmt.Sprintf("%s %%%s", TypesToLLVMType(paramEntry.Ty), paramEntry.Name))
		}

		out.WriteString(")")
		out.WriteString("\n")
		out.WriteString("{")
		out.WriteString("\n")

		visited := make(map[*Block]bool)
		queue := []*Block{llvmFunc.Block}
		if !stackBased {llvmFunc.Block.TranslateToRegister(funcEntry, nil, make(map[string]bool), make(map[*Block]int))}

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]

			if !visited[block] {
				visited[block] = true
				out.WriteString(block.Label)
				out.WriteString(":")
				out.WriteString("\n")

				if stackBased { for _, instr := range block.StackInstrs { out.WriteString(fmt.Sprintf("\t%s", instr)) } }
				if !stackBased { for _, instr := range block.RegisterInstrs { out.WriteString(fmt.Sprintf("\t%s", instr)) } }
				for _, next := range block.Next { if !visited[next] { queue = append(queue, next) } }
			}
		}

		out.WriteString("}")
		out.WriteString("\n\n")
	}

	out.WriteString("declare i8* @malloc(i32)\n")
	out.WriteString("declare i32 @printf(i8*, ...)\n")
	out.WriteString("declare i32 @scanf(i8*, ...)\n")
	out.WriteString("declare void @free(i8*)\n")
	out.WriteString(fmt.Sprintf("@.read = private unnamed_addr constant [4 x i8] c\"%%ld\\00\", align 1\n"))
	for _, fstring := range program.FStrings { out.WriteString(fstring.String()) }

	return out.String()
}

func (program *LLVMProgram) TranslateToOutOfSSA() {
	for funcEntry, llvmFunc := range program.FuncDefs {
		var intervals, active, spilled []string
		spilledCount := 1

		for _, parameter := range funcEntry.Parameters {
			llvmFunc.Allocation["%" + parameter.Name] = NewRegisterAlloc(0)
			intervals = append(intervals, "%" + parameter.Name)
		}

		visited := make(map[*Block]bool)
		queue := []*Block{llvmFunc.Block}
		position := 1

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]

			if !visited[block] {
				visited[block] = true
				for _, next := range block.Next { if !visited[next] { queue = append(queue, next) } }

				for _, instr := range block.RegisterInstrs {
					instr.ComputeLiveRange(llvmFunc, position)
					position++
				}
			}
		}

		for key := range llvmFunc.Allocation { if !slices.Contains(intervals, key) { intervals = append(intervals, key) } }
		sort.SliceStable(intervals, func(i, j int) bool { return llvmFunc.Allocation[intervals[i]].Start < llvmFunc.Allocation[intervals[j]].Start })


		registers := []string{
            "x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7", "x8", "x11", "x12", "x13", "x14",
            "x15", "x19", "x20", "x21", "x22", "x23", "x24", "x25", "x26", "x27", "x28",
		}

		for _, currKey := range intervals {
			curr := llvmFunc.Allocation[currKey]
			var next []string

			for _, activeKey := range active {
				if llvmFunc.Allocation[activeKey].End > curr.Start {
					next = append(next, activeKey)
				} else {
					if llvmFunc.Allocation[activeKey].Spilled {
						spilled = append(spilled, llvmFunc.Allocation[activeKey].Register)
					} else { registers = append(registers, llvmFunc.Allocation[activeKey].Register) }
				}
			}

			active = next

			if len(registers) > 0 {
				curr.Register = registers[0]
				registers = registers[1:]
			} else {
                if len(spilled) > 0 {
                    curr.Register = spilled[0]
                    spilled = spilled[1:]
                } else {
                    curr.Register = fmt.Sprintf("spill_%d", spilledCount)
                    spilledCount++
                }
				curr.Spilled = true
			}

			active = append(active, currKey)
		}

        visited = make(map[*Block]bool)
		queue = []*Block{llvmFunc.Block}
        
        for i := 1; i < spilledCount; i++ {
            spillRegister := NewPhysicalRegister(fmt.Sprintf("spill_%d", i))
            llvmFunc.Block.OutOfSSAInstrs = append(llvmFunc.Block.OutOfSSAInstrs, NewAlloca(spillRegister, types.IntTySig))
        }

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]

			if !visited[block] {
				visited[block] = true
				for _, instr := range block.RegisterInstrs { instr.TranslateToOutOfSSA(llvmFunc, block) }
				for _, next := range block.Next { if !visited[next] { queue = append(queue, next) } }
			}
		}

		visited = make(map[*Block]bool)
		queue = []*Block{llvmFunc.Block}

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]

			if !visited[block] {
				visited[block] = true
				for _, instr := range block.RegisterInstrs { if phi, ok := instr.(*Phi); ok { phi.TranslateToOutOfSSAPhi(llvmFunc, block) } }
				for _, next := range block.Next { if !visited[next] { queue = append(queue, next) } }
			}
		}
	}
}

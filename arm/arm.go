package arm

import ("golite/codegen" ; "golite/llvm" ; "golite/st" ; "bytes" ; "fmt" ; "strings")

type ARMProgram struct {
	Tables      *st.SymbolTables
	FuncDefs    map[*st.FuncEntry]*llvm.LLVMFunction
	FStrings	[]*llvm.FString
}

func NewARMProgram(LLVMProgram *llvm.LLVMProgram) *ARMProgram {
	LLVMProgram.TranslateToOutOfSSA()
	for _, llvmFunc := range LLVMProgram.FuncDefs {
		offset := len(llvmFunc.Allocation) * 8 + 8 + 8
		offset = ((offset + 15) / 16) * 16

		llvmFunc.Block.AssemblyInstrs = append(llvmFunc.Block.AssemblyInstrs, codegen.NewBinaryOp(codegen.SUB, "sp", "sp", fmt.Sprintf("#%d", offset)))
		llvmFunc.Block.AssemblyInstrs = append(llvmFunc.Block.AssemblyInstrs, codegen.NewStp("x29", "x30", "sp"))
		llvmFunc.Block.AssemblyInstrs = append(llvmFunc.Block.AssemblyInstrs, codegen.NewMove("x29", "sp"))

		visited := make(map[*llvm.Block]bool)
		queue := []*llvm.Block{llvmFunc.Block}

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]
			if !visited[block] {
				visited[block] = true
				for _, instr := range block.OutOfSSAInstrs { 
					instr.TranslateToAssembly(llvmFunc, block) 
				}
				for _, next := range block.Next { 
					if !visited[next] { queue = append(queue, next) } 
				}
			}
		}
	}

	return &ARMProgram{LLVMProgram.Tables, LLVMProgram.FuncDefs, LLVMProgram.FStrings}
}

func (program *ARMProgram) String() string {
	var out bytes.Buffer
	out.WriteString("\t.arch armv8-a\n")

	for _, varEntry := range program.Tables.Globals.Table { out.WriteString(fmt.Sprintf("\t.comm\t%s,8,8\n", varEntry.Name)) }
	out.WriteString(fmt.Sprintf("\t.comm\t.read_scratch,8,8\n"))

	out.WriteString("\t.text\n")
	out.WriteString("\n")

	for funcEntry, llvmFunc := range program.FuncDefs {
		out.WriteString(fmt.Sprintf("\t.type %s, %%function\n", funcEntry.Name))
		out.WriteString(fmt.Sprintf("\t.global %s\n", funcEntry.Name))
		out.WriteString("\t.p2align 2\n")
		out.WriteString(fmt.Sprintf("%s:\n", funcEntry.Name))

		visited := make(map[*llvm.Block]bool)
		queue := []*llvm.Block{llvmFunc.Block}

		for len(queue) > 0 {
			block := queue[0]
			queue = queue[1:]

			if !visited[block] {
				visited[block] = true
				out.WriteString(fmt.Sprintf(".%s:\n", block.Label))
				for _, instr := range block.AssemblyInstrs { out.WriteString(fmt.Sprintf("\t%s", instr)) }
				for _, next := range block.Next { if !visited[next] { queue = append(queue, next) } }
			}
		}

		out.WriteString(fmt.Sprintf("\t.size %s, (. - %s)\n", funcEntry.Name, funcEntry.Name))
		out.WriteString("\n")
	}

	out.WriteString(".READ:\n")
	out.WriteString("\t.asciz\t\"%ld\"\n")
	out.WriteString("\t.size\t.READ, 4\n")
	out.WriteString("\n")

	for _, fstring := range program.FStrings {
		content := fstring.Content
		content = strings.ReplaceAll(content, "\\0A", "\\n")
		content = strings.ReplaceAll(content, "\\00", "")

		out.WriteString(fmt.Sprintf(".%s:\n", strings.ToUpper(fstring.Name)))
		out.WriteString(fmt.Sprintf("\t.asciz\t\"%s\"\n", content))
		out.WriteString(fmt.Sprintf("\t.size\t.PRINT, %d\n", fstring.Length))
		out.WriteString("\n")
	}

	return out.String()
}

package llvm

import ("golite/codegen" ; "golite/st" ; "golite/types" ; "fmt")

var blockCount int
func init() { blockCount = 0 }

func NewBlockLabel() string {
	label := fmt.Sprintf("L%d", blockCount)
	blockCount++
	return label
}

type Block struct {
	Label          string
	Prev           []*Block
	Next           []*Block
	HasReturn      bool
	ReturnType     types.Type
	DefsMap	       map[string]LLVMOperand
	PhiInstrs      map[string]*Phi
	StackInstrs    []LLVMInstruction
	RegisterInstrs []LLVMInstruction
	OutOfSSAInstrs []LLVMInstruction
	AssemblyInstrs []codegen.Instruction
}

func NewBlock() *Block {
	return &Block{NewBlockLabel(), nil, nil, false, types.VoidTySig, make(map[string]LLVMOperand), make(map[string]*Phi), nil, nil, nil, nil}
}

func (block *Block) String() string {
	return fmt.Sprintf("%%%s", block.Label)
}

func (block *Block) AddNext(next *Block) {
	block.Next = append(block.Next, next)
	next.Prev = append(next.Prev, block)
}

func (block *Block) AddReturn(returnType types.Type) {
	block.HasReturn = true
	block.ReturnType = returnType
}

func (block *Block) TranslateToRegister(funcEntry *st.FuncEntry, predecessor *Block, taken map[string]bool, visited map[*Block]int) {
	if predecessor != nil {
		if taken[fmt.Sprintf("%s-%s", predecessor.Label, block.Label)] { return }
		if visited[block] == len(block.Prev) { return }
		taken[fmt.Sprintf("%s-%s", predecessor.Label, block.Label)] = true
	}

	visited[block]++

	if len(block.Prev) > 1 {
		if visited[block] == 1 {
			i := NewPhiLabel()

			for key, _ := range block.DefsMap {
				if varEntry := funcEntry.Variables.Contains(key[1:]); varEntry != nil {
					entry := st.NewVarEntry(fmt.Sprintf("%s_%s", varEntry.Name, i), varEntry.Ty, st.LOCAL, varEntry.Token)
					register := NewRegister(entry)
					block.PhiInstrs[varEntry.Name] = NewPhi(register, entry.Ty)
					block.DefsMap["%" + varEntry.Name] = block.PhiInstrs[varEntry.Name].result
				}
			}
		}

		for key, _ := range predecessor.DefsMap {
			if varEntry := funcEntry.Variables.Contains(key[1:]); varEntry != nil {
				if value, exists := predecessor.DefsMap["%" + varEntry.Name]; exists {
					block.PhiInstrs[varEntry.Name].AddValue(value, predecessor)
					block.DefsMap[key] = block.PhiInstrs[varEntry.Name].result
				}
			}
		}
	}

	if visited[block] == 1 {
		for _, instr := range block.PhiInstrs { instr.MemtoReg(block) }
		for _, instr := range block.StackInstrs { instr.MemtoReg(block) }
	}

	for _, next := range block.Next {
		if visited[next] == len(next.Prev) { continue }
		next.DefsMap = make(map[string]LLVMOperand)
		for key, value := range block.DefsMap { next.DefsMap[key] = value }
		next.TranslateToRegister(funcEntry, block, taken, visited)
	}
}

package llvm

import ("golite/types" ; "bytes")

type Alloca struct {
	result LLVMOperand
	ty     types.Type
}

func NewAlloca(result LLVMOperand, ty types.Type) *Alloca {
	return &Alloca{result, ty}
}

func (a *Alloca) String() string {
	var out bytes.Buffer
	out.WriteString(a.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("alloca")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(a.ty))
	out.WriteString("\n")
	return out.String()
}

func (a *Alloca) MemtoReg(block *Block) {
	block.DefsMap[a.result.String()] = NewImmediate(0, a.ty)
}

func (a *Alloca) ComputeLiveRange(function *LLVMFunction, position int) {}

func (a *Alloca) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {}

func (a *Alloca) TranslateToAssembly(function *LLVMFunction, block *Block) {}

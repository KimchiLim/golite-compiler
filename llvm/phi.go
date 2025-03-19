package llvm

import ("golite/types" ; "bytes" ; "fmt")

type Phi struct {
	result LLVMOperand
	ty     types.Type
	values []LLVMOperand
	blocks []*Block
}

var phiCount int
func init() { phiCount = 0 }

func NewPhiLabel() string {
	label := fmt.Sprintf("%d", phiCount)
	phiCount++
	return label
}

func NewPhi(result LLVMOperand, ty types.Type) *Phi {
	return &Phi{result, ty, nil, nil}
}

func (p *Phi) AddValue(value LLVMOperand, block *Block) {
	p.values = append(p.values, value)
	p.blocks = append(p.blocks, block)
}

func (p *Phi) String() string {
	var out bytes.Buffer
	out.WriteString(p.result.String())
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("phi")
	out.WriteString(" ")
	out.WriteString(TypesToLLVMType(p.ty))
	out.WriteString(" ")

	for i, _ := range p.values {
		if i > 0 { out.WriteString(", ") }
		out.WriteString("[")
		out.WriteString(p.values[i].String())
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString(p.blocks[i].String())
		out.WriteString("]")
	}

	out.WriteString("\n")
	return out.String()
}

func (p *Phi) MemtoReg(block *Block) {
	if len(p.values) > 0 { block.RegisterInstrs = append(block.RegisterInstrs, p) }
}

func (p *Phi) ComputeLiveRange(function *LLVMFunction, position int) {
	for _, val := range p.values {
		if _, ok := val.(*LLVMRegister); ok {
			if value, exists := function.Allocation[val.String()]; !exists {
				function.Allocation[val.String()] = NewRegisterAlloc(position)
			} else { value.End = max(value.End, position - 1) }
		}
	}

	if value, exists := function.Allocation[p.result.String()]; !exists {
		function.Allocation[p.result.String()] = NewRegisterAlloc(position)
	} else { value.End = max(value.End, position - 1) }
}

func (p *Phi) TranslateToOutOfSSA(function *LLVMFunction, block *Block) {}

func (p *Phi) TranslateToOutOfSSAPhi(function *LLVMFunction, block *Block) {
	if value, exists := function.Allocation[p.result.String()]; exists {
		if value.Spilled {
			p.result = NewPhysicalRegister("x10")
		} else { p.result = NewPhysicalRegister(value.Register) }
	}

	for i, _ := range p.values {
		ty := p.values[i].GetType()
		length := len(p.blocks[i].OutOfSSAInstrs)

		if value, exists := function.Allocation[p.values[i].String()]; exists {
			if value.Spilled {
				p.values[i] = NewPhysicalRegister("x9")
				spillRegister := NewPhysicalRegister(value.Register)
				p.blocks[i].OutOfSSAInstrs = append(p.blocks[i].OutOfSSAInstrs, NewLoad(p.values[i], ty, spillRegister))
			} else { p.values[i] = NewPhysicalRegister(value.Register) }
		}

		if length == 0 {
			p.blocks[i].OutOfSSAInstrs = append(p.blocks[i].OutOfSSAInstrs, NewMove(p.values[i], p.result))

		} else if br, ok := p.blocks[i].OutOfSSAInstrs[length - 1].(*Branch); ok {
			p.blocks[i].OutOfSSAInstrs = append(p.blocks[i].OutOfSSAInstrs[:length - 1], NewMove(p.values[i], p.result))
			p.blocks[i].OutOfSSAInstrs = append(p.blocks[i].OutOfSSAInstrs, br)

		} else { p.blocks[i].OutOfSSAInstrs = append(p.blocks[i].OutOfSSAInstrs, NewMove(p.values[i], p.result)) }
	}

	if value, exists := function.Allocation[p.result.String()]; exists {
		if value.Spilled {
			spillRegister := NewPhysicalRegister(value.Register)
			block.OutOfSSAInstrs = append(block.OutOfSSAInstrs, NewStore(p.ty, p.result, spillRegister))
		}
	}
}

func (p *Phi) TranslateToAssembly(function *LLVMFunction, block *Block) {}

package codegen

import ("bytes")

type Adrp struct {
	address string
	label   string
}

func NewAdrp(address string, label string) *Adrp {
	return &Adrp{address, label}
}

func (a *Adrp) String() string {
	var out bytes.Buffer
	out.WriteString("adrp")
	out.WriteString(" ")
	out.WriteString(a.address)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(a.label)
	out.WriteString("\n")
	return out.String()
}

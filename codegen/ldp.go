package codegen

import ("bytes")

type Ldp struct {
	r1      string
	r2      string
	address string
}

func NewLdp(r1 string, r2 string, address string) *Ldp {
	return &Ldp{r1, r2, address}
}

func (l *Ldp) String() string {
	var out bytes.Buffer
	out.WriteString("ldp")
	out.WriteString(" ")
	out.WriteString(l.r1)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(l.r2)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("[")
	out.WriteString(l.address)
	out.WriteString("]")
	out.WriteString("\n")
	return out.String()
}

package codegen

import ("bytes")

type Load struct {
	register string
	address  string
}

func NewLoad(register string, address string) *Load {
	return &Load{register, address}
}

func (l *Load) String() string {
	var out bytes.Buffer
	out.WriteString("ldr")
	out.WriteString(" ")
	out.WriteString(l.register)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("[")
	out.WriteString(l.address)
	out.WriteString("]")
	out.WriteString("\n")
	return out.String()
}

package codegen

import ("bytes")

type Return struct {
	register string
}

func NewReturn(register string) *Return {
	return &Return{register}
}

func (r *Return) String() string {
	var out bytes.Buffer
	out.WriteString("ret")

	if r.register != "" {
		out.WriteString(" ")
		out.WriteString(r.register)
	}

	out.WriteString("\n")
	return out.String()
}

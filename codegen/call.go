package codegen

import ("bytes")

type Call struct {
	label string
}

func NewCall(label string) *Call {
	return &Call{label}
}

func (c *Call) String() string {
	var out bytes.Buffer
	out.WriteString("bl")
	out.WriteString(" ")
	out.WriteString(c.label)
	out.WriteString("\n")
	return out.String()
}

package codegen

import ("bytes")

type Move struct {
	dest string
	op   string
}

func NewMove(dest string, op string) *Move {
	return &Move{dest, op}
}

func (m *Move) String() string {
	var out bytes.Buffer
	out.WriteString("mov")
	out.WriteString(" ")
	out.WriteString(m.dest)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(m.op)
	out.WriteString("\n")
	return out.String()
}

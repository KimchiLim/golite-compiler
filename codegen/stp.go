package codegen

import ("bytes")

type Stp struct {
	r1      string
	r2      string
	address string
}

func NewStp(r1 string, r2 string, address string) *Stp {
	return &Stp{r1, r2, address}
}

func (s *Stp) String() string {
	var out bytes.Buffer
	out.WriteString("stp")
	out.WriteString(" ")
	out.WriteString(s.r1)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(s.r2)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("[")
	out.WriteString(s.address)
	out.WriteString("]")
	out.WriteString("\n")
	return out.String()
}

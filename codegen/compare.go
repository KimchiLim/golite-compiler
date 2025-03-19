package codegen

import ("bytes")

type Compare struct {
	cond   Operator
	result string
	op1    string
	op2    string
}

func NewCompare(cond Operator, result string, op1 string, op2 string) *Compare {
	return &Compare{cond, result, op1, op2}
}

func (c *Compare) String() string {
	var out bytes.Buffer
	out.WriteString("cmp")
	out.WriteString(" ")
	out.WriteString(c.op1)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(c.op2)
	out.WriteString("\n\t")
	out.WriteString("cset")
	out.WriteString(" ")
	out.WriteString(c.result)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(c.cond.String())
	out.WriteString("\n")
	return out.String()
}

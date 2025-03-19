package codegen

import ("bytes")

type BinaryOp struct {
	operator Operator
	result   string
	op1      string
	op2      string
}

func NewBinaryOp(operator Operator, result string, op1 string, op2 string) *BinaryOp {
	return &BinaryOp{operator, result, op1, op2}
}

func (b *BinaryOp) String() string {
	var out bytes.Buffer
	out.WriteString(b.operator.String())
	out.WriteString(" ")
	out.WriteString(b.result)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(b.op1)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString(b.op2)
	out.WriteString("\n")
	return out.String()
}

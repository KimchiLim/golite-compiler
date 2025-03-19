package codegen

import ("bytes")

type Branch struct {
	cond  string
	label string
}

func NewBranch(cond string, label string) *Branch {
	return &Branch{cond, label}
}

func (b *Branch) String() string {
	var out bytes.Buffer

	if b.cond != "" {
		out.WriteString("cmp")
		out.WriteString(" ")
		out.WriteString(b.cond)
		out.WriteString(",")
		out.WriteString(" ")
		out.WriteString("#0")
		out.WriteString("\n\t")
	}

	out.WriteString("b")

	if b.cond != "" {
		out.WriteString(".")
		out.WriteString("eq")
	}

	out.WriteString(" ")
	out.WriteString(b.label)
	out.WriteString("\n")
	return out.String()
}

package ast

import ("golite/token" ; "golite/types" ; "bytes")

type Declaration struct {
	*token.Token
	identifiers []*Variable
	ty          types.Type
}

func NewDeclaration(identifiers []*Variable, ty types.Type, token *token.Token) *Declaration {
	return &Declaration{token, identifiers, ty}
}

func (d *Declaration) String() string {
	var out bytes.Buffer
	out.WriteString("var")
	out.WriteString(" ")
	
	for i, identifier := range d.identifiers {
		if (i > 0) {out.WriteString(", ")}
		out.WriteString(identifier.String())
	}

	out.WriteString(" ")
	out.WriteString(d.ty.String())
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

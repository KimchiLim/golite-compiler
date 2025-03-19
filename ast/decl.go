package ast

import ("golite/token" ; "golite/types" ; "bytes")

type Decl struct {
	*token.Token
	identifier *Variable
	ty         types.Type
}

func NewDecl(identifier *Variable, ty types.Type, token *token.Token) *Decl {
	return &Decl{token, identifier, ty}
}

func (d *Decl) String() string {
	var out bytes.Buffer
	out.WriteString(d.identifier.String())
	out.WriteString(" ")
	out.WriteString(d.ty.String())
	return out.String()
}

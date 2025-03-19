package ast

import ("golite/token" ; "bytes")

type TypeDeclaration struct {
	*token.Token
	identifier *Variable
	fields     []*Decl
}

func NewTypeDeclaration(identifier *Variable, fields []*Decl, token *token.Token) *TypeDeclaration {
	return &TypeDeclaration{token, identifier, fields}
}

func (t *TypeDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString("type")
	out.WriteString(" ")
	out.WriteString(t.identifier.String())
	out.WriteString(" ")
	out.WriteString("struct")
	out.WriteString(" ")
	out.WriteString("{")
	out.WriteString("\n")

	for _, decl := range t.fields {
		out.WriteString("    ")
		out.WriteString(decl.String())
		out.WriteString(";")
		out.WriteString("\n")
	}
	
	out.WriteString("}")
	out.WriteString(";")
	out.WriteString("\n")
	return out.String()
}

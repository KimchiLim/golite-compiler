package ast

import ("golite/token" ; "golite/types" ; "bytes" ; "strings")

type Function struct {
	*token.Token
	identifier   *Variable
	parameters   []*Decl
	returnType   types.Type
	declarations []*Declaration
	statements   []Statement	
}

func NewFunction(identifier *Variable, parameters []*Decl, returnType types.Type, declarations []*Declaration, statements []Statement, token *token.Token) *Function {
	return &Function{token, identifier, parameters, returnType, declarations, statements}
}

func (f *Function) String() string {
	var out bytes.Buffer
	var prev = false

	out.WriteString("func")
	out.WriteString(" ")
	out.WriteString(f.identifier.String())
	out.WriteString("(")

	for i, decl := range f.parameters {
		if (i > 0) {out.WriteString(", ")}
		out.WriteString(decl.String())
	}

	out.WriteString(")")

	if (f.returnType != nil) {
		out.WriteString(" ")
		out.WriteString(f.returnType.String())
	}

	out.WriteString(" ")
	out.WriteString("{")
	if len(f.declarations) > 0 {out.WriteString("\n\n")}

	for _, declaration := range f.declarations {
		out.WriteString("    ")
		out.WriteString(declaration.String())
	}

	if len(f.declarations) == 0 && len(f.statements) > 0 {out.WriteString("\n")}

	for i, statement := range f.statements {
		switch statement.(type) {
			case *Conditional, *Loop:
				out.WriteString(formatStatement(statement, 1, prev))
				if i < len(f.statements) - 1 {out.WriteString("\n")}
				prev = true
			default:
				out.WriteString(formatStatement(statement, 1, prev))				
				if i == len(f.statements) - 1 {out.WriteString("\n")}
				prev = false
		}
	}

	out.WriteString("}")
	out.WriteString("\n")
	return out.String()
}

func formatStatement(statement Statement, indentLevel int, line bool) string {
	var out bytes.Buffer
	indent := strings.Repeat("    ", indentLevel)

	switch stmt := statement.(type) {
		case *Conditional:
			if !line {out.WriteString("\n")}
			out.WriteString(indent + "if " + stmt.expression.String() + " {")

			if len(stmt.ifBlock) > 0 {
				out.WriteString("\n")
				for _, nestedStmt := range stmt.ifBlock {
					out.WriteString(formatStatement(nestedStmt, indentLevel + 1, false))
				}
			}

			if len(stmt.elseBlock) > 0 {
				out.WriteString("\n" + indent + "} else {\n")
				for _, nestedStmt := range stmt.elseBlock {
					out.WriteString(formatStatement(nestedStmt, indentLevel + 1, false))
				}
			}

			out.WriteString(indent + "}\n")

		case *Loop:
			if !line {out.WriteString("\n")}
			out.WriteString(indent + "for " + stmt.expression.String() + " {")

			if len(stmt.forBlock) > 0 {
				out.WriteString("\n")
				for _, nestedStmt := range stmt.forBlock {
					out.WriteString(formatStatement(nestedStmt, indentLevel + 1, false))
				}
			}

			out.WriteString(indent + "}\n")

		default: 
			if !line {out.WriteString("\n")}
			out.WriteString(indent + statement.String() + "\n")
	}

	return out.String()
}

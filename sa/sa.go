package sa

import ("golite/ast" ; "golite/cfg" ; "golite/context" ; "golite/st")

func Execute(program *ast.Program) *st.SymbolTables {
	// Define the compiler symbol tables
	tables := st.NewSymbolTables()
	var blocks []*cfg.Block

	// First build the symbol table(s) for all global declarations
	errors := make([]*context.CompilerError, 0)
	errors = program.BuildSymbolTable(errors, tables)

	if !context.HasErrors(errors) {
		// Second perform type checking!
		errors := make([]*context.CompilerError, 0)
		errors = program.TypeCheck(errors, tables, blocks)
		if !context.HasErrors(errors) { return tables }
	}
	return nil
}

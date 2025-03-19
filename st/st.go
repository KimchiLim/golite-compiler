package st

import ("golite/context" ; "golite/token" ; "golite/types" ; "fmt")

type VarScope int

const (GLOBAL VarScope = iota ; LOCAL)

type VarEntry struct {
	*token.Token
	Name  string
	Ty    types.Type
	Scope VarScope
}

type StructEntry struct {
	*token.Token
	Name   string
	Ty     types.Type
	Fields []*VarEntry
}

type FuncEntry struct {
	*token.Token
	Name       string
	ReturnTy   types.Type
	Parameters []*VarEntry
	Variables  *SymbolTable[*VarEntry]
}

func NewVarEntry(name string, ty types.Type, scope VarScope, token *token.Token) *VarEntry {
	return &VarEntry{token, name, ty, scope}
}

func NewStructEntry(name string, ty types.Type, fields []*VarEntry, token *token.Token) *StructEntry {
	return &StructEntry{token, name, ty, fields}
}

func NewFuncEntry(name string, returnTy types.Type, parameters []*VarEntry, variables *SymbolTable[*VarEntry], token *token.Token) *FuncEntry {
	return &FuncEntry{token, name, returnTy, parameters, variables}
}

func (entry *VarEntry) String() string { return fmt.Sprintf("Var: %s", entry.Name) }
func (entry *StructEntry) String() string { return fmt.Sprintf("Struct: %s", entry.Name) }
func (entry *FuncEntry) String() string { return fmt.Sprintf("Func: %s", entry.Name) }

type SymbolTable[T fmt.Stringer] struct {
	Parent *SymbolTable[T]
	Table  map[string]T
}

func NewSymbolTable[T fmt.Stringer](parent *SymbolTable[T]) *SymbolTable[T] {
	return &SymbolTable[T]{parent, make(map[string]T)}
}

func (st *SymbolTable[T]) Insert(id string, entry T) {
	st.Table[id] = entry
}

func (st *SymbolTable[T]) Contains(id string) T {
	var empty T
	if value, exists := st.Table[id]; exists {return value}
	if st.Parent != nil {return st.Parent.Contains(id)}
	return empty
}

type SymbolTables struct {
	Structs *SymbolTable[*StructEntry]
	Globals *SymbolTable[*VarEntry]
	Funcs   *SymbolTable[*FuncEntry]
}

func NewSymbolTables() *SymbolTables {
	return &SymbolTables{NewSymbolTable[*StructEntry](nil), NewSymbolTable[*VarEntry](nil), NewSymbolTable[*FuncEntry](nil)}
}

func SemanticError(errors []*context.CompilerError, line int, column int, msg string) []*context.CompilerError {
	errors = append(errors, &context.CompilerError{
		Line:   line,
		Column: column,
		Msg:    msg,
		Phase:  context.SEMANTIC,
	})
	return errors
}

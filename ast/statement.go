package ast

import ("golite/cfg" ; "golite/context" ; "golite/st" ; "golite/types" ; "golite/llvm")

type Statement interface {
	StatementNode()
	String() string
	TypeCheck([]*context.CompilerError, *st.SymbolTables, *st.FuncEntry, *cfg.Block) ([]*context.CompilerError, types.Type, *cfg.Block)
	TranslateToLLVMStack(*llvm.Block, *st.FuncEntry, *llvm.LLVMProgram) *llvm.Block
}

func (a *Assignment) StatementNode() {}
func (p *Print) StatementNode() {}
func (r *Read) StatementNode() {}
func (d *Delete) StatementNode() {}
func (c *Conditional) StatementNode() {}
func (l *Loop) StatementNode() {}
func (r *Return) StatementNode() {}
func (i *Invocation) StatementNode() {}

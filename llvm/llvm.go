package llvm

import ("golite/codegen" ; "golite/st" ; "golite/types" ; "fmt")

type LLVMOperator int
const (ADD LLVMOperator = iota ; SUB ; MUL ; SDIV ; OR ; AND ; EQ ; NE ; SLT ; SGT ; SLE ; SGE ; XOR)

type LLVMInstruction interface {
	String() string
	MemtoReg(*Block)
	ComputeLiveRange(*LLVMFunction, int)
	TranslateToOutOfSSA(*LLVMFunction, *Block)
	TranslateToAssembly(*LLVMFunction, *Block)
}

type LLVMOperand interface {
	String()        string
	GetType()       types.Type
	PointsToStack() bool
	GetValue()      string
}

type LLVMRegister struct {
	entry         *st.VarEntry
	pointsToStack bool
}

type LLVMImmediate struct {
	value int
	ty    types.Type
}

type LLVMPhysicalRegister struct {
	name string
}

var registerCount int
func init() { registerCount = 1 }

func NewRegisterLabel() string {
	label := fmt.Sprintf("r%d", registerCount)
	registerCount++
	return label
}

func NewRegister(entry *st.VarEntry) *LLVMRegister {
	if entry.Scope == st.LOCAL {
		return &LLVMRegister{entry, true}
	} else {
		return &LLVMRegister{entry, false}
	}
}

func NewImmediate(value int, ty types.Type) *LLVMImmediate {
	return &LLVMImmediate{value, ty}
}

func NewPhysicalRegister(name string) *LLVMPhysicalRegister {
	return &LLVMPhysicalRegister{name}
}

func (register *LLVMRegister) String() string {
	if register.entry.Scope == st.GLOBAL { return fmt.Sprintf("@%s", register.entry.Name) }
	return fmt.Sprintf("%%%s", register.entry.Name)
}

func (immediate *LLVMImmediate) String() string {
	if immediate.ty == types.IntTySig || immediate.ty == types.BoolTySig { return fmt.Sprintf("%d", immediate.value) }
	return "null"
}

func (physicalRegister *LLVMPhysicalRegister) String() string {
	return fmt.Sprintf("%%%s", physicalRegister.name)
}

func (register *LLVMRegister) GetType() types.Type {
	if structTy, ok := register.entry.Ty.(*types.StructTy); ok { return &types.PointerTy{structTy} }
	return register.entry.Ty
}

func (immediate *LLVMImmediate) GetType() types.Type {
	return immediate.ty
}

func (physicalRegister *LLVMPhysicalRegister) GetType() types.Type {
	return types.UnknownTySig
}

func (register *LLVMRegister) PointsToStack() bool {
	return register.pointsToStack
}

func (immediate *LLVMImmediate) PointsToStack() bool {
	return true
}

func (physicalRegister *LLVMPhysicalRegister) PointsToStack() bool {
	return true
}

func (register *LLVMRegister) GetValue() string {
	return register.entry.Name
}

func (immediate *LLVMImmediate) GetValue() string {
	return fmt.Sprintf("#%d", immediate.value)
}

func (physicalRegister *LLVMPhysicalRegister) GetValue() string {
	return physicalRegister.name
}

func TypesToLLVMType(ty types.Type) string {
	switch ty {
		case types.IntTySig: return "i64"
		case types.BoolTySig: return "i64"
		case types.VoidTySig: return "void"
	}

	if ty == types.Int1TySig { return "i1" }
	if pointerTy, ok := ty.(*types.PointerTy); ok {return "%struct." + pointerTy.BaseType.String() + "*" }
	if structTy, ok := ty.(*types.StructTy); ok {return "%struct." + structTy.String() }
	return "i8*"
}

func (operator LLVMOperator) String() string {
	switch operator {
		case ADD: return "add"
		case SUB: return "sub"
		case MUL: return "mul"
		case SDIV: return "sdiv"
		case OR: return "or"
		case AND: return "and"
		case EQ: return "eq"
		case NE: return "ne"
		case SLT: return "slt"
		case SGT: return "sgt"
		case SLE: return "sle"
		case SGE: return "sge"
		case XOR: return "xor"
	}
	panic("Not found operator")
}

func (operator LLVMOperator) Translate() codegen.Operator {
	switch operator {
		case ADD: return codegen.ADD
		case SUB: return codegen.SUB
		case MUL: return codegen.MUL
		case SDIV: return codegen.SDIV
		case OR: return codegen.OR
		case AND: return codegen.AND
		case EQ: return codegen.EQ
		case NE: return codegen.NE
		case SLT: return codegen.SLT
		case SGT: return codegen.SGT
		case SLE: return codegen.SLE
		case SGE: return codegen.SGE
		case XOR: return codegen.XOR
	}
	panic("Not found operator")
}

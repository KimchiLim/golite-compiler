package codegen

type Operator int
const (ADD Operator = iota ; SUB ; MUL ; SDIV ; OR ; AND ; EQ ; NE ; SLT ; SGT ; SLE ; SGE ; XOR)

type Instruction interface {
	String() string
}

func (operator Operator) String() string {
	switch operator {
		case ADD: return "add"
		case SUB: return "sub"
		case MUL: return "mul"
		case SDIV: return "sdiv"
		case OR: return "or"
		case AND: return "and"
		case EQ: return "eq"
		case NE: return "ne"
		case SLT: return "lt"
		case SGT: return "gt"
		case SLE: return "le"
		case SGE: return "ge"
		case XOR: return "eor"
	}
	panic("Not found operator")
}

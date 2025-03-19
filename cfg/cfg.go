package cfg

import ("golite/context" ; "golite/st" ; "golite/token" ; "golite/types" ; "fmt")

var blockCount int
func init() { blockCount = 0 }

func newBlockLabel() string {
	label := fmt.Sprintf("L%d", blockCount)
	blockCount++
	return label
}

type Block struct {
	*token.Token
	Label      string
	Prev       []*Block
	Next       []*Block
	HasReturn  bool
	ReturnType types.Type
}

func NewBlock(token *token.Token) *Block {
	return &Block{token, newBlockLabel(), nil, nil, false, types.VoidTySig}
}

func (block *Block) String() string {
	return fmt.Sprintf("%%%s", block.Label)
}

func (block *Block) AddNext(next *Block) {
	block.Next = append(block.Next, next)
	next.Prev = append(next.Prev, block)
}

func (block *Block) AddReturn(returnType types.Type) {
	block.HasReturn = true
	block.ReturnType = returnType
}

func (block *Block) Validate(errors []*context.CompilerError, returnType types.Type, visited map[string]bool) ([]*context.CompilerError, bool) {
	var result, valid bool
	result = true

    if block.HasReturn && returnType.String() != block.ReturnType.String() {
		errors = st.SemanticError(errors, block.Token.Line, block.Token.Column, fmt.Sprintf("invalid return: expected %s but got %s", returnType, block.ReturnType))
	}

	if block.HasReturn { return errors, true }
	if len(block.Next) == 0 { return errors, false }

	for _, next := range block.Next {
		if _, exists := visited[next.Label]; !exists {
            visited[next.Label] = true
            errors, valid = next.Validate(errors, returnType, visited)
            if !valid { result = false }
        }
	}

	return errors, result
}

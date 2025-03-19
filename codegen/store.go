package codegen

import ("bytes")

type Store struct {
	register string
	address  string
}

func NewStore(register string, address string) *Store {
	return &Store{register, address}
}

func (s *Store) String() string {
	var out bytes.Buffer
	out.WriteString("str")
	out.WriteString(" ")
	out.WriteString(s.register)
	out.WriteString(",")
	out.WriteString(" ")
	out.WriteString("[")
	out.WriteString(s.address)
	out.WriteString("]")
	out.WriteString("\n")
	return out.String()
}

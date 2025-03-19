package types

// Type includes information about working on Types
type Type interface {
	String() string
}

type IntTy struct {}
type BoolTy struct {}
type VoidTy struct {}
type StructTy struct { Name string }
type PointerTy struct { BaseType Type }
type UnknownTy struct {}

func (intTy *IntTy) String() string { return "int" }
func (boolTy *BoolTy) String() string { return "bool" }
func (voidTy *VoidTy) String() string { return "void" }
func (structTy *StructTy) String() string { return structTy.Name }
func (pointerTy *PointerTy) String() string { return "*" + pointerTy.BaseType.String() }
func (unknownTy *UnknownTy) String() string { return "unknown" }

var Int1TySig *IntTy
var IntTySig *IntTy
var BoolTySig *BoolTy
var VoidTySig *VoidTy
var NilTySig *StructTy
var UnknownTySig *UnknownTy

// The init() function will only be called once per package. This is where you can setup singletons for types
func init() {
	Int1TySig = &IntTy{}
	IntTySig = &IntTy{}
	BoolTySig = &BoolTy{}
	VoidTySig = &VoidTy{}
	NilTySig = &StructTy{"nil"}
	UnknownTySig = &UnknownTy{}
}

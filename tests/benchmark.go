package tests

import (
	"math/big"

	"gopkg.in/inf.v0"
)

type BenchmarkStruct struct {
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Bool    bool
	String  string
	Float32 float32
	Float64 float64
	BigInt  big.Int
	Dec     inf.Dec

	NamedInt     NamedInt
	NamedInt8    NamedInt8
	NamedInt16   NamedInt16
	NamedInt32   NamedInt32
	NamedInt64   NamedInt64
	NamedUint    NamedUint
	NamedUint8   NamedUint8
	NamedUint16  NamedUint16
	NamedUint32  NamedUint32
	NamedUint64  NamedUint64
	NamedBool    NamedBool
	NamedString  NamedString
	NamedFloat32 NamedFloat32
	NamedDouble  NamedFloat64
}

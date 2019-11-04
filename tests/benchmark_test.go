package tests

import (
	"math/big"
	"testing"

	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"

	"github.com/kiwicom/easycql/marshal"
)

var benchmarkUnmarshalOut BenchmarkStruct
var benchmarkMarshalOut []byte

var benchTypeInfo = gocql.UDTTypeInfo{
	NativeType: gocql.NewNativeType(3, gocql.TypeUDT, ""),
	KeySpace:   "ks",
	Name:       "MyUDT",
	Elements: []gocql.UDTField{
		{
			Name: "Int",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "Int8",
			Type: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		},
		{
			Name: "Int16",
			Type: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		},
		{
			Name: "Int32",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "Int64",
			Type: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		},
		{
			Name: "Uint",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "Uint8",
			Type: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		},
		{
			Name: "Uint16",
			Type: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		},
		{
			Name: "Uint32",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "Uint64",
			Type: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		},
		{
			Name: "Bool",
			Type: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		},
		{
			Name: "String",
			Type: gocql.NewNativeType(2, gocql.TypeVarchar, ""),
		},
		{
			Name: "Float32",
			Type: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		},
		{
			Name: "Float64",
			Type: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		},
		{
			Name: "BigInt",
			Type: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		},
		{
			Name: "NamedInt",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "NamedInt8",
			Type: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		},
		{
			Name: "NamedInt16",
			Type: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		},
		{
			Name: "NamedInt32",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "NamedInt64",
			Type: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		},
		{
			Name: "NamedUint",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "NamedUint8",
			Type: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		},
		{
			Name: "NamedUint16",
			Type: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		},
		{
			Name: "NamedUint32",
			Type: gocql.NewNativeType(2, gocql.TypeInt, ""),
		},
		{
			Name: "NamedUint64",
			Type: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		},
		{
			Name: "NamedBool",
			Type: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		},
		{
			Name: "NamedString",
			Type: gocql.NewNativeType(2, gocql.TypeVarchar, ""),
		},
		{
			Name: "NamedFloat32",
			Type: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		},
		{
			Name: "NamedFloat64",
			Type: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		},
	},
}

func BenchmarkUnmarshal(b *testing.B) {
	fieldValues := [][]byte{
		[]byte("\x12\x34\x56\x78"),                 // int
		[]byte("\x12"),                             // int8
		[]byte("\x12\x34"),                         // int16
		[]byte("\x12\x34\x56\x78"),                 // int32
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"), // int64
		[]byte("\x12\x34\x56\x78"),                 // uint
		[]byte("\x12"),                             // uint8
		[]byte("\x12\x34"),                         // uint16
		[]byte("\x12\x34\x56\x78"),                 // uint32
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"), // uint64
		[]byte("\x01"),                             // bool
		[]byte("hello world"),                      // string
		[]byte("\x00\x00\x00\x00"),                 // float32
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00"), // float64
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef\x12\x34\x56\x78\x9a"), // BigInt
		[]byte("\x00\x00\x00\x01\xde\xef\x12\x34\x56\x78\x9a"),         // Dec
		[]byte("\x12\x34\x56\x78"),                                     // NamedInt
		[]byte("\x12"),                                                 // NamedInt8
		[]byte("\x12\x34"),                                             // NamedInt16
		[]byte("\x12\x34\x56\x78"),                                     // NamedInt32
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"),                     // NamedInt64
		[]byte("\x12\x34\x56\x78"),                                     // NamedUint
		[]byte("\x12"),                                                 // NamedUint8
		[]byte("\x12\x34"),                                             // NamedUint16
		[]byte("\x12\x34\x56\x78"),                                     // NamedUint32
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"),                     // NamedUint64
		[]byte("\x01"),                                                 // NamedBool
		[]byte("hello world"),                                          // NamedString
		[]byte("\x00\x00\x00\x00"),                                     // NamedFloat32
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00"),                     // NamedFloat64
	}

	var data []byte
	for fieldIdx := range fieldValues {
		data = marshal.AppendBytes(data, fieldValues[fieldIdx])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := gocql.Unmarshal(benchTypeInfo, data, &benchmarkUnmarshalOut)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	s := BenchmarkStruct{
		Int:          0x12345678,
		Int8:         0x12,
		Int16:        0x1234,
		Int32:        0x12345678,
		Int64:        0x123456789abcdeef,
		Uint:         0x12345678,
		Uint8:        0x12,
		Uint16:       0x1234,
		Uint32:       0x12345678,
		Uint64:       0x123456789abcdeef,
		Bool:         true,
		String:       "hello world",
		Float32:      0,
		Float64:      0,
		BigInt:       *big.NewInt(0x123456789abcdeef),
		Dec:          *inf.NewDec(0x123456789abcdeef, 1),
		NamedInt:     0x12345678,
		NamedInt8:    0x12,
		NamedInt16:   0x1234,
		NamedInt32:   0x12345678,
		NamedInt64:   0x123456789abcdeef,
		NamedUint:    0x12345678,
		NamedUint8:   0x12,
		NamedUint16:  0x1234,
		NamedUint32:  0x12345678,
		NamedUint64:  0x123456789abcdeef,
		NamedBool:    true,
		NamedString:  "hello world",
		NamedFloat32: 0,
		NamedDouble:  0,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := gocql.Marshal(benchTypeInfo, s)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMarshalOut = data
	}
}

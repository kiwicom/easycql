package tests

import (
	"testing"

	"github.com/gocql/gocql"

	"github.com/kiwicom/easycql/marshal"
)

var benchmarkUnmarshalOut BenchmarkStruct

func BenchmarkUnmarshal(b *testing.B) {
	info := gocql.UDTTypeInfo{
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

	fieldValues := [][]byte{
		[]byte("\x12\x34\x56\x78"),
		[]byte("\x12"),
		[]byte("\x12\x34"),
		[]byte("\x12\x34\x56\x78"),
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"),
		[]byte("\x01"),
		[]byte("hello world"),
		[]byte("\x00\x00\x00\x00"),
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef\x12\x34\x56\x78\x9a"),
		[]byte("\x00\x00\x00\x01\xde\xef\x12\x34\x56\x78\x9a"),
		[]byte("\x12\x34\x56\x78"),
		[]byte("\x12"),
		[]byte("\x12\x34"),
		[]byte("\x12\x34\x56\x78"),
		[]byte("\x12\x34\x56\x78\x9a\xbc\xde\xef"),
		[]byte("\x01"),
		[]byte("hello world"),
		[]byte("\x00\x00\x00\x00"),
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
	}

	var data []byte
	for fieldIdx := range fieldValues {
		data = marshal.AppendBytes(data, fieldValues[fieldIdx])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := gocql.Unmarshal(info, data, &benchmarkUnmarshalOut)
		if err != nil {
			b.Fatal(err)
		}
	}
}

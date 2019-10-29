package tests

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/require"
	"gopkg.in/inf.v0"

	"github.com/kiwicom/easycql/marshal"
)

var varcharTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Values        CQLVarcharTypes
}{
	{
		Name:          "varchar value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarchar, ""),
		Data:          []byte("hello world"),
		Values: CQLVarcharTypes{
			String:          "hello world",
			StringPtr:       newStringPtr("hello world"),
			NamedString:     NamedString("hello world"),
			NamedStringPtr:  newNamedStringPtr("hello world"),
			Bytes:           []byte("hello world"),
			BytesPtr:        newBytesPtr([]byte("hello world")),
			NamedBytes:      NamedBytes("hello world"),
			NamedBytesPtr:   newNamedBytesPtr(NamedBytes("hello world")),
			CustomString:    CustomString("HELLO WORLD"),
			CustomStringPtr: newCustomStringPtr(CustomString("HELLO WORLD")),
		},
	},
	{
		Name:          "varchar null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarchar, ""),
		Data:          []byte(nil),
		Values: CQLVarcharTypes{
			String:          "",
			StringPtr:       nil,
			NamedString:     NamedString(""),
			NamedStringPtr:  nil,
			Bytes:           nil,
			BytesPtr:        nil,
			NamedBytes:      NamedBytes(nil),
			NamedBytesPtr:   nil,
			CustomString:    CustomString(""),
			CustomStringPtr: nil,
		},
	},
	{
		Name:          "ascii value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeAscii, ""),
		Data:          []byte("hello world"),
		Values: CQLVarcharTypes{
			String:          "hello world",
			StringPtr:       newStringPtr("hello world"),
			NamedString:     NamedString("hello world"),
			NamedStringPtr:  newNamedStringPtr("hello world"),
			Bytes:           []byte("hello world"),
			BytesPtr:        newBytesPtr([]byte("hello world")),
			NamedBytes:      NamedBytes("hello world"),
			NamedBytesPtr:   newNamedBytesPtr(NamedBytes("hello world")),
			CustomString:    CustomString("HELLO WORLD"),
			CustomStringPtr: newCustomStringPtr(CustomString("HELLO WORLD")),
		},
	},
	{
		Name:          "ascii null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeAscii, ""),
		Data:          []byte(nil),
		Values: CQLVarcharTypes{
			String:          "",
			StringPtr:       nil,
			NamedString:     NamedString(""),
			NamedStringPtr:  nil,
			Bytes:           nil,
			BytesPtr:        nil,
			NamedBytes:      NamedBytes(nil),
			NamedBytesPtr:   nil,
			CustomString:    CustomString(""),
			CustomStringPtr: nil,
		},
	},
	{
		Name:          "text value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeText, ""),
		Data:          []byte("hello world"),
		Values: CQLVarcharTypes{
			String:          "hello world",
			StringPtr:       newStringPtr("hello world"),
			NamedString:     NamedString("hello world"),
			NamedStringPtr:  newNamedStringPtr("hello world"),
			Bytes:           []byte("hello world"),
			BytesPtr:        newBytesPtr([]byte("hello world")),
			NamedBytes:      NamedBytes("hello world"),
			NamedBytesPtr:   newNamedBytesPtr(NamedBytes("hello world")),
			CustomString:    CustomString("HELLO WORLD"),
			CustomStringPtr: newCustomStringPtr(CustomString("HELLO WORLD")),
		},
	},
	{
		Name:          "text null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeText, ""),
		Data:          []byte(nil),
		Values: CQLVarcharTypes{
			String:          "",
			StringPtr:       nil,
			NamedString:     NamedString(""),
			NamedStringPtr:  nil,
			Bytes:           nil,
			BytesPtr:        nil,
			NamedBytes:      NamedBytes(nil),
			NamedBytesPtr:   nil,
			CustomString:    CustomString(""),
			CustomStringPtr: nil,
		},
	},
	{
		Name:          "blob value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBlob, ""),
		Data:          []byte("hello world\x00abc"),
		Values: CQLVarcharTypes{
			String:          "hello world\x00abc",
			StringPtr:       newStringPtr("hello world\x00abc"),
			NamedString:     NamedString("hello world\x00abc"),
			NamedStringPtr:  newNamedStringPtr("hello world\x00abc"),
			Bytes:           []byte("hello world\x00abc"),
			BytesPtr:        newBytesPtr([]byte("hello world\x00abc")),
			NamedBytes:      NamedBytes("hello world\x00abc"),
			NamedBytesPtr:   newNamedBytesPtr(NamedBytes("hello world\x00abc")),
			CustomString:    CustomString("HELLO WORLD\x00ABC"),
			CustomStringPtr: newCustomStringPtr(CustomString("HELLO WORLD\x00ABC")),
		},
	},
	{
		Name:          "blob null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBlob, ""),
		Data:          []byte(nil),
		Values: CQLVarcharTypes{
			String:          "",
			StringPtr:       nil,
			NamedString:     NamedString(""),
			NamedStringPtr:  nil,
			Bytes:           nil,
			BytesPtr:        nil,
			NamedBytes:      NamedBytes(nil),
			NamedBytesPtr:   nil,
			CustomString:    CustomString(""),
			CustomStringPtr: nil,
		},
	},
}

func buildUDT(typ reflect.Type, fieldType gocql.TypeInfo, fieldData []byte) (gocql.UDTTypeInfo, []byte) {
	var data []byte
	typeInfo := gocql.UDTTypeInfo{
		NativeType: gocql.NewNativeType(3, gocql.TypeUDT, ""),
		KeySpace:   "myKeyspace",
		Name:       "CQLVarcharTypesUDT",
		Elements:   make([]gocql.UDTField, 0, typ.NumField()),
	}
	for i := 0; i < typ.NumField(); i++ {
		udtField := gocql.UDTField{
			Name: typ.Field(i).Name,
			Type: fieldType,
		}
		typeInfo.Elements = append(typeInfo.Elements, udtField)
		data = marshal.AppendBytes(data, fieldData)
	}
	return typeInfo, data
}

func TestUnmarshalVarchar(t *testing.T) {
	for _, test := range varcharTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.TypeOf((*CQLVarcharTypes)(nil)).Elem(), test.FieldTypeInfo, test.Data)
			var value CQLVarcharTypes
			err := gocql.Unmarshal(typeInfo, data, &value)
			require.NoError(t, err)
			require.Equal(t, test.Values, value)
		})
	}
}

var integerTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// ------------------------------- int → integers -------------------------------------
	// int → int
	{
		Name:          "int to int zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt{Int: 0},
	},
	{
		Name:          "int to int byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleInt{Int: 0x01020304},
	},
	{
		Name:          "int to int MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleInt{Int: math.MinInt32},
	},
	{
		Name:          "int to int MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleInt{Int: math.MaxInt32},
	},
	{
		Name:          "int to int null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt{Int: 0},
	},
	{
		Name:          "int to int empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt{Int: 0},
	},
	// int → *int
	{
		Name:          "int to intPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleIntPtr{IntPtr: newIntPtr(0)},
	},
	{
		Name:          "int to intPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleIntPtr{IntPtr: newIntPtr(0x01020304)},
	},
	{
		Name:          "int to intPtr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleIntPtr{IntPtr: newIntPtr(math.MinInt32)},
	},
	{
		Name:          "int to intPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleIntPtr{IntPtr: newIntPtr(math.MaxInt32)},
	},
	{
		Name:          "int to intPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleIntPtr{IntPtr: nil},
	},
	{
		Name:          "int to intPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleIntPtr{IntPtr: newIntPtr(0)},
	},
	// int → int8
	{
		Name:          "int to int8 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt8{Int8: 0},
	},
	{
		Name:          "int to int8 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleInt8{Int8: 0x01},
	},
	{
		Name:          "int to int8 MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleInt8{Int8: math.MinInt8},
	},
	{
		Name:          "int to int8 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleInt8{},
		Error:         true,
	},
	{
		Name:          "int to int8 MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleInt8{Int8: math.MaxInt8},
	},
	{
		Name:          "int to int8 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleInt8{},
		Error:         true,
	},
	{
		Name:          "int to int8 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt8{Int8: 0},
	},
	{
		Name:          "int to int8 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt8{Int8: 0},
	},
	// int → *int8
	{
		Name:          "int to int8Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt8Ptr{Int8Ptr: newInt8Ptr(0)},
	},
	{
		Name:          "int to int8Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleInt8Ptr{Int8Ptr: newInt8Ptr(0x01)},
	},
	{
		Name:          "int to int8Ptr MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleInt8Ptr{Int8Ptr: newInt8Ptr(math.MinInt8)},
	},
	{
		Name:          "int to int8Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to int8Ptr MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleInt8Ptr{Int8Ptr: newInt8Ptr(math.MaxInt8)},
	},
	{
		Name:          "int to int8Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to int8Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt8Ptr{Int8Ptr: nil},
	},
	{
		Name:          "int to int8Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt8Ptr{Int8Ptr: newInt8Ptr(0)},
	},
	// int → int16
	{
		Name:          "int to int16 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt16{Int16: 0},
	},
	{
		Name:          "int to int16 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleInt16{Int16: 0x0102},
	},
	{
		Name:          "int to int16 MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleInt16{Int16: math.MinInt16},
	},
	{
		Name:          "int to int16 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleInt16{},
		Error:         true,
	},
	{
		Name:          "int to int16 MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleInt16{Int16: math.MaxInt16},
	},
	{
		Name:          "int to int16 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleInt16{},
		Error:         true,
	},
	{
		Name:          "int to int16 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt16{Int16: 0},
	},
	{
		Name:          "int to int16 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt16{Int16: 0},
	},
	// int → *int16
	{
		Name:          "int to int16Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt16Ptr{Int16Ptr: newInt16Ptr(0)},
	},
	{
		Name:          "int to int16Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleInt16Ptr{Int16Ptr: newInt16Ptr(0x0102)},
	},
	{
		Name:          "int to int16Ptr MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleInt16Ptr{Int16Ptr: newInt16Ptr(math.MinInt16)},
	},
	{
		Name:          "int to int16Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to int16Ptr MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleInt16Ptr{Int16Ptr: newInt16Ptr(math.MaxInt16)},
	},
	{
		Name:          "int to int16Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to int16Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt16Ptr{Int16Ptr: nil},
	},
	{
		Name:          "int to int16Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt16Ptr{Int16Ptr: newInt16Ptr(0)},
	},
	// int → int32
	{
		Name:          "int to int32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt32{Int32: 0},
	},
	{
		Name:          "int to int32 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleInt32{Int32: 0x01020304},
	},
	{
		Name:          "int to int32 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleInt32{Int32: math.MinInt32},
	},
	{
		Name:          "int to int32 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleInt32{Int32: math.MaxInt32},
	},
	{
		Name:          "int to int32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt32{Int32: 0},
	},
	{
		Name:          "int to int32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt32{Int32: 0},
	},
	// int → *int32
	{
		Name:          "int to int32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt32Ptr{Int32Ptr: newInt32Ptr(0)},
	},
	{
		Name:          "int to int32Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleInt32Ptr{Int32Ptr: newInt32Ptr(0x01020304)},
	},
	{
		Name:          "int to int32Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleInt32Ptr{Int32Ptr: newInt32Ptr(math.MinInt32)},
	},
	{
		Name:          "int to int32Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleInt32Ptr{Int32Ptr: newInt32Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to int32Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt32Ptr{Int32Ptr: nil},
	},
	{
		Name:          "int to int32Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt32Ptr{Int32Ptr: newInt32Ptr(0)},
	},
	// int → int64
	{
		Name:          "int to int64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt64{Int64: 0},
	},
	{
		Name:          "int to int64 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleInt64{Int64: 0x01020304},
	},
	{
		Name:          "int to int64 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleInt64{Int64: math.MinInt32},
	},
	{
		Name:          "int to int64 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleInt64{Int64: math.MaxInt32},
	},
	{
		Name:          "int to int64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt64{Int64: 0},
	},
	{
		Name:          "int to int64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt64{Int64: 0},
	},
	// int → *int64
	{
		Name:          "int to int64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleInt64Ptr{Int64Ptr: newInt64Ptr(0)},
	},
	{
		Name:          "int to int64Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleInt64Ptr{Int64Ptr: newInt64Ptr(0x01020304)},
	},
	{
		Name:          "int to int64Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleInt64Ptr{Int64Ptr: newInt64Ptr(math.MinInt32)},
	},
	{
		Name:          "int to int64Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleInt64Ptr{Int64Ptr: newInt64Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to int64Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleInt64Ptr{Int64Ptr: nil},
	},
	{
		Name:          "int to int64Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleInt64Ptr{Int64Ptr: newInt64Ptr(0)},
	},
	// ------------------------------- int → named integers -------------------------------------
	// int → NamedInt
	{
		Name:          "int to NamedInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	{
		Name:          "int to NamedInt byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt{NamedInt: 0x01020304},
	},
	{
		Name:          "int to NamedInt MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt{NamedInt: math.MinInt32},
	},
	{
		Name:          "int to NamedInt MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt{NamedInt: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	{
		Name:          "int to NamedInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	// int → *NamedInt
	{
		Name:          "int to NamedIntPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0)},
	},
	{
		Name:          "int to NamedIntPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0x01020304)},
	},
	{
		Name:          "int to NamedIntPtr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(math.MinInt32)},
	},
	{
		Name:          "int to NamedIntPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedIntPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedIntPtr{NamedIntPtr: nil},
	},
	{
		Name:          "int to NamedIntPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0)},
	},
	// int → int8
	{
		Name:          "int to NamedInt8 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	{
		Name:          "int to NamedInt8 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedInt8{NamedInt8: 0x01},
	},
	{
		Name:          "int to NamedInt8 MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleNamedInt8{NamedInt8: math.MinInt8},
	},
	{
		Name:          "int to NamedInt8 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleNamedInt8{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8 MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedInt8{NamedInt8: math.MaxInt8},
	},
	{
		Name:          "int to NamedInt8 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleNamedInt8{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	{
		Name:          "int to NamedInt8 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	// int → *NamedInt8
	{
		Name:          "int to NamedInt8Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0)},
	},
	{
		Name:          "int to NamedInt8Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0x01)},
	},
	{
		Name:          "int to NamedInt8Ptr MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(math.MinInt8)},
	},
	{
		Name:          "int to NamedInt8Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleNamedInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8Ptr MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(math.MaxInt8)},
	},
	{
		Name:          "int to NamedInt8Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleNamedInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: nil},
	},
	{
		Name:          "int to NamedInt8Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0)},
	},
	// int → int16
	{
		Name:          "int to NamedInt16 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	{
		Name:          "int to NamedInt16 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedInt16{NamedInt16: 0x0102},
	},
	{
		Name:          "int to NamedInt16 MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleNamedInt16{NamedInt16: math.MinInt16},
	},
	{
		Name:          "int to NamedInt16 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleNamedInt16{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16 MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedInt16{NamedInt16: math.MaxInt16},
	},
	{
		Name:          "int to NamedInt16 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleNamedInt16{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	{
		Name:          "int to NamedInt16 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	// int → *NamedInt16
	{
		Name:          "int to NamedInt16Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0)},
	},
	{
		Name:          "int to NamedInt16Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0x0102)},
	},
	{
		Name:          "int to NamedInt16Ptr MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(math.MinInt16)},
	},
	{
		Name:          "int to NamedInt16Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleNamedInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16Ptr MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(math.MaxInt16)},
	},
	{
		Name:          "int to NamedInt16Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleNamedInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: nil},
	},
	{
		Name:          "int to NamedInt16Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0)},
	},
	// int → int32
	{
		Name:          "int to NamedInt32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	{
		Name:          "int to NamedInt32 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt32{NamedInt32: 0x01020304},
	},
	{
		Name:          "int to NamedInt32 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt32{NamedInt32: math.MinInt32},
	},
	{
		Name:          "int to NamedInt32 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt32{NamedInt32: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	{
		Name:          "int to NamedInt32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	// int → *NamedInt32
	{
		Name:          "int to NamedInt32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0)},
	},
	{
		Name:          "int to NamedInt32Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedInt32Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(math.MinInt32)},
	},
	{
		Name:          "int to NamedInt32Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedInt32Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: nil},
	},
	{
		Name:          "int to NamedInt32Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0)},
	},
	// int → int64
	{
		Name:          "int to NamedInt64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	{
		Name:          "int to NamedInt64 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt64{NamedInt64: 0x01020304},
	},
	{
		Name:          "int to NamedInt64 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt64{NamedInt64: math.MinInt32},
	},
	{
		Name:          "int to NamedInt64 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt64{NamedInt64: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	{
		Name:          "int to NamedInt64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	// int → *NamedInt64
	{
		Name:          "int to NamedInt64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0)},
	},
	{
		Name:          "int to NamedInt64Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedInt64Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(math.MinInt32)},
	},
	{
		Name:          "int to NamedInt64Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedInt64Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: nil},
	},
	{
		Name:          "int to NamedInt64Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0)},
	},
	// ------------------------------- int → unsigned integers -----------------------------
	// int → uint
	{
		Name:          "int to uint zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint{Uint: 0},
	},
	{
		Name:          "int to uint byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUint{Uint: 0x01020304},
	},
	{
		Name:          "int to uint MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUint{Uint: math.MaxInt32},
	},
	{
		Name:          "int to uint MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint{Uint: math.MaxUint32},
	},
	{
		Name:          "int to uint null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint{Uint: 0},
	},
	{
		Name:          "int to uint empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint{Uint: 0},
	},
	// int → *uint
	{
		Name:          "int to uintPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUintPtr{UintPtr: newUintPtr(0)},
	},
	{
		Name:          "int to uintPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUintPtr{UintPtr: newUintPtr(0x01020304)},
	},
	{
		Name:          "int to uintPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUintPtr{UintPtr: newUintPtr(math.MaxInt32)},
	},
	{
		Name:          "int to uintPtr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUintPtr{UintPtr: newUintPtr(math.MaxUint32)},
	},
	{
		Name:          "int to uintPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUintPtr{UintPtr: nil},
	},
	{
		Name:          "int to uintPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUintPtr{UintPtr: newUintPtr(0)},
	},
	// int → uint8
	{
		Name:          "int to uint8 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint8{Uint8: 0},
	},
	{
		Name:          "int to uint8 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleUint8{Uint8: 0x01},
	},
	{
		Name:          "int to uint8 MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleUint8{Uint8: math.MaxInt8},
	},
	{
		Name:          "int to uint8 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint8{},
		Error:         true,
	},
	{
		Name:          "int to uint8 MaxUint8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\xff"),
		Value:         SingleUint8{Uint8: math.MaxUint8},
	},
	{
		Name:          "int to uint8 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x00"),
		Value:         SingleUint8{},
		Error:         true,
	},
	{
		Name:          "int to uint8 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint8{Uint8: 0},
	},
	{
		Name:          "int to uint8 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint8{Uint8: 0},
	},
	// int → *uint8
	{
		Name:          "int to uint8Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint8Ptr{Uint8Ptr: newUint8Ptr(0)},
	},
	{
		Name:          "int to uint8Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleUint8Ptr{Uint8Ptr: newUint8Ptr(0x01)},
	},
	{
		Name:          "int to uint8Ptr MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleUint8Ptr{Uint8Ptr: newUint8Ptr(math.MaxInt8)},
	},
	{
		Name:          "int to uint8Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to uint8Ptr MaxUint8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\xff"),
		Value:         SingleUint8Ptr{Uint8Ptr: newUint8Ptr(math.MaxUint8)},
	},
	{
		Name:          "int to uint8Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x00"),
		Value:         SingleUint8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to uint8Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint8Ptr{Uint8Ptr: nil},
	},
	{
		Name:          "int to uint8Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint8Ptr{Uint8Ptr: newUint8Ptr(0)},
	},
	// int → uint16
	{
		Name:          "int to uint16 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint16{Uint16: 0},
	},
	{
		Name:          "int to uint16 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleUint16{Uint16: 0x0102},
	},
	{
		Name:          "int to uint16 MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleUint16{Uint16: math.MaxInt16},
	},
	{
		Name:          "int to uint16 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint16{},
		Error:         true,
	},
	{
		Name:          "int to uint16 MaxUint16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\xff\xff"),
		Value:         SingleUint16{Uint16: math.MaxUint16},
	},
	{
		Name:          "int to uint16 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x01\x00\x00"),
		Value:         SingleUint16{},
		Error:         true,
	},
	{
		Name:          "int to uint16 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint16{Uint16: 0},
	},
	{
		Name:          "int to uint16 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint16{Uint16: 0},
	},
	// int → *uint16
	{
		Name:          "int to uint16Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint16Ptr{Uint16Ptr: newUint16Ptr(0)},
	},
	{
		Name:          "int to uint16Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleUint16Ptr{Uint16Ptr: newUint16Ptr(0x0102)},
	},
	{
		Name:          "int to uint16Ptr MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleUint16Ptr{Uint16Ptr: newUint16Ptr(math.MaxInt16)},
	},
	{
		Name:          "int to uint16Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to uint16Ptr MaxUint16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\xff\xff"),
		Value:         SingleUint16Ptr{Uint16Ptr: newUint16Ptr(math.MaxUint16)},
	},
	{
		Name:          "int to uint16Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x01\x00\x00"),
		Value:         SingleInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to uint16Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint16Ptr{Uint16Ptr: nil},
	},
	{
		Name:          "int to uint16Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint16Ptr{Uint16Ptr: newUint16Ptr(0)},
	},
	// int → uint32
	{
		Name:          "int to uint32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint32{Uint32: 0},
	},
	{
		Name:          "int to uint32 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUint32{Uint32: 0x01020304},
	},
	{
		Name:          "int to uint32 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleUint32{Uint32: 0x80000000},
	},
	{
		Name:          "int to int32 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUint32{Uint32: math.MaxInt32},
	},
	{
		Name:          "int to int32 MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint32{Uint32: math.MaxUint32},
	},
	{
		Name:          "int to uint32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint32{Uint32: 0},
	},
	{
		Name:          "int to uint32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint32{Uint32: 0},
	},
	// int → *uint32
	{
		Name:          "int to uint32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(0)},
	},
	{
		Name:          "int to uint32Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(0x01020304)},
	},
	{
		Name:          "int to uint32Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(0x80000000)},
	},
	{
		Name:          "int to uint32Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to uint32Ptr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(math.MaxUint32)},
	},
	{
		Name:          "int to uint32Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint32Ptr{Uint32Ptr: nil},
	},
	{
		Name:          "int to uint32Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint32Ptr{Uint32Ptr: newUint32Ptr(0)},
	},
	// int → uint64
	{
		Name:          "int to int64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint64{Uint64: 0},
	},
	{
		Name:          "int to uint64 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUint64{Uint64: 0x01020304},
	},
	{
		Name:          "int to uint64 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleUint64{Uint64: 0x80000000},
	},
	{
		Name:          "int to uint64 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUint64{Uint64: math.MaxInt32},
	},
	{
		Name:          "int to uint64 MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint64{Uint64: math.MaxUint32},
	},
	{
		Name:          "int to uint64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint64{Uint64: 0},
	},
	{
		Name:          "int to uint64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint64{Uint64: 0},
	},
	// int → *uint64
	{
		Name:          "int to uint64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(0)},
	},
	{
		Name:          "int to uint64Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(0x01020304)},
	},
	{
		Name:          "int to uint64Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(0x80000000)},
	},
	{
		Name:          "int to uint64Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to uint64Ptr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(math.MaxUint32)},
	},
	{
		Name:          "int to uint64Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleUint64Ptr{Uint64Ptr: nil},
	},
	{
		Name:          "int to uint64Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleUint64Ptr{Uint64Ptr: newUint64Ptr(0)},
	},
	// ------------------------------- int → named integers -------------------------------------
	// int → NamedInt
	{
		Name:          "int to NamedInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	{
		Name:          "int to NamedInt byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt{NamedInt: 0x01020304},
	},
	{
		Name:          "int to NamedInt MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt{NamedInt: math.MinInt32},
	},
	{
		Name:          "int to NamedInt MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt{NamedInt: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	{
		Name:          "int to NamedInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt{NamedInt: 0},
	},
	// int → *NamedInt
	{
		Name:          "int to NamedIntPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0)},
	},
	{
		Name:          "int to NamedIntPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0x01020304)},
	},
	{
		Name:          "int to NamedIntPtr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(math.MinInt32)},
	},
	{
		Name:          "int to NamedIntPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedIntPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedIntPtr{NamedIntPtr: nil},
	},
	{
		Name:          "int to NamedIntPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedIntPtr{NamedIntPtr: newNamedIntPtr(0)},
	},
	// int → int8
	{
		Name:          "int to NamedInt8 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	{
		Name:          "int to NamedInt8 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedInt8{NamedInt8: 0x01},
	},
	{
		Name:          "int to NamedInt8 MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleNamedInt8{NamedInt8: math.MinInt8},
	},
	{
		Name:          "int to NamedInt8 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleNamedInt8{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8 MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedInt8{NamedInt8: math.MaxInt8},
	},
	{
		Name:          "int to NamedInt8 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleNamedInt8{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	{
		Name:          "int to NamedInt8 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt8{NamedInt8: 0},
	},
	// int → *NamedInt8
	{
		Name:          "int to NamedInt8Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0)},
	},
	{
		Name:          "int to NamedInt8Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0x01)},
	},
	{
		Name:          "int to NamedInt8Ptr MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x80"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(math.MinInt8)},
	},
	{
		Name:          "int to NamedInt8Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\x7f"),
		Value:         SingleNamedInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8Ptr MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(math.MaxInt8)},
	},
	{
		Name:          "int to NamedInt8Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x80"),
		Value:         SingleNamedInt8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt8Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: nil},
	},
	{
		Name:          "int to NamedInt8Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt8Ptr{NamedInt8Ptr: newNamedInt8Ptr(0)},
	},
	// int → NamedInt16
	{
		Name:          "int to NamedInt16 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	{
		Name:          "int to NamedInt16 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedInt16{NamedInt16: 0x0102},
	},
	{
		Name:          "int to NamedInt16 MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleNamedInt16{NamedInt16: math.MinInt16},
	},
	{
		Name:          "int to NamedInt16 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleNamedInt16{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16 MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedInt16{NamedInt16: math.MaxInt16},
	},
	{
		Name:          "int to NamedInt16 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleNamedInt16{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	{
		Name:          "int to NamedInt16 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt16{NamedInt16: 0},
	},
	// int → *NamedInt16
	{
		Name:          "int to NamedInt16Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0)},
	},
	{
		Name:          "int to NamedInt16Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0x0102)},
	},
	{
		Name:          "int to NamedInt16Ptr MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x80\x00"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(math.MinInt16)},
	},
	{
		Name:          "int to NamedInt16Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\x7f\xff"),
		Value:         SingleNamedInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16Ptr MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(math.MaxInt16)},
	},
	{
		Name:          "int to NamedInt16Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x80\x00"),
		Value:         SingleNamedInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedInt16Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: nil},
	},
	{
		Name:          "int to NamedInt16Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt16Ptr{NamedInt16Ptr: newNamedInt16Ptr(0)},
	},
	// int → int32
	{
		Name:          "int to NamedInt32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	{
		Name:          "int to NamedInt32 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt32{NamedInt32: 0x01020304},
	},
	{
		Name:          "int to NamedInt32 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt32{NamedInt32: math.MinInt32},
	},
	{
		Name:          "int to NamedInt32 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt32{NamedInt32: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	{
		Name:          "int to NamedInt32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt32{NamedInt32: 0},
	},
	// int → *NamedInt32
	{
		Name:          "int to NamedInt32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0)},
	},
	{
		Name:          "int to NamedInt32Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedInt32Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(math.MinInt32)},
	},
	{
		Name:          "int to NamedInt32Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedInt32Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: nil},
	},
	{
		Name:          "int to NamedInt32Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt32Ptr{NamedInt32Ptr: newNamedInt32Ptr(0)},
	},
	// int → int64
	{
		Name:          "int to NamedInt64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	{
		Name:          "int to NamedInt64 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt64{NamedInt64: 0x01020304},
	},
	{
		Name:          "int to NamedInt64 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt64{NamedInt64: math.MinInt32},
	},
	{
		Name:          "int to NamedInt64 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt64{NamedInt64: math.MaxInt32},
	},
	{
		Name:          "int to NamedInt64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	{
		Name:          "int to NamedInt64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt64{NamedInt64: 0},
	},
	// int → *NamedInt64
	{
		Name:          "int to NamedInt64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0)},
	},
	{
		Name:          "int to NamedInt64Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedInt64Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(math.MinInt32)},
	},
	{
		Name:          "int to NamedInt64Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedInt64Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: nil},
	},
	{
		Name:          "int to NamedInt64Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedInt64Ptr{NamedInt64Ptr: newNamedInt64Ptr(0)},
	},

	// int → *NamedUint
	{
		Name:          "int to NamedUintPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUintPtr{NamedUintPtr: newNamedUintPtr(0)},
	},
	{
		Name:          "int to NamedUintPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedUintPtr{NamedUintPtr: newNamedUintPtr(0x01020304)},
	},
	{
		Name:          "int to NamedUintPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedUintPtr{NamedUintPtr: newNamedUintPtr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedUintPtr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUintPtr{NamedUintPtr: newNamedUintPtr(math.MaxUint32)},
	},
	{
		Name:          "int to NamedUintPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUintPtr{NamedUintPtr: nil},
	},
	{
		Name:          "int to NamedUintPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUintPtr{NamedUintPtr: newNamedUintPtr(0)},
	},
	// int → NamedUint8
	{
		Name:          "int to NamedUint8 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint8{NamedUint8: 0},
	},
	{
		Name:          "int to NamedUint8 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedUint8{NamedUint8: 0x01},
	},
	{
		Name:          "int to NamedUint8 MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedUint8{NamedUint8: math.MaxInt8},
	},
	{
		Name:          "int to NamedUint8 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint8{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint8 MaxUint8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\xff"),
		Value:         SingleNamedUint8{NamedUint8: math.MaxUint8},
	},
	{
		Name:          "int to NamedUint8 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x00"),
		Value:         SingleNamedUint8{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint8 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint8{NamedUint8: 0},
	},
	{
		Name:          "int to NamedUint8 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint8{NamedUint8: 0},
	},
	// int → *NamedUint8
	{
		Name:          "int to NamedUint8Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: newNamedUint8Ptr(0)},
	},
	{
		Name:          "int to NamedUint8Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x01"),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: newNamedUint8Ptr(0x01)},
	},
	{
		Name:          "int to NamedUint8Ptr MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x7f"),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: newNamedUint8Ptr(math.MaxInt8)},
	},
	{
		Name:          "int to NamedUint8Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint8Ptr MaxUint8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\xff"),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: newNamedUint8Ptr(math.MaxUint8)},
	},
	{
		Name:          "int to NamedUint8Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x00"),
		Value:         SingleNamedUint8Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint8Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: nil},
	},
	{
		Name:          "int to NamedUint8Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint8Ptr{NamedUint8Ptr: newNamedUint8Ptr(0)},
	},
	// int → NamedUint16
	{
		Name:          "int to NamedUint16 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint16{NamedUint16: 0},
	},
	{
		Name:          "int to NamedUint16 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedUint16{NamedUint16: 0x0102},
	},
	{
		Name:          "int to NamedUint16 MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedUint16{NamedUint16: math.MaxInt16},
	},
	{
		Name:          "int to NamedUint16 overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint16{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint16 MaxUint16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\xff\xff"),
		Value:         SingleNamedUint16{NamedUint16: math.MaxUint16},
	},
	{
		Name:          "int to NamedUint16 overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x01\x00\x00"),
		Value:         SingleNamedUint16{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint16 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint16{NamedUint16: 0},
	},
	{
		Name:          "int to NamedUint16 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint16{NamedUint16: 0},
	},
	// int → *NamedUint16
	{
		Name:          "int to NamedUint16Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: newNamedUint16Ptr(0)},
	},
	{
		Name:          "int to NamedUint16Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x01\x02"),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: newNamedUint16Ptr(0x0102)},
	},
	{
		Name:          "int to NamedUint16Ptr MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x7f\xff"),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: newNamedUint16Ptr(math.MaxInt16)},
	},
	{
		Name:          "int to NamedUint16Ptr overflow negative",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint16Ptr MaxUint16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\xff\xff"),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: newNamedUint16Ptr(math.MaxUint16)},
	},
	{
		Name:          "int to NamedUint16Ptr overflow positive",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x01\x00\x00"),
		Value:         SingleInt16Ptr{},
		Error:         true,
	},
	{
		Name:          "int to NamedUint16Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: nil},
	},
	{
		Name:          "int to NamedUint16Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint16Ptr{NamedUint16Ptr: newNamedUint16Ptr(0)},
	},
	// int → NamedUint32
	{
		Name:          "int to NamedUint32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint32{NamedUint32: 0},
	},
	{
		Name:          "int to NamedUint32 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedUint32{NamedUint32: 0x01020304},
	},
	{
		Name:          "int to NamedUint32 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedUint32{NamedUint32: 0x80000000},
	},
	{
		Name:          "int to int32 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedUint32{NamedUint32: math.MaxInt32},
	},
	{
		Name:          "int to int32 MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint32{NamedUint32: math.MaxUint32},
	},
	{
		Name:          "int to NamedUint32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint32{NamedUint32: 0},
	},
	{
		Name:          "int to NamedUint32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint32{NamedUint32: 0},
	},
	// int → *NamedUint32
	{
		Name:          "int to NamedUint32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(0)},
	},
	{
		Name:          "int to NamedUint32Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedUint32Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(0x80000000)},
	},
	{
		Name:          "int to NamedUint32Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedUint32Ptr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(math.MaxUint32)},
	},
	{
		Name:          "int to NamedUint32Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: nil},
	},
	{
		Name:          "int to NamedUint32Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint32Ptr{NamedUint32Ptr: newNamedUint32Ptr(0)},
	},
	// int → NamedUint64
	{
		Name:          "int to int64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint64{NamedUint64: 0},
	},
	{
		Name:          "int to NamedUint64 byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedUint64{NamedUint64: 0x01020304},
	},
	{
		Name:          "int to NamedUint64 MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedUint64{NamedUint64: 0x80000000},
	},
	{
		Name:          "int to NamedUint64 MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedUint64{NamedUint64: math.MaxInt32},
	},
	{
		Name:          "int to NamedUint64 MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint64{NamedUint64: math.MaxUint32},
	},
	{
		Name:          "int to NamedUint64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint64{NamedUint64: 0},
	},
	{
		Name:          "int to NamedUint64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint64{NamedUint64: 0},
	},
	// int → *NamedUint64
	{
		Name:          "int to NamedUint64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(0)},
	},
	{
		Name:          "int to NamedUint64Ptr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(0x01020304)},
	},
	{
		Name:          "int to NamedUint64Ptr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(0x80000000)},
	},
	{
		Name:          "int to NamedUint64Ptr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(math.MaxInt32)},
	},
	{
		Name:          "int to NamedUint64Ptr MaxUint32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(math.MaxUint32)},
	},
	{
		Name:          "int to NamedUint64Ptr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: nil},
	},
	{
		Name:          "int to NamedUint64Ptr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleNamedUint64Ptr{NamedUint64Ptr: newNamedUint64Ptr(0)},
	},
	// int → string
	{
		Name:          "int to string zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleString{String: "0"},
	},
	{
		Name:          "int to string byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleString{String: "16909060"},
	},
	{
		Name:          "int to string MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleString{String: "-2147483648"},
	},
	{
		Name:          "int to string MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleString{String: "2147483647"},
	},
	{
		Name:          "int to string null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleString{String: "0"},
	},
	{
		Name:          "int to string empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleString{String: "0"},
	},
	// int → *string
	{
		Name:          "int to stringPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleStringPtr{StringPtr: newStringPtr("0")},
	},
	{
		Name:          "int to stringPtr byte order",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x01\x02\x03\x04"),
		Value:         SingleStringPtr{StringPtr: newStringPtr("16909060")},
	},
	{
		Name:          "int to stringPtr MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleStringPtr{StringPtr: newStringPtr("-2147483648")},
	},
	{
		Name:          "int to stringPtr MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleStringPtr{StringPtr: newStringPtr("2147483647")},
	},
	{
		Name:          "int to stringPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleStringPtr{StringPtr: nil},
	},
	{
		Name:          "int to stringPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleStringPtr{StringPtr: newStringPtr("0")},
	},
	// int → NamedString
	{
		Name:          "int to NamedString zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedString{},
		Error:         true,
	},
	// int → *NamedString
	{
		Name:          "int to NamedStringPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedStringPtr{},
		Error:         true,
	},
}

func TestUnmarshalInteger(t *testing.T) {
	for _, test := range integerTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Value, value.Elem().Interface())
			}

		})
	}
}

var bigIntTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// tinyint → big.Int
	{
		Name:          "tinyint to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "tinyint to bigInt MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x7f"),
		Value:         SingleBigInt{BigInt: *newBigInt("127")},
	},
	{
		Name:          "tinyint to bigInt MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x80"),
		Value:         SingleBigInt{BigInt: *newBigInt("-128")},
	},
	{
		Name:          "tinyint to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("-1")},
	},
	{
		Name:          "tinyint to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte(nil),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "tinyint to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte(""),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	// tinyint → *big.Int
	{
		Name:          "tinyint to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x00"),
		Value:         SingleBigIntPtr{BigIntPtr: newBigInt("0")},
	},
	{
		Name:          "tinyint to bigInt MaxInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x7f"),
		Value:         SingleBigIntPtr{BigIntPtr: newBigInt("127")},
	},
	{
		Name:          "tinyint to bigInt MinInt8",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\x80"),
		Value:         SingleBigIntPtr{BigIntPtr: newBigInt("-128")},
	},
	{
		Name:          "tinyint to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte("\xff"),
		Value:         SingleBigIntPtr{BigIntPtr: newBigInt("-1")},
	},
	{
		Name:          "tinyint to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte(nil),
		Value:         SingleBigIntPtr{BigIntPtr: nil},
	},
	{
		Name:          "tinyint to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeTinyInt, ""),
		Data:          []byte(""),
		Value:         SingleBigIntPtr{BigIntPtr: newBigInt("0")},
	},
	// smallint → big.Int
	{
		Name:          "smallint to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte("\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "smallint to bigInt MaxInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte("\x7f\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("32767")},
	},
	{
		Name:          "smallint to bigInt MinInt16",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte("\x80\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("-32768")},
	},
	{
		Name:          "smallint to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte("\xff\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("-1")},
	},
	{
		Name:          "smallint to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte(nil),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "smallint to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeSmallInt, ""),
		Data:          []byte(""),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	// int → big.Int
	{
		Name:          "int to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "int to bigInt MaxInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x7f\xff\xff\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("2147483647")},
	},
	{
		Name:          "int to bigInt MinInt32",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\x80\x00\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("-2147483648")},
	},
	{
		Name:          "int to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte("\xff\xff\xff\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("-1")},
	},
	{
		Name:          "int to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(nil),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "int to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeInt, ""),
		Data:          []byte(""),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	// bigint → big.Int
	{
		Name:          "bigint to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "bigint to bigInt MaxInt64",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte("\x7f\xff\xff\xff\xff\xff\xff\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("9223372036854775807")},
	},
	{
		Name:          "bigint to bigInt MinInt64",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte("\x80\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("-9223372036854775808")},
	},
	{
		Name:          "bigint to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte("\xff\xff\xff\xff\xff\xff\xff\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("-1")},
	},
	{
		Name:          "bigint to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte(nil),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "bigint to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBigInt, ""),
		Data:          []byte(""),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	// varint → big.Int
	{
		Name:          "varint to bigInt zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte("\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "varint to bigInt 2**71",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte("\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleBigInt{BigInt: *newBigInt("2361183241434822606848")}, // 2**71
	},
	{
		Name:          "varint to bigInt -2**63-1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte("\xFF\x7F\xFF\xFF\xFF\xFF\xFF\xFF\xFF"),
		Value:         SingleBigInt{BigInt: *newBigInt("-9223372036854775809")},
	},
	{
		Name:          "varint to bigInt -1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte("\xff"),
		Value:         SingleBigInt{BigInt: *newBigInt("-1")},
	},
	{
		Name:          "varint to bigInt null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte(nil),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
	{
		Name:          "varint to bigInt empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeVarint, ""),
		Data:          []byte(""),
		Value:         SingleBigInt{BigInt: *newBigInt("0")},
	},
}

func extractBigInt(iface interface{}) *big.Int {
	value := reflect.ValueOf(iface)
	if value.NumField() != 1 {
		panic("only single field structs supported")
	}
	fieldValue := value.Field(0).Interface()
	switch x := fieldValue.(type) {
	case *big.Int:
		return x
	case big.Int:
		return &x
	default:
		panic("unsupported type for comparison")
	}
}

func requireBigIntStructEqual(t *testing.T, expectedIface interface{}, actualIface interface{}) {
	t.Helper()
	expected := extractBigInt(expectedIface)
	actual := extractBigInt(actualIface)
	if (expected == nil) != (actual == nil) {
		t.Fatalf("%s != %s", expected, actual)
	}
	if expected != nil && expected.Cmp(actual) != 0 {
		t.Fatalf("%s != %s", expected, actual)
	}
}

func TestUnmarshalBigInt(t *testing.T) {
	for _, test := range bigIntTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				requireBigIntStructEqual(t, test.Value, value.Elem().Interface())
			}
		})
	}
}

var booleanTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// boolean → bool
	{
		Name:          "boolean to bool false",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x00"),
		Value:         SingleBool{Bool: false},
	},
	{
		Name:          "boolean to bool true",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x01"),
		Value:         SingleBool{Bool: true},
	},
	{
		Name:          "boolean to bool nonzero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x42"),
		Value:         SingleBool{Bool: true},
	},
	{
		Name:          "boolean to bool null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(nil),
		Value:         SingleBool{Bool: false},
	},
	{
		Name:          "boolean to bool empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(""),
		Value:         SingleBool{Bool: false},
	},
	// boolean → *bool
	{
		Name:          "boolean to boolPtr false",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x00"),
		Value:         SingleBoolPtr{BoolPtr: newBoolPtr(false)},
	},
	{
		Name:          "boolean to boolPtr true",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x01"),
		Value:         SingleBoolPtr{BoolPtr: newBoolPtr(true)},
	},
	{
		Name:          "boolean to boolPtr nonzero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x42"),
		Value:         SingleBoolPtr{BoolPtr: newBoolPtr(true)},
	},
	{
		Name:          "boolean to boolPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(nil),
		Value:         SingleBoolPtr{BoolPtr: nil},
	},
	{
		Name:          "boolean to boolPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(""),
		Value:         SingleBoolPtr{BoolPtr: newBoolPtr(false)},
	},
	// boolean → NamedBool
	{
		Name:          "boolean to NamedBool false",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x00"),
		Value:         SingleNamedBool{NamedBool: false},
	},
	{
		Name:          "boolean to NamedBool true",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x01"),
		Value:         SingleNamedBool{NamedBool: true},
	},
	{
		Name:          "boolean to NamedBool nonzero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x42"),
		Value:         SingleNamedBool{NamedBool: true},
	},
	{
		Name:          "boolean to NamedBool null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(nil),
		Value:         SingleNamedBool{NamedBool: false},
	},
	{
		Name:          "boolean to NamedBool empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(""),
		Value:         SingleNamedBool{NamedBool: false},
	},
	// boolean → *NamedBool
	{
		Name:          "boolean to NamedBoolPtr false",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x00"),
		Value:         SingleNamedBoolPtr{NamedBoolPtr: newNamedBoolPtr(false)},
	},
	{
		Name:          "boolean to NamedBoolPtr true",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x01"),
		Value:         SingleNamedBoolPtr{NamedBoolPtr: newNamedBoolPtr(true)},
	},
	{
		Name:          "boolean to NamedBoolPtr nonzero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte("\x42"),
		Value:         SingleNamedBoolPtr{NamedBoolPtr: newNamedBoolPtr(true)},
	},
	{
		Name:          "boolean to NamedBoolPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(nil),
		Value:         SingleNamedBoolPtr{NamedBoolPtr: nil},
	},
	{
		Name:          "boolean to NamedBoolPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeBoolean, ""),
		Data:          []byte(""),
		Value:         SingleNamedBoolPtr{NamedBoolPtr: newNamedBoolPtr(false)},
	},
}

func TestUnmarshalBoolean(t *testing.T) {
	for _, test := range booleanTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Value, value.Elem().Interface())
			}

		})
	}
}

var decimalTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// decimal → inf.Dec
	{
		Name:          "decimal to Dec zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x00\x00"),
		Value:         SingleDec{Dec: *newDec("0")},
	},
	{
		Name:          "decimal to Dec 100",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x00\x64"),
		Value:         SingleDec{Dec: *newDec("100")},
	},
	{
		Name:          "decimal to Dec 0.25",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x02\x19"),
		Value:         SingleDec{Dec: *newDec("0.25")},
	},
	{
		Name:          "decimal to Dec cqlrb1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x13\xD5\a;\x20\x14\xA2\x91"),
		Value:         SingleDec{Dec: *newDec("-0.0012095473475870063")}, // From the iconara/cql-rb test suite
	},
	{
		Name:          "decimal to Dec cqlrb2",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x13*\xF8\xC4\xDF\xEB]o"),
		Value:         SingleDec{Dec: *newDec("0.0012095473475870063")}, // From the iconara/cql-rb test suite
	},
	{
		Name:          "decimal to Dec cqlrb3",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x12\xF2\xD8\x02\xB6R\x7F\x99\xEE\x98#\x99\xA9V"),
		Value:         SingleDec{Dec: *newDec("-1042342234234.123423435647768234")}, // From the iconara/cql-rb
	},
	{
		Name:          "decimal to Dec python1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\r\nJ\x04\"^\x91\x04\x8a\xb1\x18\xfe"),
		Value:         SingleDec{Dec: *newDec("1243878957943.1234124191998")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to Dec python2",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x06\xe5\xde]\x98Y"),
		Value:         SingleDec{Dec: *newDec("-112233.441191")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to Dec python3",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x14\x00\xfa\xce"),
		Value:         SingleDec{Dec: *newDec("0.00000000000000064206")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to Dec python4",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x14\xff\x052"),
		Value:         SingleDec{Dec: *newDec("-0.00000000000000064206")},
	},
	{
		Name:          "decimal to Dec python5",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\xff\xff\xff\x9c\x00\xfa\xce"),
		Value:         SingleDec{Dec: *inf.NewDec(64206, -100)},
	},
	{
		Name:          "decimal to Dec null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte(nil),
		Value:         SingleDec{},
		Error:         true,
	},
	{
		Name:          "decimal to Dec empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte(""),
		Value:         SingleDec{},
		Error:         true,
	},
	// decimal → *inf.Dec
	{
		Name:          "decimal to DecPtr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x00\x00"),
		Value:         SingleDecPtr{DecPtr: newDec("0")},
	},
	{
		Name:          "decimal to DecPtr 100",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x00\x64"),
		Value:         SingleDecPtr{DecPtr: newDec("100")},
	},
	{
		Name:          "decimal to DecPtr 0.25",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x02\x19"),
		Value:         SingleDecPtr{DecPtr: newDec("0.25")},
	},
	{
		Name:          "decimal to DecPtr cqlrb1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x13\xD5\a;\x20\x14\xA2\x91"),
		Value:         SingleDecPtr{DecPtr: newDec("-0.0012095473475870063")}, // From the iconara/cql-rb test suite
	},
	{
		Name:          "decimal to DecPtr cqlrb2",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x13*\xF8\xC4\xDF\xEB]o"),
		Value:         SingleDecPtr{DecPtr: newDec("0.0012095473475870063")}, // From the iconara/cql-rb test suite
	},
	{
		Name:          "decimal to DecPtr cqlrb3",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x12\xF2\xD8\x02\xB6R\x7F\x99\xEE\x98#\x99\xA9V"),
		Value:         SingleDecPtr{DecPtr: newDec("-1042342234234.123423435647768234")}, // From the iconara/cql-rb
	},
	{
		Name:          "decimal to DecPtr python1",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\r\nJ\x04\"^\x91\x04\x8a\xb1\x18\xfe"),
		Value:         SingleDecPtr{DecPtr: newDec("1243878957943.1234124191998")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to DecPtr python2",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x06\xe5\xde]\x98Y"),
		Value:         SingleDecPtr{DecPtr: newDec("-112233.441191")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to DecPtr python3",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x14\x00\xfa\xce"),
		Value:         SingleDecPtr{DecPtr: newDec("0.00000000000000064206")}, // From datastax/python-driver tests
	},
	{
		Name:          "decimal to DecPtr python4",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\x00\x00\x00\x14\xff\x052"),
		Value:         SingleDecPtr{DecPtr: newDec("-0.00000000000000064206")},
	},
	{
		Name:          "decimal to DecPtr python5",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte("\xff\xff\xff\x9c\x00\xfa\xce"),
		Value:         SingleDecPtr{DecPtr: inf.NewDec(64206, -100)},
	},
	{
		Name:          "decimal to DecPtr null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte(nil),
		Value:         SingleDecPtr{DecPtr: nil},
	},
	{
		Name:          "decimal to DecPtr empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDecimal, ""),
		Data:          []byte(""),
		Value:         SingleDecPtr{},
		Error:         true,
	},
}

func extractDec(iface interface{}) *inf.Dec {
	value := reflect.ValueOf(iface)
	if value.NumField() != 1 {
		panic("only single field structs supported")
	}
	fieldValue := value.Field(0).Interface()
	switch x := fieldValue.(type) {
	case *inf.Dec:
		return x
	case inf.Dec:
		return &x
	default:
		panic(fmt.Sprintf("unsupported type for comparison: %T", iface))
	}
}

func requireDecStructEqual(t *testing.T, expectedIface interface{}, actualIface interface{}) {
	t.Helper()
	expected := extractDec(expectedIface)
	actual := extractDec(actualIface)
	if (expected == nil) != (actual == nil) {
		t.Fatalf("%s != %s", expected, actual)
	}
	if expected != nil && expected.Cmp(actual) != 0 {
		t.Fatalf("%s != %s", expected, actual)
	}
}

func TestUnmarshalDec(t *testing.T) {
	for _, test := range decimalTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				requireDecStructEqual(t, test.Value, value.Elem().Interface())
			}
		})
	}
}

var floatTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// float → float32
	{
		Name:          "float to float32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleFloat32{Float32: 0},
	},
	{
		Name:          "float to float32 value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x40\x49\x0f\xdb"),
		Value:         SingleFloat32{Float32: 3.14159265},
	},
	{
		Name:          "float to float32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(nil),
		Value:         SingleFloat32{Float32: 0},
	},
	{
		Name:          "float to float32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(""),
		Value:         SingleFloat32{Float32: 0},
	},
	// float → *float32
	{
		Name:          "float to float32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleFloat32Ptr{Float32Ptr: newFloat32Ptr(0)},
	},
	{
		Name:          "float to float32Ptr value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x40\x49\x0f\xdb"),
		Value:         SingleFloat32Ptr{Float32Ptr: newFloat32Ptr(3.14159265)},
	},
	{
		Name:          "float to float32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(nil),
		Value:         SingleFloat32Ptr{Float32Ptr: nil},
	},
	{
		Name:          "float to float32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(""),
		Value:         SingleFloat32Ptr{Float32Ptr: newFloat32Ptr(0)},
	},
	// float → NamedFloat32
	{
		Name:          "float to NamedFloat32 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedFloat32{NamedFloat32: 0},
	},
	{
		Name:          "float to NamedFloat32 value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x40\x49\x0f\xdb"),
		Value:         SingleNamedFloat32{NamedFloat32: 3.14159265},
	},
	{
		Name:          "float to NamedFloat32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(nil),
		Value:         SingleNamedFloat32{NamedFloat32: 0},
	},
	{
		Name:          "float to NamedFloat32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(""),
		Value:         SingleNamedFloat32{NamedFloat32: 0},
	},
	// float → *NamedFloat32
	{
		Name:          "float to NamedFloat32Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x00\x00\x00\x00"),
		Value:         SingleNamedFloat32Ptr{NamedFloat32Ptr: newNamedFloat32Ptr(0)},
	},
	{
		Name:          "float to NamedFloat32Ptr value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte("\x40\x49\x0f\xdb"),
		Value:         SingleNamedFloat32Ptr{NamedFloat32Ptr: newNamedFloat32Ptr(3.14159265)},
	},
	{
		Name:          "float to NamedFloat32 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(nil),
		Value:         SingleNamedFloat32Ptr{NamedFloat32Ptr: nil},
	},
	{
		Name:          "float to NamedFloat32 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeFloat, ""),
		Data:          []byte(""),
		Value:         SingleNamedFloat32Ptr{NamedFloat32Ptr: newNamedFloat32Ptr(0)},
	},
}

func TestUnmarshalFloat(t *testing.T) {
	for _, test := range floatTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Value, value.Elem().Interface())
			}

		})
	}
}

var doubleTests = []struct {
	Name          string
	FieldTypeInfo gocql.TypeInfo
	Data          []byte
	Value         interface{}
	Error         bool
}{
	// double → float64
	{
		Name:          "double to float64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleFloat64{Float64: 0},
	},
	{
		Name:          "double to float64 value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x40\x09\x21\xfb\x53\xc8\xd4\xf1"),
		Value:         SingleFloat64{Float64: 3.14159265},
	},
	{
		Name:          "double to float64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(nil),
		Value:         SingleFloat64{Float64: 0},
	},
	{
		Name:          "double to float64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(""),
		Value:         SingleFloat64{Float64: 0},
	},
	// double → *float64
	{
		Name:          "double to float64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleFloat64Ptr{Float64Ptr: newFloat64Ptr(0)},
	},
	{
		Name:          "double to float64Ptr value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x40\x09\x21\xfb\x53\xc8\xd4\xf1"),
		Value:         SingleFloat64Ptr{Float64Ptr: newFloat64Ptr(3.14159265)},
	},
	{
		Name:          "double to float64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(nil),
		Value:         SingleFloat64Ptr{Float64Ptr: nil},
	},
	{
		Name:          "double to float64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(""),
		Value:         SingleFloat64Ptr{Float64Ptr: newFloat64Ptr(0)},
	},
	// double → NamedFloat64
	{
		Name:          "double to NamedFloat64 zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleNamedFloat64{NamedFloat64: 0},
	},
	{
		Name:          "double to NamedFloat64 value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x40\x09\x21\xfb\x53\xc8\xd4\xf1"),
		Value:         SingleNamedFloat64{NamedFloat64: 3.14159265},
	},
	{
		Name:          "double to NamedFloat64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(nil),
		Value:         SingleNamedFloat64{NamedFloat64: 0},
	},
	{
		Name:          "double to NamedFloat64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(""),
		Value:         SingleNamedFloat64{NamedFloat64: 0},
	},
	// double → *NamedFloat64
	{
		Name:          "double to NamedFloat64Ptr zero",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		Value:         SingleNamedFloat64Ptr{NamedFloat64Ptr: newNamedFloat64Ptr(0)},
	},
	{
		Name:          "double to NamedFloat64Ptr value",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte("\x40\x09\x21\xfb\x53\xc8\xd4\xf1"),
		Value:         SingleNamedFloat64Ptr{NamedFloat64Ptr: newNamedFloat64Ptr(3.14159265)},
	},
	{
		Name:          "double to NamedFloat64 null",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(nil),
		Value:         SingleNamedFloat64Ptr{NamedFloat64Ptr: nil},
	},
	{
		Name:          "double to NamedFloat64 empty",
		FieldTypeInfo: gocql.NewNativeType(2, gocql.TypeDouble, ""),
		Data:          []byte(""),
		Value:         SingleNamedFloat64Ptr{NamedFloat64Ptr: newNamedFloat64Ptr(0)},
	},
}

func TestUnmarshalDouble(t *testing.T) {
	for _, test := range doubleTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			typeInfo, data := buildUDT(reflect.ValueOf(test.Value).Type(), test.FieldTypeInfo, test.Data)
			value := reflect.New(reflect.TypeOf(test.Value))
			err := gocql.Unmarshal(typeInfo, data, value.Interface())
			if test.Error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Value, value.Elem().Interface())
			}

		})
	}
}

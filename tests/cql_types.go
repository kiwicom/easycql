package tests

import (
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
)

type CQLVarcharTypes struct {
	String    string
	StringPtr *string

	NamedString    NamedString
	NamedStringPtr *NamedString

	Bytes    []byte
	BytesPtr *[]byte

	NamedBytes    NamedBytes
	NamedBytesPtr *NamedBytes

	CustomString    CustomString
	CustomStringPtr *CustomString
}

type SingleInt struct {
	Int int
}

type SingleInt8 struct {
	Int8 int8
}

type SingleInt16 struct {
	Int16 int16
}

type SingleInt32 struct {
	Int32 int32
}

type SingleInt64 struct {
	Int64 int64
}

type SingleIntPtr struct {
	IntPtr *int
}

type SingleInt8Ptr struct {
	Int8Ptr *int8
}

type SingleInt16Ptr struct {
	Int16Ptr *int16
}

type SingleInt32Ptr struct {
	Int32Ptr *int32
}

type SingleInt64Ptr struct {
	Int64Ptr *int64
}

type SingleNamedInt struct {
	NamedInt NamedInt
}

type SingleNamedInt8 struct {
	NamedInt8 NamedInt8
}

type SingleNamedInt16 struct {
	NamedInt16 NamedInt16
}

type SingleNamedInt32 struct {
	NamedInt32 NamedInt32
}

type SingleNamedInt64 struct {
	NamedInt64 NamedInt64
}

type SingleNamedIntPtr struct {
	NamedIntPtr *NamedInt
}

type SingleNamedInt8Ptr struct {
	NamedInt8Ptr *NamedInt8
}

type SingleNamedInt16Ptr struct {
	NamedInt16Ptr *NamedInt16
}

type SingleNamedInt32Ptr struct {
	NamedInt32Ptr *NamedInt32
}

type SingleNamedInt64Ptr struct {
	NamedInt64Ptr *NamedInt64
}

type SingleUint struct {
	Uint uint
}

type SingleUint8 struct {
	Uint8 uint8
}

type SingleUint16 struct {
	Uint16 uint16
}

type SingleUint32 struct {
	Uint32 uint32
}

type SingleUint64 struct {
	Uint64 uint64
}

type SingleUintPtr struct {
	UintPtr *uint
}

type SingleUint8Ptr struct {
	Uint8Ptr *uint8
}

type SingleUint16Ptr struct {
	Uint16Ptr *uint16
}

type SingleUint32Ptr struct {
	Uint32Ptr *uint32
}

type SingleUint64Ptr struct {
	Uint64Ptr *uint64
}

type SingleNamedUint struct {
	NamedUint NamedUint
}

type SingleNamedUint8 struct {
	NamedUint8 NamedUint8
}

type SingleNamedUint16 struct {
	NamedUint16 NamedUint16
}

type SingleNamedUint32 struct {
	NamedUint32 NamedUint32
}

type SingleNamedUint64 struct {
	NamedUint64 NamedUint64
}

type SingleNamedUintPtr struct {
	NamedUintPtr *NamedUint
}

type SingleNamedUint8Ptr struct {
	NamedUint8Ptr *NamedUint8
}

type SingleNamedUint16Ptr struct {
	NamedUint16Ptr *NamedUint16
}

type SingleNamedUint32Ptr struct {
	NamedUint32Ptr *NamedUint32
}

type SingleNamedUint64Ptr struct {
	NamedUint64Ptr *NamedUint64
}

type SingleString struct {
	String string
}

type SingleStringPtr struct {
	StringPtr *string
}

type SingleNamedString struct {
	NamedString NamedString
}

type SingleNamedStringPtr struct {
	NamedStringPtr *NamedString
}

type SingleBigInt struct {
	BigInt big.Int
}

type SingleBigIntPtr struct {
	BigIntPtr *big.Int
}

type SingleBool struct {
	Bool bool
}

type SingleBoolPtr struct {
	BoolPtr *bool
}

type SingleNamedBool struct {
	NamedBool NamedBool
}

type SingleNamedBoolPtr struct {
	NamedBoolPtr *NamedBool
}

type SingleFloat32 struct {
	Float32 float32
}

type SingleFloat32Ptr struct {
	Float32Ptr *float32
}

type SingleNamedFloat32 struct {
	NamedFloat32 NamedFloat32
}

type SingleNamedFloat32Ptr struct {
	NamedFloat32Ptr *NamedFloat32
}

type SingleFloat64 struct {
	Float64 float64
}

type SingleFloat64Ptr struct {
	Float64Ptr *float64
}

type SingleNamedFloat64 struct {
	NamedFloat64 NamedFloat64
}

type SingleNamedFloat64Ptr struct {
	NamedFloat64Ptr *NamedFloat64
}

type SingleDec struct {
	Dec inf.Dec
}

type SingleDecPtr struct {
	DecPtr *inf.Dec
}

type CQLTimestampTypes struct {
	Int64    int64
	Int64Ptr *int64

	NamedInt64    NamedInt64
	NamedInt64Ptr *NamedInt64

	Time    time.Time
	TimePtr *time.Time
}

type CQLUUIDTypes struct {
	String    string
	StringPtr *string

	Bytes    []byte
	BytesPtr *[]byte

	UUID    gocql.UUID
	UUIDPtr *gocql.UUID
}

type CQLTimeUUIDTypes struct {
	String    string
	StringPtr *string

	Bytes    []byte
	BytesPtr *[]byte

	UUID    gocql.UUID
	UUIDPtr *gocql.UUID

	Time    time.Time
	TimePtr *time.Time
}

type CQLInetTypes struct {
	String    string
	StringPtr *string

	IP    net.IP
	IPPtr *net.IP
}

type CQLDateTypes struct {
	String    string
	StringPtr *string

	Time    time.Time
	TimePtr *time.Time
}

type CQLTimeTypes struct {
	Int64    int64
	Int64Ptr *int64

	NamedInt64    NamedInt64
	NamedInt64Ptr *NamedInt64

	Duration    time.Duration
	DurationPtr *time.Duration
}

type CQLDurationTypes struct {
	Duration    gocql.Duration
	DurationPtr *gocql.Duration
}

type CQLListTypes struct {
	StringSlice    []string
	StringSlicePtr *[]string

	StringArray    [5]string
	StringArrayPtr *[5]string

	NamedStringSlice    NamedStringSlice
	NamedStringSlicePtr *NamedStringSlice

	NamedStringArray    NamedStringArray
	NamedStringArrayPtr *NamedStringArray
}

type (
	NamedBytes       []byte
	NamedStringSlice []string
	NamedStringArray [5]string
	CustomString     string
)

func (c CustomString) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(strings.ToLower(string(c))), nil
}

func (c *CustomString) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	*c = CustomString(strings.ToUpper(string(data)))
	return nil
}

func newBigInt(s string) *big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return i
}

func newDec(s string) *inf.Dec {
	i, _ := new(inf.Dec).SetString(s)
	return i
}

package tests

import (
	"math"
	"net"
	"time"

	"github.com/mailru/easyjson/opt"
)

type PrimitiveTypes struct {
	String string
	Bool   bool

	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64

	Float32 float32
	Float64 float64

	Ptr    *string
	PtrNil *string
}

var str = "bla"

var primitiveTypesValue = PrimitiveTypes{
	String: "test", Bool: true,

	Int:   math.MinInt32,
	Int8:  math.MinInt8,
	Int16: math.MinInt16,
	Int32: math.MinInt32,
	Int64: math.MinInt64,

	Uint:   math.MaxUint32,
	Uint8:  math.MaxUint8,
	Uint16: math.MaxUint16,
	Uint32: math.MaxUint32,
	Uint64: math.MaxUint64,

	Float32: 1.5,
	Float64: math.MaxFloat64,

	Ptr: &str,
}

type (
	NamedString string
	NamedBool   bool

	NamedInt   int
	NamedInt8  int8
	NamedInt16 int16
	NamedInt32 int32
	NamedInt64 int64

	NamedUint   uint
	NamedUint8  uint8
	NamedUint16 uint16
	NamedUint32 uint32
	NamedUint64 uint64

	NamedFloat32 float32
	NamedFloat64 float64

	NamedStrPtr *string
)

type NamedPrimitiveTypes struct {
	String NamedString
	Bool   NamedBool

	Int   NamedInt
	Int8  NamedInt8
	Int16 NamedInt16
	Int32 NamedInt32
	Int64 NamedInt64

	Uint   NamedUint
	Uint8  NamedUint8
	Uint16 NamedUint16
	Uint32 NamedUint32
	Uint64 NamedUint64

	Float32 NamedFloat32
	Float64 NamedFloat64

	Ptr    NamedStrPtr
	PtrNil NamedStrPtr
}

var namedPrimitiveTypesValue = NamedPrimitiveTypes{
	String: "test",
	Bool:   true,

	Int:   math.MinInt32,
	Int8:  math.MinInt8,
	Int16: math.MinInt16,
	Int32: math.MinInt32,
	Int64: math.MinInt64,

	Uint:   math.MaxUint32,
	Uint8:  math.MaxUint8,
	Uint16: math.MaxUint16,
	Uint32: math.MaxUint32,
	Uint64: math.MaxUint64,

	Float32: 1.5,
	Float64: math.MaxFloat64,

	Ptr: NamedStrPtr(&str),
}

type SubStruct struct {
	Value     string
	Value2    string
	unexpored bool //nolint:structcheck // unused field for testing
}

type SubP struct {
	V string
}

type SubStructAlias SubStruct

type Structs struct {
	SubStruct
	*SubP

	Value2 int

	Sub1   SubStruct `json:"substruct"`
	Sub2   *SubStruct
	SubNil *SubStruct

	SubSlice    []SubStruct
	SubSliceNil []SubStruct

	SubPtrSlice    []*SubStruct
	SubPtrSliceNil []*SubStruct

	SubA1 SubStructAlias
	SubA2 *SubStructAlias

	Anonymous struct {
		V string
		I int
	}
	Anonymous1 *struct {
		V string
	}

	AnonymousSlice    []struct{ V int }
	AnonymousPtrSlice []*struct{ V int }

	Slice []string

	unexported bool //nolint:structcheck // unused field for testing
}

var structsValue = Structs{
	SubStruct: SubStruct{Value: "test"},
	SubP:      &SubP{V: "subp"},

	Value2: 5,

	Sub1: SubStruct{Value: "test1", Value2: "v"},
	Sub2: &SubStruct{Value: "test2", Value2: "v2"},

	SubSlice: []SubStruct{
		{Value: "s1"},
		{Value: "s2"},
	},

	SubPtrSlice: []*SubStruct{
		{Value: "p1"},
		{Value: "p2"},
	},

	SubA1: SubStructAlias{Value: "test3", Value2: "v3"},
	SubA2: &SubStructAlias{Value: "test4", Value2: "v4"},

	Anonymous: struct {
		V string
		I int
	}{V: "bla", I: 5},

	Anonymous1: &struct {
		V string
	}{V: "bla1"},

	AnonymousSlice:    []struct{ V int }{{1}, {2}},
	AnonymousPtrSlice: []*struct{ V int }{{3}, {4}},

	Slice: []string{"test5", "test6"},
}

type Opts struct {
	StrNull  opt.String
	StrEmpty opt.String
	Str      opt.String

	IntNull opt.Int
	IntZero opt.Int
	Int     opt.Int
}

var optsValue = Opts{
	StrEmpty: opt.OString(""),
	Str:      opt.OString("test"),

	IntZero: opt.OInt(0),
	Int:     opt.OInt(5),
}

type StdMarshaler struct {
	T  time.Time
	IP net.IP
}

var stdMarshalerValue = StdMarshaler{
	//nolint:gocritic // octalLiteral didn't like the leading zeroes but here they actually help with readability
	T:  time.Date(2016, 01, 02, 14, 15, 10, 0, time.UTC),
	IP: net.IPv4(192, 168, 0, 1),
}

type unexportedStruct struct {
	Value string
}

var unexportedStructValue = unexportedStruct{"test"}
var unexportedStructString = `{"Value":"test"}`

type ExcludedField struct {
	Process       bool `json:"process"`
	DoNotProcess  bool `json:"-"`
	DoNotProcess1 bool `json:"-"`
}

var excludedFieldValue = ExcludedField{
	Process:       true,
	DoNotProcess:  false,
	DoNotProcess1: false,
}
var excludedFieldString = `{"process":true}`

type Slices struct {
	ByteSlice      []byte
	EmptyByteSlice []byte
	NilByteSlice   []byte
	IntSlice       []int
	EmptyIntSlice  []int
	NilIntSlice    []int
}

var sliceValue = Slices{
	ByteSlice:      []byte("abc"),
	EmptyByteSlice: []byte{},
	NilByteSlice:   []byte(nil),
	IntSlice:       []int{1, 2, 3, 4, 5},
	EmptyIntSlice:  []int{},
	NilIntSlice:    []int(nil),
}

type Arrays struct {
	ByteArray      [3]byte
	EmptyByteArray [0]byte
	IntArray       [5]int
	EmptyIntArray  [0]int
}

var arrayValue = Arrays{
	ByteArray:      [3]byte{'a', 'b', 'c'},
	EmptyByteArray: [0]byte{},
	IntArray:       [5]int{1, 2, 3, 4, 5},
	EmptyIntArray:  [0]int{},
}

var arrayUnderflowValue = Arrays{
	ByteArray:      [3]byte{'x', 0, 0},
	EmptyByteArray: [0]byte{},
	IntArray:       [5]int{1, 2, 0, 0, 0},
	EmptyIntArray:  [0]int{},
}

type Str string

type Maps struct {
	Map          map[string]string
	InterfaceMap map[string]interface{}
	NilMap       map[string]string

	CustomMap map[Str]Str
}

var mapsValue = Maps{
	Map:          map[string]string{"A": "b"}, // only one item since map iteration is randomized
	InterfaceMap: map[string]interface{}{"G": float64(1)},

	CustomMap: map[Str]Str{"c": "d"},
}

type NamedSlice []Str
type NamedMap map[Str]Str

type DeepNest struct {
	SliceMap         map[Str][]Str
	SliceMap1        map[Str][]Str
	SliceMap2        map[Str][]Str
	NamedSliceMap    map[Str]NamedSlice
	NamedMapMap      map[Str]NamedMap
	MapSlice         []map[Str]Str
	NamedSliceSlice  []NamedSlice
	NamedMapSlice    []NamedMap
	NamedStringSlice []NamedString
}

var deepNestValue = DeepNest{
	SliceMap: map[Str][]Str{
		"testSliceMap": {
			"0",
			"1",
		},
	},
	SliceMap1: map[Str][]Str{
		"testSliceMap1": []Str(nil),
	},
	SliceMap2: map[Str][]Str{
		"testSliceMap2": {},
	},
	NamedSliceMap: map[Str]NamedSlice{
		"testNamedSliceMap": {
			"2",
			"3",
		},
	},
	NamedMapMap: map[Str]NamedMap{
		"testNamedMapMap": {
			"key1": "value1",
		},
	},
	MapSlice: []map[Str]Str{
		{
			"testMapSlice": "someValue",
		},
	},
	NamedSliceSlice: []NamedSlice{
		{
			"someValue1",
			"someValue2",
		},
		{
			"someValue3",
			"someValue4",
		},
	},
	NamedMapSlice: []NamedMap{
		{
			"key2": "value2",
		},
		{
			"key3": "value3",
		},
	},
	NamedStringSlice: []NamedString{
		"value4", "value5",
	},
}

//easycql:cql
type Ints []int

var IntsValue = Ints{1, 2, 3, 4, 5}

var IntsString = `[1,2,3,4,5]`

//easycql:cql
type MapStringString map[string]string

var mapStringStringValue = MapStringString{"a": "b"}

//nolint:staticcheck
type RequiredOptionalStruct struct {
	FirstName string `json:"first_name,required"`
	Lastname  string `json:"last_name"`
}

//easycql:cql
type EncodingFlagsTestMap struct {
	F map[string]string
}

//easycql:cql
type EncodingFlagsTestSlice struct {
	F []string
}

type StructWithInterface struct {
	Field1 int         `json:"f1"`
	Field2 interface{} `json:"f2"`
	Field3 string      `json:"f3"`
}

type EmbeddedStruct struct {
	Field1 int    `json:"f1"`
	Field2 string `json:"f2"`
}

var structWithInterfaceValueFilled = StructWithInterface{1, &EmbeddedStruct{11, "22"}, "3"}

//easycql:cql
type MapIntString map[int]string

var mapIntStringValue = MapIntString{3: "hi"}

//easycql:cql
type MapInt32String map[int32]string

var mapInt32StringValue = MapInt32String{-354634382: "life"}

//easycql:cql
type MapInt64String map[int64]string

var mapInt64StringValue = MapInt64String{-3546343826724305832: "life"}

//easycql:cql
type MapUintString map[uint]string

var mapUintStringValue = MapUintString{42: "life"}

//easycql:cql
type MapUint32String map[uint32]string

var mapUint32StringValue = MapUint32String{354634382: "life"}

//easycql:cql
type MapUint64String map[uint64]string

var mapUint64StringValue = MapUint64String{3546343826724305832: "life"}

//easycql:cql
type MapUintptrString map[uintptr]string

var mapUintptrStringValue = MapUintptrString{272679208: "obj"}

type MyInt int

//easycql:cql
type MapMyIntString map[MyInt]string

var mapMyIntStringValue = MapMyIntString{MyInt(42): "life"}

//easycql:cql
type IntKeyedMapStruct struct {
	Foo MapMyIntString            `json:"foo"`
	Bar map[int16]MapUint32String `json:"bar"`
}

var intKeyedMapStructValue = IntKeyedMapStruct{
	Foo: mapMyIntStringValue,
	Bar: map[int16]MapUint32String{32: mapUint32StringValue},
}

type IntArray [2]int

//easycql:cql
type IntArrayStruct struct {
	Pointer *IntArray `json:"pointer"`
	Value   IntArray  `json:"value"`
}

var intArrayStructValue = IntArrayStruct{
	Pointer: &IntArray{1, 2},
	Value:   IntArray{1, 2},
}

type MyUInt8 uint8

//easycql:cql
type MyUInt8Slice []MyUInt8

var myUInt8SliceValue = MyUInt8Slice{1, 2, 3, 4, 5}

//easycql:cql
type MyUInt8Array [2]MyUInt8

var myUInt8ArrayValue = MyUInt8Array{1, 2}

type PreferredType struct {
	Value int64 `easycql:"value,smallint"`
}

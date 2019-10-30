package gen

import (
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
)

func (g *Generator) getDecoderName(t reflect.Type) string {
	return g.functionName("decode", t)
}

// decoderGen is a function that generates unmarshaller code.
// It unmarshals t from bytes stored in variable which name is stored in `in` and stores the result into `out`.
// The gocql.TypeInfo is stored in variable named info.
// tags describe the field tags for the field being unmarshaled and indent specifies how much to indent the output.
type decoderGen func(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error

var stringType = reflect.TypeOf((*string)(nil)).Elem()
var byteSliceType = reflect.TypeOf((*[]byte)(nil)).Elem()
var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
var infDecType = reflect.TypeOf((*inf.Dec)(nil)).Elem()

var decodersByKind = map[reflect.Kind]decoderMeta{
	reflect.String: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeVarchar: varcharToStringDecoder,
			gocql.TypeAscii:   varcharToStringDecoder,
			gocql.TypeBlob:    varcharToStringDecoder,
			gocql.TypeText:    varcharToStringDecoder,
		},
		preferredType: gocql.TypeVarchar,
		complete:      true,
	},
	reflect.Int: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeInt,
	},
	reflect.Int8: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeTinyInt,
	},
	reflect.Int16: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeSmallInt,
	},
	reflect.Int32: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeInt,
	},
	reflect.Int64: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeBigInt,
	},
	reflect.Uint: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeInt,
	},
	reflect.Uint8: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeTinyInt,
	},
	reflect.Uint16: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeSmallInt,
	},
	reflect.Uint32: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeInt,
	},
	reflect.Uint64: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  intLikeToIntDecoder(gocql.TypeTinyInt),
			gocql.TypeSmallInt: intLikeToIntDecoder(gocql.TypeSmallInt),
			gocql.TypeInt:      intLikeToIntDecoder(gocql.TypeInt),
			gocql.TypeBigInt:   intLikeToIntDecoder(gocql.TypeBigInt),
		},
		preferredType: gocql.TypeBigInt,
	},
	reflect.Bool: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeBoolean: booleanToBoolDecoder,
		},
		preferredType: gocql.TypeBoolean,
		complete:      true,
	},
	reflect.Float32: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeFloat: floatToFloat32Decoder,
		},
		preferredType: gocql.TypeFloat,
		complete:      true,
	},
	reflect.Float64: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeDouble: doubleToFloat64Decoder,
		},
		preferredType: gocql.TypeDouble,
		complete:      true,
	},
}

type decoderMeta struct {
	// map of implemented generators per gocql type
	cqlTypes map[gocql.Type]decoderGen
	// default preferred gocql type.
	// preferred type will be put first in the switch statement.
	preferredType gocql.Type
	// complete indicates whether we have decoder for all CQL types supported by gocql for this Go type.
	complete bool
}

var decodersByType = map[reflect.Type]decoderMeta{
	stringType: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeVarchar:  varcharToStringDecoder,
			gocql.TypeAscii:    varcharToStringDecoder,
			gocql.TypeBlob:     varcharToStringDecoder,
			gocql.TypeText:     varcharToStringDecoder,
			gocql.TypeTinyInt:  intLikeToStringDecoder("DecTiny"),
			gocql.TypeSmallInt: intLikeToStringDecoder("DecShort"),
			gocql.TypeInt:      intLikeToStringDecoder("DecInt"),
			gocql.TypeBigInt:   intLikeToStringDecoder("DecBigInt"),
			gocql.TypeVarint:   varIntToStringDecoder,
			gocql.TypeUUID:     uuidToStringDecoder,
			gocql.TypeTimeUUID: uuidToStringDecoder,
			gocql.TypeInet:     inetToStringDecoder,
		},
		preferredType: gocql.TypeVarchar,
		complete:      true,
	},
	byteSliceType: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeVarchar:  varcharToBytesDecoder,
			gocql.TypeAscii:    varcharToBytesDecoder,
			gocql.TypeBlob:     varcharToBytesDecoder,
			gocql.TypeText:     varcharToBytesDecoder,
			gocql.TypeUUID:     uuidToBytesDecoder,
			gocql.TypeTimeUUID: uuidToBytesDecoder,
		},
		preferredType: gocql.TypeVarchar,
	},
	bigIntType: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeTinyInt:  bigIntDecoder,
			gocql.TypeSmallInt: bigIntDecoder,
			gocql.TypeInt:      bigIntDecoder,
			gocql.TypeBigInt:   bigIntDecoder,
			gocql.TypeCounter:  bigIntDecoder,
			gocql.TypeVarint:   bigIntDecoder,
		},
		preferredType: gocql.TypeVarint,
		complete:      true,
	},
	infDecType: {
		cqlTypes: map[gocql.Type]decoderGen{
			gocql.TypeDecimal: decimalDecoder,
		},
		preferredType: gocql.TypeDecimal,
		complete:      true,
	},
}

var gocqlTypes = map[gocql.Type]string{
	gocql.TypeCustom:    "TypeCustom",
	gocql.TypeAscii:     "TypeAscii",
	gocql.TypeBigInt:    "TypeBigInt",
	gocql.TypeBlob:      "TypeBlob",
	gocql.TypeBoolean:   "TypeBoolean",
	gocql.TypeCounter:   "TypeCounter",
	gocql.TypeDecimal:   "TypeDecimal",
	gocql.TypeDouble:    "TypeDouble",
	gocql.TypeFloat:     "TypeFloat",
	gocql.TypeInt:       "TypeInt",
	gocql.TypeText:      "TypeText",
	gocql.TypeTimestamp: "TypeTimestamp",
	gocql.TypeUUID:      "TypeUUID",
	gocql.TypeVarchar:   "TypeVarchar",
	gocql.TypeVarint:    "TypeVarint",
	gocql.TypeTimeUUID:  "TypeTimeUUID",
	gocql.TypeInet:      "TypeInet",
	gocql.TypeDate:      "TypeDate",
	gocql.TypeTime:      "TypeTime",
	gocql.TypeSmallInt:  "TypeSmallInt",
	gocql.TypeTinyInt:   "TypeTinyInt",
	gocql.TypeDuration:  "TypeDuration",
	gocql.TypeList:      "TypeList",
	gocql.TypeMap:       "TypeMap",
	gocql.TypeSet:       "TypeSet",
	gocql.TypeUDT:       "TypeUDT",
	gocql.TypeTuple:     "TypeTuple",
}

type gocqlIntType struct {
	goType       reflect.Type
	decodeHelper string
}

var gocqlIntTypes = map[gocql.Type]gocqlIntType{
	gocql.TypeTinyInt: {
		goType:       reflect.TypeOf((*int8)(nil)).Elem(),
		decodeHelper: "DecTiny",
	},
	gocql.TypeSmallInt: {
		goType:       reflect.TypeOf((*int16)(nil)).Elem(),
		decodeHelper: "DecShort",
	},
	gocql.TypeInt: {
		goType:       reflect.TypeOf((*int32)(nil)).Elem(),
		decodeHelper: "DecInt",
	},
	gocql.TypeBigInt: {
		goType:       reflect.TypeOf((*int64)(nil)).Elem(),
		decodeHelper: "DecBigInt",
	},
}

func varcharToStringDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintln(g.out, ws+out+" = "+g.getType(t)+"("+in+")")
	return nil
}

func intLikeToStringDecoder(decodeHelper string) decoderGen {
	return func(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
		ws := strings.Repeat("  ", indent)
		fmt.Fprintf(g.out, "%s%s = %s(strconv.FormatInt(int64(marshal.%s(%s)), 10))\n",
			ws, out, g.getType(t), decodeHelper, in)
		return nil
	}
}

func intLikeToIntDecoder(gocqlType gocql.Type) decoderGen {
	return func(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
		ws := strings.Repeat("  ", indent)

		gocqlTypeMeta := gocqlIntTypes[gocqlType]
		nativeVal := g.uniqueVarName()

		switch t.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Fprintf(g.out, "%s%s := marshal.%s(%s)\n", ws, nativeVal, gocqlTypeMeta.decodeHelper, in)
			if gocqlTypeMeta.goType.Bits() > t.Bits() {
				fmt.Fprintf(g.out, "%sif %s < math.MinInt%d || %s > math.MaxInt%d {\n", ws, nativeVal, t.Bits(),
					nativeVal, t.Bits())
				fmt.Fprintf(g.out, "%s  return fmt.Errorf(\"unmarshal int: value %%d out of range for %s\", %s)\n",
					ws, t.Name(), nativeVal)
				fmt.Fprintf(g.out, "%s}\n", ws)
			}
			fmt.Fprintf(g.out, "%s%s = %s(%s)\n",
				ws, out, g.getType(t), nativeVal)
		case reflect.Int:
			fmt.Fprintf(g.out, "%s%s := marshal.%s(%s)\n", ws, nativeVal, gocqlTypeMeta.decodeHelper, in)
			if gocqlTypeMeta.goType.Bits() > 32 {
				fmt.Fprintf(g.out, "%sif ^uint(0) == math.MaxUint32 && (%s < math.MinInt32 || %s > math.MaxInt32) {\n",
					ws, nativeVal, nativeVal)
				fmt.Fprintf(g.out, "%s  return fmt.Errorf(\"unmarshal int: value %%d out of range for %s\", %s)\n",
					ws, t.Name(), nativeVal)
				fmt.Fprintf(g.out, "%s}\n", ws)
			}
			fmt.Fprintf(g.out, "%s%s = %s(%s)\n",
				ws, out, g.getType(t), nativeVal)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fmt.Fprintf(g.out, "%s%s := marshal.%s(%s)\n", ws, nativeVal, gocqlTypeMeta.decodeHelper, in)
			if gocqlTypeMeta.goType.Bits() > t.Bits() {
				fmt.Fprintf(g.out, "%sif %s < 0 || %s > math.MaxUint%d {\n", ws, nativeVal, nativeVal, t.Bits())
				fmt.Fprintf(g.out, "%s  return fmt.Errorf(\"unmarshal int: value %%d out of range for %s\", %s)\n",
					ws, t.Name(), nativeVal)
				fmt.Fprintf(g.out, "%s}\n", ws)
			}
			if gocqlTypeMeta.goType.Bits() < t.Bits() {
				fmt.Fprintf(g.out, "%s%s = %s(%s) & 0x%x\n",
					ws, out, g.getType(t), nativeVal, (uint64(1)<<uint(gocqlTypeMeta.goType.Bits()))-1)
			} else {
				fmt.Fprintf(g.out, "%s%s = %s(%s)\n",
					ws, out, g.getType(t), nativeVal)
			}
		case reflect.Uint:
			fmt.Fprintf(g.out, "%s%s := marshal.%s(%s)\n", ws, nativeVal, gocqlTypeMeta.decodeHelper, in)
			if gocqlTypeMeta.goType.Bits() > 32 {
				fmt.Fprintf(g.out, "%sif ^uint(0) == math.MaxUint32 && (%s < 0 || %s > math.MaxUint32) {\n",
					ws, nativeVal, nativeVal)
				fmt.Fprintf(g.out, "%s  return fmt.Errorf(\"unmarshal int: value %%d out of range for %s\", %s)\n",
					ws, t.Name(), nativeVal)
				fmt.Fprintf(g.out, "%s}\n", ws)
			}
			if gocqlTypeMeta.goType.Bits() < 32 {
				fmt.Fprintf(g.out, "%s%s = %s(%s) & 0x%x\n",
					ws, out, g.getType(t), nativeVal, (uint64(1)<<uint(gocqlTypeMeta.goType.Bits()))-1)
			} else {
				fmt.Fprintf(g.out, "%s%s = %s(%s) & 0xffffffff\n",
					ws, out, g.getType(t), nativeVal)
			}

		default:
			return fmt.Errorf("cannot unmarshal %s into %s", gocqlType, t.Name())
		}
		return nil
	}
}

func varIntToStringDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	value := g.uniqueVarName()
	ok := g.uniqueVarName()
	fmt.Fprintf(g.out, "%s%s, %s := marshal.VarIntToInt64("+in+")\n", ws, value, ok)
	fmt.Fprintf(g.out, "%sif !%s {\n", ws, ok)
	fmt.Fprintf(g.out, "%s  return fmt.Errorf(\"unmarshal int: varint value %%v out of range for int64\", "+in+
		")\n", ws)
	fmt.Fprintf(g.out, "%s}\n", ws)
	fmt.Fprintf(g.out, "%s%s = %s(strconv.FormatInt(%s, 10))\n", ws, out, g.getType(t), value)
	return nil
}

func uuidDecoder(g *Generator, in, out string, indent int) {
	ws := strings.Repeat("  ", indent)
	err := g.uniqueVarName()
	fmt.Fprintf(g.out, "%svar %s error\n", ws, err)
	fmt.Fprintf(g.out, "%s%s, %s = gocql.UUIDFromBytes(%s)\n", ws, out, err, in)
	fmt.Fprintf(g.out, "%sif %s != nil {\n", ws, err)
	fmt.Fprintf(g.out, "%s    return fmt.Errorf(\"unable to parse UUID: %%s\", %s)\n", ws, err)
	fmt.Fprintf(g.out, "%s}\n", ws)
}

func uuidToStringDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	uuid := g.uniqueVarName()
	fmt.Fprintf(g.out, "%sif len("+in+") == 0 {\n", ws)
	fmt.Fprintf(g.out, "%s  %s = %s(\"\")\n", ws, out, g.getType(t))
	fmt.Fprintf(g.out, "%s} else {\n", ws)
	fmt.Fprintf(g.out, "%s  var %s gocql.UUID\n", ws, uuid)
	uuidDecoder(g, in, uuid, indent+1)
	fmt.Fprintf(g.out, "%s  %s = %s(%s.String())\n", ws, out, g.getType(t), uuid)
	fmt.Fprintf(g.out, "%s}\n", ws)
	return nil
}

func inetToStringDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	ip := g.uniqueVarName()
	ip4 := g.uniqueVarName()
	fmt.Fprintf(g.out, "%sif len("+in+") == 0 {\n", ws)
	fmt.Fprintf(g.out, "%s  %s = %s(\"\")\n", ws, out, g.getType(t))
	fmt.Fprintf(g.out, "%s} else {\n", ws)
	fmt.Fprintf(g.out, "%s  %s := net.IP("+in+")\n", ws, ip)
	fmt.Fprintf(g.out, "%s  if %s := %s.To4(); %s != nil {\n", ws, ip4, ip, ip4)
	fmt.Fprintf(g.out, "%s    %s = %s(%s.String())\n", ws, out, g.getType(t), ip4)
	fmt.Fprintf(g.out, "%s  } else {\n", ws)
	fmt.Fprintf(g.out, "%s    %s = %s(%s.String())\n", ws, out, g.getType(t), ip)
	fmt.Fprintf(g.out, "%s  }\n", ws)
	fmt.Fprintf(g.out, "%s}\n", ws)
	return nil
}

func varcharToBytesDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintf(g.out, "%sif "+in+" != nil {\n", ws)
	fmt.Fprintf(g.out, "%s  %s = append((%s)[:0], "+in+"...)\n", ws, out, out)
	fmt.Fprintf(g.out, "%s} else {\n", ws)
	fmt.Fprintf(g.out, "%s  %s = nil\n", ws, out)
	fmt.Fprintf(g.out, "%s}\n", ws)
	return nil
}

func uuidToBytesDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	uuid := g.uniqueVarName()
	fmt.Fprintf(g.out, "%sif len("+in+") == 0 {\n", ws)
	fmt.Fprintf(g.out, "%s  %s = %s(nil)\n", ws, out, g.getType(t))
	fmt.Fprintf(g.out, "%s} else {\n", ws)
	fmt.Fprintf(g.out, "%svar %s gocql.UUID\n", ws, uuid)

	uuidDecoder(g, in, uuid, indent+1)
	fmt.Fprintf(g.out, "%s  %s = %s[:]\n", ws, out, uuid)
	fmt.Fprintf(g.out, "%s}\n", ws)
	return nil
}

func bigIntDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintf(g.out, "%smarshal.DecBigInt2C(%s, %s)\n", ws, in, reference(out))
	return nil
}

func booleanToBoolDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintf(g.out, "%s%s = %s(marshal.DecBool(%s))\n",
		ws, out, g.getType(t), in)
	return nil
}

func decimalDecoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	scale := g.uniqueVarName()
	unscaled := g.uniqueVarName()
	fmt.Fprintf(g.out, "%sif len(%s) < 4 {\n", ws, in)
	fmt.Fprintf(g.out, "%s return fmt.Errorf(\"malformed decimal value\")\n", ws)
	fmt.Fprintf(g.out, "%s}\n", ws)
	fmt.Fprintf(g.out, "%s%s := marshal.DecInt(%s[0:4])\n", ws, scale, in)
	fmt.Fprintf(g.out, "%svar %s big.Int\n", ws, unscaled)
	fmt.Fprintf(g.out, "%smarshal.DecBigInt2C(%s[4:], %s)\n", ws, in, reference(unscaled))
	fmt.Fprintf(g.out, "%s(%s).SetUnscaledBig(%s)\n", ws, out, reference(unscaled))
	fmt.Fprintf(g.out, "%s(%s).SetScale(inf.Scale(%s))\n", ws, out, scale)
	return nil
}

func floatToFloat32Decoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintf(g.out, "%s%s = %s(math.Float32frombits(uint32(marshal.DecInt(%s))))\n",
		ws, out, g.getType(t), in)
	return nil
}

func doubleToFloat64Decoder(g *Generator, t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)
	fmt.Fprintf(g.out, "%s%s = %s(math.Float64frombits(uint64(marshal.DecBigInt(%s))))\n",
		ws, out, g.getType(t), in)
	return nil
}

// genTypeDecoder generates decoding code for the type t, but uses unmarshaler interface if implemented by t.
func (g *Generator) genTypeDecoder(t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)

	unmarshalerIface := reflect.TypeOf((*gocql.Unmarshaler)(nil)).Elem()
	if reflect.PtrTo(t).Implements(unmarshalerIface) {
		fallbackErr := g.uniqueVarName()
		fmt.Fprintln(g.out, ws+"if "+fallbackErr+" := ("+out+").UnmarshalCQL("+info+", "+in+"); "+
			fallbackErr+" != nil {")
		fmt.Fprintln(g.out, ws+"  return "+fallbackErr)
		fmt.Fprintln(g.out, ws+"}")
		return nil
	}

	err := g.genTypeDecoderNoCheck(t, info, in, out, tags, indent)
	return err
}

// sortTypes sorts types and puts preferred to the first index
func sortTypes(types []gocql.Type, preferred gocql.Type) {
	if len(types) == 0 {
		return
	}
	var sortFrom int
	for i := range types {
		if types[i] == preferred {
			types[i], types[0] = types[0], types[i]
			sortFrom = 1
		}
	}
	toSort := types[sortFrom:]
	sort.Slice(toSort, func(i, j int) bool {
		return gocqlTypes[toSort[i]] < gocqlTypes[toSort[j]]
	})
}

func decoderTypeKeys(m map[gocql.Type]decoderGen) []gocql.Type {
	keys := make([]gocql.Type, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func (g *Generator) genCQLTypeSwitch(t reflect.Type, info, in, out string, tags fieldTags, indent int, dm decoderMeta) error {
	ws := strings.Repeat("  ", indent)
	if g.conservative {
		fallbackErr := g.uniqueVarName()
		fmt.Fprintln(g.out, ws+"  if "+fallbackErr+" := gocql.Unmarshal("+info+", "+in+", "+reference(out)+
			"); "+fallbackErr+" != nil {")
		fmt.Fprintln(g.out, ws+"    return "+fallbackErr)
		fmt.Fprintln(g.out, ws+"  }")
		return nil
	}

	fmt.Fprintln(g.out, ws+"switch "+info+".Type() {")
	sortedTypes := decoderTypeKeys(dm.cqlTypes)
	preferredType := dm.preferredType
	if tags.cqlTypeSet {
		preferredType = tags.cqlType
	}
	sortTypes(sortedTypes, preferredType)
	for _, cqlType := range sortedTypes {
		gen := dm.cqlTypes[cqlType]
		fmt.Fprintln(g.out, ws+"  case gocql."+gocqlTypes[cqlType]+":")
		err := gen(g, t, info, in, out, tags, indent+1)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(g.out, ws+"  default:")
	if dm.complete {
		fmt.Fprintln(g.out, ws+"    return fmt.Errorf(\"cannot decode %s\", "+info+".Type())")
	} else {
		fallbackErr := g.uniqueVarName()
		fmt.Fprintln(g.out, ws+"    if "+fallbackErr+" := gocql.Unmarshal("+info+", "+in+", "+
			reference(out)+"); "+fallbackErr+" != nil {")
		fmt.Fprintln(g.out, ws+"      return "+fallbackErr)
		fmt.Fprintln(g.out, ws+"    }")
	}
	fmt.Fprintln(g.out, ws+"}")
	return nil
}

// genTypeDecoderNoCheck generates decoding code for the type t.
func (g *Generator) genTypeDecoderNoCheck(t reflect.Type, info, in, out string, tags fieldTags, indent int) error {
	ws := strings.Repeat("  ", indent)

	if decoderMeta, ok := decodersByType[t]; ok {
		return g.genCQLTypeSwitch(t, info, in, out, tags, indent, decoderMeta)
	}

	if decoderMeta, ok := decodersByKind[t.Kind()]; ok {
		return g.genCQLTypeSwitch(t, info, in, out, tags, indent, decoderMeta)
	}

	if t.Kind() == reflect.Ptr {
		fmt.Fprintln(g.out, ws+"if "+in+" == nil {")
		fmt.Fprintln(g.out, ws+"  "+out+" = nil")
		fmt.Fprintln(g.out, ws+"} else {")
		fmt.Fprintln(g.out, ws+"  "+out+" = new("+g.getType(t.Elem())+")")
		if err := g.genTypeDecoder(t.Elem(), info, in, "*"+out, tags, indent+1); err != nil {
			return err
		}
		fmt.Fprintln(g.out, ws+"}")
		return nil
	}

	fallbackErr := g.uniqueVarName()
	fmt.Fprintln(g.out, ws+"// fallback to gocql for "+t.String())
	fmt.Fprintln(g.out, ws+"if "+fallbackErr+" := gocql.Unmarshal("+info+", "+in+", "+reference(out)+"); "+
		fallbackErr+" != nil {")
	fmt.Fprintln(g.out, ws+"  return "+fallbackErr)
	fmt.Fprintln(g.out, ws+"}")
	return nil
}

func reference(out string) string {
	if len(out) > 0 && out[0] == '*' {
		// NOTE: In order to remove an extra reference to a pointer
		return out[1:]
	}
	return "&" + out
}

//nolint:gocritic // parameter f is huge
func (g *Generator) genStructFieldDecoder(t reflect.Type, f reflect.StructField) error {
	cqlName := g.fieldNamer.GetCQLFieldName(t, f)
	tags, err := parseFieldTags(f)
	if err != nil {
		return err
	}

	if tags.omit {
		return nil
	}

	fmt.Fprintf(g.out, "    case %q:\n", cqlName)
	if err := g.genTypeDecoder(f.Type, "udtElement.Type", "elementData", "out."+f.Name, tags, 3); err != nil {
		return err
	}

	if tags.required {
		fmt.Fprintf(g.out, "%sSet = true\n", f.Name)
	}

	return nil
}

//nolint:gocritic // parameter f is huge
func (g *Generator) genRequiredFieldSet(_ reflect.Type, f reflect.StructField) error {
	tags, err := parseFieldTags(f)
	if err != nil {
		return err
	}

	if !tags.required {
		return nil
	}

	fmt.Fprintf(g.out, "var %sSet bool\n", f.Name)
	return nil
}

//nolint:gocritic // parameter f is huge
func (g *Generator) genRequiredFieldCheck(t reflect.Type, f reflect.StructField) error {
	cqlName := g.fieldNamer.GetCQLFieldName(t, f)
	tags, err := parseFieldTags(f)
	if err != nil {
		return err
	}

	if !tags.required {
		return nil
	}

	g.imports["fmt"] = "fmt"

	fmt.Fprintf(g.out, "if !%sSet {\n", f.Name)
	fmt.Fprintf(g.out, "    return fmt.Errorf(\"key '%s' is required\")\n", cqlName)
	fmt.Fprintf(g.out, "}\n")

	return nil
}

func mergeStructFields(fields1, fields2 []reflect.StructField) (fields []reflect.StructField) {
	used := map[string]bool{}
	for _, f := range fields2 {
		used[f.Name] = true
		fields = append(fields, f)
	}

	for _, f := range fields1 {
		if !used[f.Name] {
			fields = append(fields, f)
		}
	}
	return
}

func getStructFields(t reflect.Type) ([]reflect.StructField, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("got %v; expected a struct", t)
	}

	var efields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags, err := parseFieldTags(f)
		if err != nil {
			return nil, err
		}
		if !f.Anonymous || tags.name != "" {
			continue
		}

		t1 := f.Type
		if t1.Kind() == reflect.Ptr {
			t1 = t1.Elem()
		}

		fs, err := getStructFields(t1)
		if err != nil {
			return nil, fmt.Errorf("error processing embedded field: %v", err)
		}
		efields = mergeStructFields(efields, fs)
	}

	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags, err := parseFieldTags(f)
		if err != nil {
			return nil, err
		}
		if f.Anonymous && tags.name == "" {
			continue
		}

		c := []rune(f.Name)[0]
		if unicode.IsUpper(c) {
			fields = append(fields, f)
		}
	}
	return mergeStructFields(efields, fields), nil
}

func (g *Generator) genDecoder(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return g.genSliceArrayDecoder(t)
	default:
		return g.genStructDecoder(t)
	}
}

func (g *Generator) genSliceArrayDecoder(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
	default:
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a slice/array/map type", t)
	}

	fname := g.getDecoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "func "+fname+"(info gocql.TypeInfo, data []byte, out *"+typ+") error {")
	err := g.genTypeDecoderNoCheck(t, "info", "data", "*out", fieldTags{}, 1)
	if err != nil {
		return err
	}
	fmt.Fprintln(g.out, "  return nil")
	fmt.Fprintln(g.out, "}")

	return nil
}

func (g *Generator) genStructDecoder(t reflect.Type) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct type", t)
	}

	fname := g.getDecoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "func "+fname+"(info gocql.TypeInfo, data []byte, out *"+typ+") error {")
	fmt.Fprintln(g.out, "  if data == nil {")
	fmt.Fprintln(g.out, "    return nil")
	fmt.Fprintln(g.out, "  }")
	fmt.Fprintln(g.out, "  udt, ok := info.(gocql.UDTTypeInfo)")
	fmt.Fprintln(g.out, "  if !ok {")
	fmt.Fprintf(g.out, "    return fmt.Errorf(\"cannot unmarshal non-udt type %%s to %%T\", info, out)")
	fmt.Fprintln(g.out, "  }")

	// Init embedded pointer fields.
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.Anonymous || f.Type.Kind() != reflect.Ptr {
			continue
		}
		fmt.Fprintln(g.out, "  out."+f.Name+" = new("+g.getType(f.Type.Elem())+")")
	}

	fs, err := getStructFields(t)
	if err != nil {
		return fmt.Errorf("cannot generate decoder for %v: %v", t, err)
	}

	for _, f := range fs {
		err := g.genRequiredFieldSet(t, f)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(g.out, "  for _, udtElement := range udt.Elements {")
	fmt.Fprintln(g.out, "    if len(data) == 0 {")
	fmt.Fprintln(g.out, "      return nil")
	fmt.Fprintln(g.out, "    }")
	fmt.Fprintln(g.out, "    var elementData []byte")
	fmt.Fprintln(g.out, "    elementData, data = marshal.ReadBytes(data)")

	fmt.Fprintln(g.out, "    switch udtElement.Name {")
	for _, f := range fs {
		if err := g.genStructFieldDecoder(t, f); err != nil {
			return err
		}
	}

	if g.disallowUnknownFields {
		fmt.Fprintln(g.out, "    default:")
		fmt.Fprintf(g.out, "      return fmt.Errorf(\"unknown field: %%s\", udtElement.Name)")
	}
	fmt.Fprintln(g.out, "    }")
	fmt.Fprintln(g.out, "  }")

	for _, f := range fs {
		err := g.genRequiredFieldCheck(t, f)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(g.out, "  return nil")
	fmt.Fprintln(g.out, "}")

	return nil
}

//nolint:dupl // this function is very similar to genStructMarshaler but does the opposite
func (g *Generator) genStructUnmarshaler(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
	default:
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct/slice/array/map type", t)
	}

	fname := g.getDecoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "// UnmarshalCQL implements custom unmarshaler as gocql.UnmarshalCQL")
	fmt.Fprintln(g.out, "func (v *"+typ+") UnmarshalCQL(info gocql.TypeInfo, data []byte) error {")
	fmt.Fprintln(g.out, "  return "+fname+"(info, data, v)")
	fmt.Fprintln(g.out, "}")

	return nil
}

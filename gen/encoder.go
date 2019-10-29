package gen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gocql/gocql"
)

func (g *Generator) getEncoderName(t reflect.Type) string {
	return g.functionName("encode", t)
}

// fieldTags contains parsed version of json struct field tags.
type fieldTags struct {
	name string

	omit       bool
	required   bool
	cqlTypeSet bool
	cqlType    gocql.Type
}

// parseFieldTags parses the json field tag into a structure.
func parseFieldTags(f reflect.StructField) (fieldTags, error) {
	var ret fieldTags

	for i, s := range strings.Split(f.Tag.Get("easycql"), ",") {
		switch {
		case i == 0 && s == "-":
			ret.omit = true
		case i == 0:
			ret.name = s
		case s == "required":
			ret.required = true
		case isCQLTypeName(s):
			if ret.cqlTypeSet {
				return ret, fmt.Errorf("easycql tags %s and %s conflict", ret.cqlType.String(), s)
			}
			ret.cqlType = gocqlTypeNameToID[s]
			ret.cqlTypeSet = true
		}
	}

	return ret, nil
}

var gocqlTypeNameToID = map[string]gocql.Type{
	"custom":    gocql.TypeCustom,
	"ascii":     gocql.TypeAscii,
	"bigint":    gocql.TypeBigInt,
	"blob":      gocql.TypeBlob,
	"boolean":   gocql.TypeBoolean,
	"counter":   gocql.TypeCounter,
	"decimal":   gocql.TypeDecimal,
	"double":    gocql.TypeDouble,
	"float":     gocql.TypeFloat,
	"int":       gocql.TypeInt,
	"text":      gocql.TypeText,
	"timestamp": gocql.TypeTimestamp,
	"uuid":      gocql.TypeUUID,
	"varchar":   gocql.TypeVarchar,
	"timeuuid":  gocql.TypeTimeUUID,
	"inet":      gocql.TypeInet,
	"date":      gocql.TypeDate,
	"duration":  gocql.TypeDuration,
	"time":      gocql.TypeTime,
	"smallint":  gocql.TypeSmallInt,
	"tinyint":   gocql.TypeTinyInt,
	"list":      gocql.TypeList,
	"map":       gocql.TypeMap,
	"set":       gocql.TypeSet,
	"varint":    gocql.TypeVarint,
	"tuple":     gocql.TypeTuple,
}

func isCQLTypeName(s string) bool {
	_, ok := gocqlTypeNameToID[s]
	return ok
}

// genTypeEncoder generates code that encodes in of type t into the writer, but uses marshaler interface if implemented by t.
func (g *Generator) genTypeEncoder(t reflect.Type, in string, tags fieldTags, indent int, assumeNonEmpty bool) error {
	err := g.genTypeEncoderNoCheck(t, in, tags, indent, assumeNonEmpty)
	return err
}

// genTypeEncoderNoCheck generates code that encodes in of type t into the writer.
func (g *Generator) genTypeEncoderNoCheck(t reflect.Type, in string, tags fieldTags, indent int, assumeNonEmpty bool) error {
	ws := strings.Repeat("  ", indent)

	fallbackErr := g.uniqueVarName()
	marshaledBytes := g.uniqueVarName()
	fmt.Fprintln(g.out, ws+marshaledBytes+", "+fallbackErr+" := gocql.Marshal(info, "+in+")")
	fmt.Fprintln(g.out, ws+"if "+fallbackErr+" != nil {")
	fmt.Fprintln(g.out, ws+"  return nil, "+fallbackErr)
	fmt.Fprintln(g.out, ws+"}")
	fmt.Fprintln(g.out, ws+"buf = marshal.AppendBytes(buf, "+marshaledBytes+")")

	return nil
}

func (g *Generator) notEmptyCheck(t reflect.Type, v string) string {
	switch t.Kind() {
	case reflect.Slice, reflect.Map:
		return "len(" + v + ") != 0"
	case reflect.Interface, reflect.Ptr:
		return v + " != nil"
	case reflect.Bool:
		return v
	case reflect.String:
		return v + ` != ""`
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		return v + " != 0"

	default:
		// note: Array types don't have a useful empty value
		return "true"
	}
}

func (g *Generator) genStructFieldEncoder(t reflect.Type, f reflect.StructField, first, firstCondition bool) (bool, error) {
	cqlName := g.fieldNamer.GetCQLFieldName(t, f)
	tags, err := parseFieldTags(f)
	if err != nil {
		return false, err
	}

	if tags.omit {
		return firstCondition, nil
	}

	toggleFirstCondition := firstCondition

	fmt.Fprintf(g.out, "    case %q:\n", cqlName)
	if err := g.genTypeEncoder(f.Type, "in."+f.Name, tags, 2, false); err != nil {
		return toggleFirstCondition, err
	}
	return toggleFirstCondition, nil
}

func (g *Generator) genEncoder(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return g.genSliceArrayMapEncoder(t)
	default:
		return g.genStructEncoder(t)
	}
}

func (g *Generator) genSliceArrayMapEncoder(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
	default:
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a slice/array/map type", t)
	}

	fname := g.getEncoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "func "+fname+"(info gocql.TypeInfo, in "+typ+") ([]byte, error) {")
	fmt.Fprintln(g.out, "var buf []byte")
	err := g.genTypeEncoderNoCheck(t, "in", fieldTags{}, 1, false)
	if err != nil {
		return err
	}
	fmt.Fprintln(g.out, "return buf, nil")
	fmt.Fprintln(g.out, "}")
	return nil
}

func (g *Generator) genStructEncoder(t reflect.Type) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct type", t)
	}

	fname := g.getEncoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "func "+fname+"(info gocql.TypeInfo, in "+typ+") ([]byte, error) {")
	fmt.Fprintln(g.out, "  udt, ok := info.(gocql.UDTTypeInfo)")
	fmt.Fprintln(g.out, "  if !ok {")
	fmt.Fprintf(g.out, "    return nil, fmt.Errorf(\"cannot marshal %%T to non-udt type %%s\", in, info)")
	fmt.Fprintln(g.out, "  }")
	fmt.Fprintln(g.out, "  var buf []byte")

	fs, err := getStructFields(t)
	if err != nil {
		return fmt.Errorf("cannot generate encoder for %v: %v", t, err)
	}

	fmt.Fprintln(g.out, "  for _, udtElement := range udt.Elements {")
	fmt.Fprintln(g.out, "    switch udtElement.Name {")
	firstCondition := true
	for i, f := range fs {
		firstCondition, err = g.genStructFieldEncoder(t, f, i == 0, firstCondition)

		if err != nil {
			return err
		}
	}

	if g.disallowUnknownFields {
		fmt.Fprintln(g.out, "    default:")
		fmt.Fprintf(g.out, "      return fmt.Errorf(\"unknown field: %%s\", udtElement.Name)")
	} else {
		fmt.Fprintln(g.out, "    default:")
		fmt.Fprintln(g.out, "      marshal.AppendBytes(buf, nil)")
	}
	fmt.Fprintln(g.out, "    }")
	fmt.Fprintln(g.out, "  }")
	fmt.Fprintln(g.out, "  return buf, nil")
	fmt.Fprintln(g.out, "}")

	return nil
}

func (g *Generator) genStructMarshaler(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
	default:
		return fmt.Errorf("cannot generate encoder/decoder for %v, not a struct/slice/array/map type", t)
	}

	fname := g.getEncoderName(t)
	typ := g.getType(t)

	fmt.Fprintln(g.out, "// MarshalCQL supports gocql.Marshaler interface")
	fmt.Fprintln(g.out, "func (v "+typ+") MarshalCQL(info gocql.TypeInfo) ([]byte, error) {")
	fmt.Fprintln(g.out, "  return "+fname+"(info, v)")
	fmt.Fprintln(g.out, "}")

	return nil
}

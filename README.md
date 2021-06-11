# EasyCQL

EasyCQL provides a fast way to unmarshal CQL user defined types without the use of reflection.
It generates a `UnmarshalCQL` or `MarshalCQL` methods for types.

## Usage

```
# install
go get -u github.com/kiwicom/easycql/...

# run
easycql -all <file>.go
```

This will generate marshaler and unmarshaler implementations for all structs in the given file.

If you want to generate code for some types only, you can annotate the types with `easycql:cql` on a separate
line of the comment:

```go
// MyStruct will have the code generated.
// easycql:cql
type MyStruct struct {
    Value int
}

// AnotherStruct won't have the code generated (unless you use -all)
type AnotherStruct struct {
    Value string
}
``` 

Please note that easycql requires a full Go build environment and the GOPATH environment variable
to be set. This is because easyjson code generation invokes go run on a temporary file
(an approach to code generation borrowed from easyjson, which borrowed it from ffjson).

## Modes of operation

easycql currently has two modes of operation:

* conservative mode generates code similar to what you'd get by implementing simple
  `UnmarshalUDT`/`MarshalUDT` implementations, falling back to gocql for (un)marshaling each
  field. Conservative mode can be invoked by using `-conservative` flag.
* optimized mode (the default) generates code based on the actual Go and gocql types of the fields.

## Compatibility with gocql

easycql aims for the generated code to be compatible with gocql behavior so that the generated
code can be seamlessly plugged into existing software projects.

## easycql struct tags

easycql supports struct tags on fields to guide it's behavior.

The value of `easycql` tag is a comma separated list, the first item is a cql field name:

```go
type MyStruct struct {
    MyField string `easycql:"my_field"`
}
```

When the first item is `-`, the field is ignored when marshaling/unmarshaling:

```go
type MyStruct struct {
    MyField string
    OmittedField string `easycql:"-"`
}
```

If you want to specify other items, but not the cql name, you can leave the name empty:

```go
type MyStruct struct {
    MyField string `easycql:",ascii"`
}
```

### Annotating the CQL type

While easycql knows the static type of fields in Go structs, it does not know which CQL type
will be used with that Go field most often. Knowing the cql type of a field allows easycql to
put the generated code for that type in the hot code path. 

By default easycql assumes the cql type based on Go type as in the following table:

| Go type | CQL type |
| --- | --- |
| string | varchar |
| int | int |
| int8 | tinyint |
| int16 | smallint |
| int32 | int |
| int64| bigint |
| uint | int |
| uint8 | tinyint |
| uint16 | smallint |
| uint32 | int |
| uint64 | bigint |
| bool | boolean |
| float32 | float |
| float64 | double |

If the type stored in the database does not match the one gocql assumes by default, you can add
the type in the easycql tag of the struct field. The following tags could be used:

`ascii`, `bigint`, `blob`, `boolean`, `counter`, `decimal`, `double`, `float`, `int`, `text`,
`timestamp`, `uuid`, `varchar`, `timeuuid`, `inet`, `date`, `duration`, `time`, `smallint`,
`tinyint`, `list`, `map`, `set`, `varint`, `tuple`

For example:

```go
type MyStruct struct {
    MyField string `easycql:"my_field,ascii"`
    MyInt int64 `easycql:"my_int,varint"`
}
```

## Issues, Notes, Limitations

* At the moment, optimized mode is available only for unmarshaling. Marshaling generates the same
  code in both conservative and optimized modes. Not all combinations of Go and cql types have
  generators for optimized code yet. When generator for optimized mode is not available the generated
  code will contain a fallback to `gocql.Unmarshal`.
  
* Tests for some implemented combinations of Go vs. cql types in optimized mode are missing at the moment.
  If you want to be extra safe, use conservative mode.

* easycql parser and codegen is based on reflection, so it wont work on package `main` and test files, because they
  cant be imported by the parser.

## Benchmarks

We have a benchmark for unmarshaling a large struct with different types that can be run with
`make bench`.
However this is a synthetic benchmark and your mileage may vary depending on the schema of your data.

When unmarshaling, conservative mode is roughly 4 times faster than using plain `gocql.Unmarshal` and optimized mode
is roughly 2 times faster than conservative mode according to our benchmark:

```
plain gocql vs easycql conservative:
name          old time/op    new time/op    delta
Unmarshal-12    8.78µs ±14%    1.26µs ± 0%  -85.67%  (p=0.008 n=5+5)
Marshal-12      6.07µs ± 1%    3.21µs ± 9%  -47.18%  (p=0.008 n=5+5)

name          old alloc/op   new alloc/op   delta
Unmarshal-12    3.26kB ± 0%    0.11kB ± 0%  -96.57%  (p=0.008 n=5+5)
Marshal-12      1.56kB ± 0%    1.18kB ± 0%  -24.62%  (p=0.008 n=5+5)

name          old allocs/op  new allocs/op  delta
Unmarshal-12      61.0 ± 0%       2.0 ± 0%  -96.72%  (p=0.008 n=5+5)
Marshal-12        94.0 ± 0%      55.0 ± 0%  -41.49%  (p=0.008 n=5+5)

plain gocql vs easycql optimized:
name          old time/op    new time/op    delta
Unmarshal-12    8.78µs ±14%    0.49µs ± 1%  -94.40%  (p=0.008 n=5+5)
Marshal-12      6.07µs ± 1%    3.03µs ±15%  -50.18%  (p=0.008 n=5+5)

name          old alloc/op   new alloc/op   delta
Unmarshal-12    3.26kB ± 0%    0.11kB ± 0%  -96.57%  (p=0.008 n=5+5)
Marshal-12      1.56kB ± 0%    1.18kB ± 0%  -24.62%  (p=0.008 n=5+5)

name          old allocs/op  new allocs/op  delta
Unmarshal-12      61.0 ± 0%       2.0 ± 0%  -96.72%  (p=0.008 n=5+5)
Marshal-12        94.0 ± 0%      55.0 ± 0%  -41.49%  (p=0.008 n=5+5)

easycql conservative vs easycql optimized:
name          old time/op    new time/op    delta
Unmarshal-12    1.26µs ± 0%    0.49µs ± 1%  -60.90%  (p=0.008 n=5+5)
Marshal-12      3.21µs ± 9%    3.03µs ±15%     ~     (p=0.151 n=5+5)

name          old alloc/op   new alloc/op   delta
Unmarshal-12      112B ± 0%      112B ± 0%     ~     (all equal)
Marshal-12      1.18kB ± 0%    1.18kB ± 0%     ~     (all equal)

name          old allocs/op  new allocs/op  delta
Unmarshal-12      2.00 ± 0%      2.00 ± 0%     ~     (all equal)
Marshal-12        55.0 ± 0%      55.0 ± 0%     ~     (all equal)

```

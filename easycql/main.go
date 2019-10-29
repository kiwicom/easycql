package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kiwicom/easycql/bootstrap"
	// Reference the gen package to be friendly to vendoring tools,
	// as it is an indirect dependency.
	// (The temporary bootstrapping code uses it.)
	_ "github.com/kiwicom/easycql/gen"
	"github.com/kiwicom/easycql/parser"
)

// available cli parameters
var (
	buildTags             = flag.String("build_tags", "", "build tags to add to generated file")
	snakeCase             = flag.Bool("snake_case", false, "use snake_case names instead of CamelCase by default")
	lowerCamelCase        = flag.Bool("lower_camel_case", false, "use lowerCamelCase names instead of CamelCase by default")
	allStructs            = flag.Bool("all", false, "generate marshaler/unmarshalers for all structs in a file")
	leaveTemps            = flag.Bool("leave_temps", false, "do not delete temporary files")
	stubs                 = flag.Bool("stubs", false, "only generate stubs for marshaler/unmarshaler funcs")
	noformat              = flag.Bool("noformat", false, "do not run 'gofmt -w' on output file")
	specifiedName         = flag.String("output_filename", "", "specify the filename of the output")
	processPkg            = flag.Bool("pkg", false, "process the whole package instead of just the given file")
	disallowUnknownFields = flag.Bool("disallow_unknown_fields", false, "return error if any unknown field in json appeared")
	conservative          = flag.Bool("conservative", false, "be conservative about generated code, mostly falls back to gocql")
)

func generate(fname string) (err error) {
	fInfo, err := os.Stat(fname)
	if err != nil {
		return err
	}

	p := parser.Parser{AllStructs: *allStructs}
	if err := p.Parse(fname, fInfo.IsDir()); err != nil {
		return fmt.Errorf("Error parsing %v: %v", fname, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fname, p.PkgName+"_easycql.go")
	} else {
		if s := strings.TrimSuffix(fname, ".go"); s == fname {
			return errors.New("Filename must end in '.go'")
		} else {
			outName = s + "_easycql.go"
		}
	}

	if *specifiedName != "" {
		outName = *specifiedName
	}

	var trimmedBuildTags string
	if *buildTags != "" {
		trimmedBuildTags = strings.TrimSpace(*buildTags)
	}

	g := bootstrap.Generator{
		BuildTags:             trimmedBuildTags,
		PkgPath:               p.PkgPath,
		PkgName:               p.PkgName,
		Types:                 p.StructNames,
		SnakeCase:             *snakeCase,
		LowerCamelCase:        *lowerCamelCase,
		DisallowUnknownFields: *disallowUnknownFields,
		Conservative:          *conservative,
		LeaveTemps:            *leaveTemps,
		OutName:               outName,
		StubsOnly:             *stubs,
		NoFormat:              *noformat,
	}

	if err := g.Run(); err != nil {
		return fmt.Errorf("Bootstrap failed: %v", err)
	}
	return nil
}

func main() {
	flag.Parse()

	files := flag.Args()

	gofile := os.Getenv("GOFILE")
	if *processPkg {
		gofile = filepath.Dir(gofile)
	}

	if len(files) == 0 && gofile != "" {
		files = []string{gofile}
	} else if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, fname := range files {
		if err := generate(fname); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

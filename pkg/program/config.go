package program

import (
	"flag"
	"fmt"
	"strings"

	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"gopkg.in/yaml.v2"
)

type OutputType string

const (
	OutputJson  OutputType = "json"
	OutputProto OutputType = "proto"
	OutputText  OutputType = "text"
	OutputYaml  OutputType = "yaml"
)

type Config struct {
	ProtoFiles  *protoregistry.Files
	ProtoTypes  *protoregistry.Types
	File        string
	Entrypoint  string
	Vars        starlark.StringDict
	Interactive bool
	OutputType  OutputType
	Marshaler   func(m protoreflect.ProtoMessage) ([]byte, error)
}

func Usage(errMsg string) error {
	if errMsg != "" {
		errMsg = fmt.Sprintf("\nerror: %s", errMsg)
	}
	return fmt.Errorf(`usage: grpcstar [OPTIONS...] [ARGS...]

github:
	https://github.com/stackb/grpc-starlark

options:
	-h, --help [optional, false]
		show this help screen
	-p, --protoset [required]
		filename name of proto descriptor set
	-f, --file [required]
		filename of entrypoint starlark script
		(conventionally named *.grpc.star)
	-e, --entrypoint [optional, "main"]
		name of function in global scope to invoke upon script start
	-o, --output [optional, "json", oneof "json|proto|text|yaml"]
		formatter for output protobufs returned by entrypoint function
	-i, --interactive [optional, false]
		start a REPL session (rather then exec the entrypoint)

example:
	grpcstar \
		-p routeguide.pb \
		-f routeguide.grpc.star \
		-e call_get_feature \
		longitude=35.0 latitude=109.1
%s`, errMsg)
}

func ParseConfig(args []string) (*Config, error) {
	cfg := new(Config)
	cfg.Vars = make(starlark.StringDict)

	flags := flag.NewFlagSet("grpcstar", flag.ExitOnError)

	var help bool
	var protosetFile string
	var output string

	flags.BoolVar(&help, "h", false, "show help")
	flags.BoolVar(&help, "help", false, "show help")

	flags.StringVar(&protosetFile, "p", "", "filepath to proto descriptor set")
	flags.StringVar(&protosetFile, "protoset", "", "filepath to proto descriptor set")

	flags.StringVar(&cfg.File, "f", "", "entrypoint file")
	flags.StringVar(&cfg.File, "file", "", "entrypoint file")

	flags.StringVar(&cfg.Entrypoint, "e", "main", "entrypoint function")
	flags.StringVar(&cfg.Entrypoint, "entrypoint", "main", "entrypoint function")

	flags.StringVar(&output, "o", "json", "output type (optional; one of json|proto|text|yaml)")
	flags.StringVar(&output, "output", "json", "output type (optional; one of json|proto|text|yaml)")

	flags.BoolVar(&cfg.Interactive, "i", false, "interactive mode")
	flags.BoolVar(&cfg.Interactive, "interactive", false, "interactive mode")

	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if help {
		return nil, Usage("")
	}

	if protosetFile != "" {
		files, err := protodescriptorset.LoadFiles(protosetFile)
		if err != nil {
			return nil, err
		}
		cfg.ProtoFiles = files
		cfg.ProtoTypes = protodescriptorset.FileTypes(files)
	}

	if cfg.File == "" {
		return nil, Usage("-file is mandatory")
	}

	switch OutputType(output) {
	case OutputJson:
		marshaler := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			UseProtoNames:   true,
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
		}
		cfg.OutputType = OutputJson
		cfg.Marshaler = marshaler.Marshal
	case OutputProto:
		cfg.OutputType = OutputProto
		cfg.Marshaler = proto.Marshal
	case OutputText:
		cfg.OutputType = OutputText
		cfg.Marshaler = prototext.Marshal
	case OutputYaml:
		cfg.OutputType = OutputYaml
		cfg.Marshaler = func(m protoreflect.ProtoMessage) ([]byte, error) {
			data, err := protojson.Marshal(m)
			if err != nil {
				return nil, err
			}
			var yamlMap yaml.MapSlice
			if err := yaml.Unmarshal(data, &yamlMap); err != nil {
				return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
			}
			return yaml.Marshal(yamlMap)
		}
	default:
		return nil, Usage(fmt.Sprintf("invalid flag -o: must be one of (%v|%v|%v|%v)",
			OutputJson,
			OutputProto,
			OutputText,
			OutputYaml,
		))
	}

	for _, arg := range args {
		fields := strings.SplitN(arg, "=", 2)
		if len(fields) == 2 {
			k := fields[0]
			v := starlark.String(fields[1])
			cfg.Vars[k] = v
		}
	}

	return cfg, nil
}

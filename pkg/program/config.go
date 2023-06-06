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

func ParseConfig(args []string) (*Config, error) {
	cfg := new(Config)
	cfg.Vars = make(starlark.StringDict)

	flags := flag.NewFlagSet("grpcstar", flag.ExitOnError)

	var protosetFile string
	var output string

	flags.StringVar(&protosetFile, "protoset", "", "filepath to proto descriptor set (mandatory)")
	flags.StringVar(&protosetFile, "p", "", "filepath to proto descriptor set (mandatory)")

	flags.StringVar(&cfg.File, "file", "", "entrypoint file (mandatory)")
	flags.StringVar(&cfg.File, "f", "", "entrypoint file (mandatory)")

	flags.StringVar(&cfg.Entrypoint, "entrypoint", "main", "entrypoint function (optional)")
	flags.StringVar(&cfg.Entrypoint, "e", "main", "entrypoint function (optional)")

	flags.StringVar(&output, "output", "json", "output type (optional; one of json|proto|text|yaml)")
	flags.StringVar(&output, "o", "json", "output type (optional; one of json|proto|text|yaml)")

	flags.BoolVar(&cfg.Interactive, "interactive", false, "interactive mode (REPL)")
	flags.BoolVar(&cfg.Interactive, "i", false, "interactive mode (REPL)")

	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
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
		return nil, fmt.Errorf("-file is mandatory")
	}

	switch OutputType(output) {
	case OutputJson:
		cfg.OutputType = OutputJson
		cfg.Marshaler = protojson.Marshal
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
		return nil, fmt.Errorf("invalid flag -o: must be one of (%v|%v|%v|%v)",
			OutputJson,
			OutputProto,
			OutputText,
			OutputYaml,
		)
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

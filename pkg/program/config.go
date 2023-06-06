package program

import (
	"flag"
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
)

type marshaler func(m protoreflect.ProtoMessage) ([]byte, error)

type OutputType string

const (
	OutputJson  OutputType = "json"
	OutputProto OutputType = "proto"
	OutputText  OutputType = "text"
	OutputYaml  OutputType = "yaml"
)

type Config struct {
	Protoset    string
	File        string
	Entrypoint  string
	Interactive bool
	Output      string
	Marshaler   marshaler
	Vars        starlark.StringDict
}

func ParseConfig(args []string) (*Config, error) {
	cfg := new(Config)
	cfg.Vars = make(starlark.StringDict)

	flags := flag.NewFlagSet("grpcstar", flag.ExitOnError)

	flags.StringVar(&cfg.Protoset, "protoset", "", "filepath to proto descriptor set (mandatory)")
	flags.StringVar(&cfg.Protoset, "p", "", "filepath to proto descriptor set (mandatory)")

	flags.StringVar(&cfg.File, "file", "", "entrypoint file (mandatory)")
	flags.StringVar(&cfg.File, "f", "", "entrypoint file (mandatory)")

	flags.StringVar(&cfg.Entrypoint, "entrypoint", "main", "entrypoint function (optional)")
	flags.StringVar(&cfg.Entrypoint, "e", "main", "entrypoint function (optional)")

	flags.StringVar(&cfg.Output, "output", "json", "output type (optional; one of json|proto|text|yaml)")
	flags.StringVar(&cfg.Output, "o", "json", "output type (optional; one of json|proto|text|yaml)")

	flags.BoolVar(&cfg.Interactive, "interactive", false, "interactive mode (REPL)")
	flags.BoolVar(&cfg.Interactive, "i", false, "interactive mode (REPL)")

	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if cfg.File == "" {
		return nil, fmt.Errorf("-file is mandatory")
	}

	if cfg.Protoset == "" {
		return nil, fmt.Errorf("-protoset is mandatory")
	}

	switch cfg.Output {
	case string(OutputJson):
		cfg.Marshaler = protojson.Marshal
	case string(OutputProto):
		cfg.Marshaler = proto.Marshal
	case string(OutputText):
		cfg.Marshaler = prototext.Marshal
	case string(OutputYaml):
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

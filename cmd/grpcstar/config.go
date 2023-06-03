package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	protosetFile string
	filename     string
	in           io.ReadCloser
}

func parseConfig(args []string) (*config, error) {
	cfg := new(config)

	flags := flag.NewFlagSet("grpc-starlark", flag.ExitOnError)
	flags.StringVar(&cfg.protosetFile, "protoset", "", "filepath to proto descriptor set")

	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if cfg.protosetFile == "" {
		return nil, fmt.Errorf("-protoset is mandatory")
	}

	switch len(flags.Args()) {
	case 0:
		cfg.filename = "<stdin>"
		cfg.in = os.Stdin
	case 1:
		cfg.filename = flags.Arg(0)
		f, err := os.Open(cfg.filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open starlark file %q: %w", cfg.filename, err)
		}
		cfg.in = f
	default:
		return nil, fmt.Errorf("grpc-starlark expects a single positional argument: got %v", flags.Args())
	}
	return cfg, nil
}

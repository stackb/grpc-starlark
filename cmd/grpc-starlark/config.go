package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type config struct {
	host            string
	port            string
	bindAddressFile string
	protosetFile    string
	loadFile        string
}

func parseConfig(args []string) (*config, error) {
	cfg := new(config)

	flags := flag.NewFlagSet("grpc-starlark", flag.ExitOnError)
	flags.StringVar(&cfg.host, "host", "localhost", "bind host for the gRPC server")
	flags.StringVar(&cfg.port, "port", "1234", "port for the gRPC server")
	flags.StringVar(&cfg.bindAddressFile, "bind_address_file", "", "optional filename to write server bind address to")
	flags.StringVar(&cfg.protosetFile, "protoset", "", "filepath to proto descriptor set")
	flags.StringVar(&cfg.loadFile, "load", "", "filepath to starlark handler implementations")

	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if cfg.protosetFile == "" {
		return nil, fmt.Errorf("-protoset is mandatory")
	}

	if cfg.loadFile == "" {
		return nil, fmt.Errorf("-load is mandatory")
	}

	return cfg, nil
}

func parseProtoSetFile(filename string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading protoset file: %w", err)
	}

	var dpb descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &dpb); err != nil {
		return nil, fmt.Errorf("parsing protoset file: %v", err)
	}

	return &dpb, nil
}

func makeProtoRegistryFiles(dpb *descriptorpb.FileDescriptorSet) (*protoregistry.Files, error) {
	return protodesc.NewFiles(dpb)
}

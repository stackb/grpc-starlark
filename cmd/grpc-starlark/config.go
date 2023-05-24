package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/stackb/grpc-starlark/pkg/program"
	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
)

type config struct {
	Port               string
	ProtosetFile       string
	ProtoRegistryFiles *protoregistry.Files
	HandlersFile       string
	Handlers           starlarkgrpc.HandlerMap
}

func parseConfig(args []string) (*config, error) {
	cfg := new(config)
	flags := flag.NewFlagSet("grpcmock", flag.ExitOnError)
	flags.StringVar(&cfg.Port, "port", "1234", "port for the gRPC server")
	flags.StringVar(&cfg.ProtosetFile, "protoset", "", "path to proto descriptor set file")
	flags.StringVar(&cfg.HandlersFile, "handlers", "", "path to mock handler file")
	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}
	if cfg.ProtosetFile != "" {
		dpb, err := parseProtoSetFile(cfg.ProtosetFile)
		if err != nil {
			return nil, err
		}
		files, err := makeProtoRegistryFiles(dpb)
		if err != nil {
			return nil, err
		}
		cfg.ProtoRegistryFiles = files
	}
	if cfg.HandlersFile != "" {
		reporter := func(msg string) {
			log.Println("starlark> ", msg)
		}
		errorReporter := func(err error) {
			log.Println("starlark error> ", err.Error())
		}
		handlers, err := program.LoadFile(cfg.HandlersFile, reporter, errorReporter, cfg.ProtoRegistryFiles)
		if err != nil {
			return nil, err
		}
		cfg.Handlers = handlers
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

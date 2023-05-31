package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/stackb/grpc-starlark/pkg/program"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetPrefix("grpc-starlark: ")
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return err
	}

	dpb, err := parseProtoSetFile(cfg.protosetFile)
	if err != nil {
		return err
	}

	files, err := protodesc.NewFiles(dpb)
	if err != nil {
		return err
	}
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if err := protoregistry.GlobalFiles.RegisterFile(fd); err != nil {
			log.Printf("global registerFile error: %v", err)
		}
		return true
	})

	reporter := func(msg string) {
		log.Println("starlark> ", msg)
	}

	errorReporter := func(err error) {
		log.Println("starlark error> ", err.Error())
	}

	if err := program.Load(cfg.filename, cfg.in, reporter, errorReporter, files); err != nil {
		return err
	}

	log.Printf("grpc-starlark ready (use SIGTERM to exit)")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("SIGTERM recv'd (exiting)")

	return nil
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

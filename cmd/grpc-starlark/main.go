package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/stackb/grpc-starlark/pkg/program"
	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
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

	dpb, err := protodescriptorset.ParseFile(cfg.protosetFile)
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
		log.Println("grpc-starlark> ", msg)
	}
	errorReporter := func(err error) {
		log.Println("grpc-starlark error> ", err.Error())
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

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/protobuf/reflect/protodesc"

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

	pg := program.NewProgram(files)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if _, _, err := pg.Init(cfg.filename, cfg.in); err != nil {
		return err
	}

	<-c

	return nil
}

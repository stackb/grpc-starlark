package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/stackb/grpc-starlark/pkg/program"
	"github.com/stackb/grpc-starlark/pkg/server"
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

	files, err := makeProtoRegistryFiles(dpb)
	if err != nil {
		return err
	}

	server, err := server.New(files)
	if err != nil {
		return err
	}

	reporter := func(msg string) {
		log.Println("starlark> ", msg)
	}

	errorReporter := func(err error) {
		log.Println("starlark error> ", err.Error())
	}

	if err := program.LoadFile(cfg.loadFile, reporter, errorReporter, files, server.OnHandler); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", cfg.port))
	if err != nil {
		return fmt.Errorf("failed to listen to port %s: %w", cfg.port, err)
	}

	if err := server.Start(listener); err != nil {
		return err
	}

	if cfg.bindAddressFile != "" {
		if err := os.WriteFile(cfg.bindAddressFile, []byte(listener.Addr().String()), os.ModePerm); err != nil {
			return fmt.Errorf("writing -bind_address_file %s: %v", cfg.bindAddressFile, err)
		}
		log.Printf("Server bind address written to <%s> (%s)", cfg.bindAddressFile, listener.Addr())
	}

	log.Printf("Ready at %s (use SIGTERM to exit)", listener.Addr())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("SIGTERM recv'd (exiting)")
	server.Stop()

	return nil
}

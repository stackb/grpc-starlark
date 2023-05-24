package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	server, err := server.New(cfg.ProtoRegistryFiles, cfg.Handlers)
	if err != nil {
		return err
	}
	if err := server.Start(cfg.Port); err != nil {
		return err
	}

	log.Printf("Ready on port %s (use SIGTERM to exit)", cfg.Port)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("SIGTERM recv'd (exiting)")
	server.Stop()

	return nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stackb/grpc-starlark/pkg/program"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if err := run(cwd, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-c
}

func run(wd string, args []string) error {
	cfg, err := program.ParseConfig(args)
	if err != nil {
		return err
	}

	program, err := program.NewProgram(cfg)
	if err != nil {
		return err
	}

	if err := program.Run(); err != nil {
		return err
	}

	return nil
}

func help() error {
	fmt.Fprint(os.Stderr, `usage: grpcstar <command> [args...]

grpcstar is a standalone interpreter for a dialect of starlark intended
for interaction with gRPC services.

grpcstar may be run with one of the commands below. If no command is given,
grpcstar defaults to "eval".

  run - grpcstar will load the given args as filenames.
  help - show this message.

For usage information for a specific command, run the command with the -h flag.
For example:

  grpcstar run -h

grpcstar is under active development, and its interface may change
without notice.

`)
	return flag.ErrHelp
}

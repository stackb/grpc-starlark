package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stackb/grpc-starlark/pkg/program"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := run(cwd, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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

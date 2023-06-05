package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/stackb/grpc-starlark/pkg/program"
	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
)

type command int

const (
	runCmd command = iota
	helpCmd
)

var commandFromName = map[string]command{
	"run":  runCmd,
	"help": helpCmd,
}

var nameFromCommand = []string{
	// keep in sync with definition above
	"run",
	"help",
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if err := run(cwd, os.Args[1:]); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			log.Fatal(err.Backtrace())
		}
		log.Fatal(err)
	}

	<-c
}

func run(wd string, args []string) error {
	cmd := runCmd
	if len(args) == 1 && (args[0] == "-h" || args[0] == "-help" || args[0] == "--help") {
		cmd = helpCmd
	} else if len(args) > 0 {
		c, ok := commandFromName[args[0]]
		if ok {
			cmd = c
			args = args[1:]
		}
	}

	switch cmd {
	case runCmd:
		return commandRun(wd, cmd, args)
	case helpCmd:
		return help()
	default:
		log.Panicf("unknown command: %v", cmd)
	}
	return nil
}

func help() error {
	fmt.Fprint(os.Stderr, `usage: grpcstar <command> [args...]

grpcstar is a standalone interpreter for a dialect of starlark intended
for interaction with gRPC services.

grpcstar may be run with one of the commands below. If no command is given,
grpcstar defaults to "eval".

  eval - grpcstar will load the given args as filenames.
  help - show this message.

For usage information for a specific command, run the command with the -h flag.
For example:

  grpcstar eval -h

grpcstar is under active development, and its interface may change
without notice.

`)
	return flag.ErrHelp
}

func commandRun(wd string, cmd command, args []string) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return err
	}

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	if cfg.logFile != "" {
		f, err := os.Create(cfg.logFile)
		if err != nil {
			return fmt.Errorf("creating log file: %w", err)
		}
		log.SetOutput(f)
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

	if _, _, err := pg.Init(cfg.filename, cfg.in); err != nil {
		return err
	}

	return nil
}

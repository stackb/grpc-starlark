package program

import (
	"context"
	"fmt"
	"os"

	"github.com/stripe/skycfg"
	"github.com/stripe/skycfg/go/protomodule"
	libtime "go.starlark.net/lib/time"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/stackb/grpc-starlark/pkg/starlarkcrypto"
	"github.com/stackb/grpc-starlark/pkg/starlarkgrpc"
	"github.com/stackb/grpc-starlark/pkg/starlarknet"
	"github.com/stackb/grpc-starlark/pkg/starlarkos"
	"github.com/stackb/grpc-starlark/pkg/starlarkthread"
)

type Program struct {
	cfg       *Config
	skyConfig *skycfg.Config
}

func NewProgram(cfg *Config) (*Program, error) {
	if cfg.File == "" {
		return nil, fmt.Errorf("entrypoint file is required")
	}
	skyConfig, err := skycfg.Load(context.Background(), cfg.File,
		skycfg.WithProtoRegistry(skycfg.NewUnstableProtobufRegistryV2(cfg.ProtoTypes)),
		skycfg.WithGlobals(newPredeclared(cfg.ProtoFiles, cfg.ProtoTypes)),
	)
	if err != nil {
		return nil, err
	}
	return &Program{cfg, skyConfig}, nil
}

func (p *Program) Run(options ...skycfg.ExecOption) error {
	if p.cfg.Interactive {
		p.REPL()
		return nil
	}
	msgs, err := p.Exec()
	if err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf("%s: %w", err.Backtrace(), err)
		}
		return err
	}
	if err := p.Format(msgs); err != nil {
		return err
	}
	return nil
}

func (p *Program) Exec() ([]protoreflect.ProtoMessage, error) {
	msgs, err := p.skyConfig.Main(context.Background(),
		skycfg.WithEntryPoint(p.cfg.Entrypoint),
		skycfg.WithVars(p.cfg.Vars),
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (p *Program) Format(msgs []protoreflect.ProtoMessage) error {
	var sep string
	if p.cfg.OutputType == OutputYaml {
		sep = "---\n"
	}
	for _, msg := range msgs {
		data, err := p.cfg.Marshaler(msg)
		if err != nil {
			return err
		}
		if p.cfg.OutputType == OutputProto {
			fmt.Print(data)
		} else {
			fmt.Printf("%s%s\n", sep, string(data))
		}
	}
	return nil
}

func (p *Program) REPL() {
	thread := &starlark.Thread{}
	globals := make(starlark.StringDict)
	globals["exit"] = starlark.NewBuiltin("exit", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		os.Exit(0)
		return starlark.None, nil
	})
	for key, value := range p.skyConfig.Globals() {
		globals[key] = value
	}
	for key, value := range p.skyConfig.Locals() {
		globals[key] = value
	}
	repl.REPL(thread, globals)
}

func newPredeclared(files *protoregistry.Files, types *protoregistry.Types) starlark.StringDict {
	protoModule := protomodule.NewModule(types)
	protoModule.Members["encode"] = protoEncode(types)
	protoModule.Members["decode"] = protoDecode(types)

	return starlark.StringDict{
		"os":     starlarkos.Module,
		"net":    starlarknet.Module,
		"thread": starlarkthread.Module,
		"time":   libtime.Module,
		"crypto": starlarkcrypto.Module,
		"grpc":   starlarkgrpc.NewModule(files),
		"proto":  protoModule,
		"struct": starlark.NewBuiltin("struct", starlarkstruct.Make),
		"module": starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
}

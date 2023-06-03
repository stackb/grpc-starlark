package starlarkgrpc

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoregistry"

	net "github.com/stackb/grpc-starlark/pkg/starlarknet"
	thread "github.com/stackb/grpc-starlark/pkg/starlarkthread"
)

func TestStarlarkGrpcModuleExpr(t *testing.T) {
	testCases := []struct {
		expr    string
		env     map[string]string
		wantErr string
		want    string
	}{
		// grpc.status
		{
			expr: "grpc.status",
			want: "grpc.status(ABORTED = 10, ALREADYEXISTS = 6, CANCELED = 1, DATALOSS = 15, DEADLINEEXCEEDED = 4, FAILEDPRECONDITION = 9, INTERNAL = 13, INVALIDARGUMENT = 3, NOTFOUND = 5, OK = 0, OUTOFRANGE = 11, PERMISSIONDENIED = 7, RESOURCEEXHAUSTED = 8, UNAUTHENTICATED = 16, UNAVAILABLE = 14, UNIMPLEMENTED = 12, UNKNOWN = 2)",
		},
		{
			expr: "grpc.status.OK",
			want: "0",
		},
		{
			expr:    "grpc.status.Foo",
			wantErr: "grpc.status struct has no .Foo attribute",
		},
		// grpc.Error
		{
			expr: "grpc.Error",
			want: "<built-in function grpc.Error>",
		},
		{
			expr: "grpc.Error()",
			want: `grpc.Error(code = 2, message = "")`,
		},
		{
			expr: "grpc.Error(code = grpc.status.ABORTED, message = 'user')",
			want: `grpc.Error(code = 10, message = "user")`,
		},
		// grpc.Server
		{
			expr: "grpc.Server",
			want: `<built-in function grpc.Server>`,
		},
		{
			expr: "grpc.Server()",
			want: `<grpc.Server []>`,
		},
		// grpc.Server.start
		{
			expr: "grpc.Server().start",
			want: `<built-in function grpc.Server.start>`,
		},
		{
			expr:    "grpc.Server().start()",
			wantErr: `grpc.Server.start: got 0 arguments, want 1`,
		},
		// grpc.Server.stop
		{
			expr: "grpc.Server().stop",
			want: `<built-in function grpc.Server.stop>`,
		},
		{
			expr: "grpc.Server().stop()",
			want: `None`,
		},
		{
			expr: "grpc.Server().stop(graceful = False)",
			want: `None`,
		},
		// grpc.Server.register
		{
			expr: "grpc.Server().register",
			want: `<built-in function grpc.Server.register>`,
		},
		{
			expr:    "grpc.Server().register()",
			wantErr: `grpc.Server.register: missing argument for service`,
		},
		{
			expr:    "grpc.Server().register('example.routeguide.Routeguide', {})",
			wantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
		{
			expr:    "grpc.Server().register(service = 'example.routeguide.Routeguide', handlers = {})",
			wantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
		// grpc.Channel
		{
			expr: "grpc.Channel",
			want: `<built-in function grpc.Channel>`,
		},
		{
			expr:    "grpc.Channel()",
			wantErr: `grpc.Channel: missing argument for target`,
		},
		// grpc.Client
		{
			expr: "grpc.Client",
			want: `<built-in function grpc.Client>`,
		},
		{
			expr:    "grpc.Client()",
			wantErr: `grpc.Client: missing argument for service`,
		},
		{
			expr:    "grpc.Client(service = 'example.routeguide.Routeguide')",
			wantErr: `grpc.Client: missing argument for channel`,
		},
		{
			expr:    "grpc.Client(service = 'example.routeguide.Routeguide', channel = grpc.Channel(':0'))",
			wantErr: `unknown service: example.routeguide.Routeguide (known: [])`,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			value, err := starlark.Eval(
				new(starlark.Thread),
				"<expr>",
				tc.expr,
				starlark.StringDict{
					"grpc":   NewModule(protoregistry.GlobalFiles),
					"net":    net.Module,
					"thread": thread.Module,
				},
			)
			if err != nil {
				if tc.wantErr == "" {
					t.Fatal("unexpected error: ", err)
				}
				gotErr := err.Error()
				if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
				return
			}

			got := value.String()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

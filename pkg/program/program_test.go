package program

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	libproto "go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protodesc"

	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtoDescriptor []byte

//go:embed routeguide.grpc.star
var routeguideGrpcStar string

func TestProgram(t *testing.T) {

	pds, err := protodescriptorset.Parse(routeguideProtoDescriptor)
	if err != nil {
		t.Fatal(err)
	}
	files, err := protodesc.NewFiles(pds)
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]struct {
		program     string
		wantErr     string
		wantPrinted string
	}{
		"simple case": {
			program:     "print(grpc.status.OK)",
			wantPrinted: "0",
		},
		"proto.package": {
			program: `
pb = proto.package("example.routeguide")
print(pb)
`,
			wantPrinted: `<proto.Package "example.routeguide">`,
		},
		"proto.MessageType": {
			program: `
pb = proto.package("example.routeguide")
print(pb.Point)
`,
			wantPrinted: `<proto.MessageType "example.routeguide.Point">`,
		},
		"skycfg-services-not-included": {
			program: `
pb = proto.package("example.routeguide")
print(pb.RouteGuide)
`,
			// this test demonstrates that services are not included in the proto
			// package from skycfg
			wantErr: `Protobuf type "example.routeguide.RouteGuide" not found`,
		},
		"unary method": {
			program: routeguideGrpcStar + `call_get_feature()`,
			wantPrinted: `
GetFeature: <example.routeguide.Feature name:"point (1,2)">
`,
		},
		"server streaming": {
			program: routeguideGrpcStar + `call_list_features()`,
			wantPrinted: `
ListFeatures: <example.routeguide.Feature name:"lo (1,2)">
ListFeatures: <example.routeguide.Feature name:"hi (1,4)">
`,
		},
		"client streaming": {
			program: routeguideGrpcStar + `call_record_route()`,
			wantPrinted: `
RecordRoute: <example.routeguide.RouteSummary point_count:2 distance:2 elapsed_time:10>
`,
		},
		"bidi streaming": {
			program: routeguideGrpcStar + `call_route_chat()`,
			wantPrinted: `
RouteChat: <example.routeguide.RouteNote message:"A">
RouteChat: <example.routeguide.RouteNote message:"B">
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var printed bytes.Buffer
			thread := new(starlark.Thread)
			thread.Name = "main:" + name
			thread.Print = func(thread *starlark.Thread, msg string) {
				t.Log(msg)
				printed.WriteString(msg)
				printed.WriteString("\n")
			}

			globals := newPredeclared(files)
			libproto.SetPool(thread, files)

			_, err = starlark.ExecFile(
				thread,
				"<in-memory>",
				strings.NewReader(tc.program),
				globals,
			)

			if err != nil {
				if tc.wantErr == "" {
					t.Error("unexpected error: ", err)
					return
				}
				gotErr := err.Error()
				if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
					t.Errorf("error (-want +got):\n%s", diff)
				}
				return
			}

			gotPrinted := strings.TrimSpace(printed.String())
			wantPrinted := strings.TrimSpace(tc.wantPrinted)

			if diff := cmp.Diff(wantPrinted, gotPrinted); diff != "" {
				t.Errorf("print (-want +got):\n%s", diff)
			}
		})
	}
}

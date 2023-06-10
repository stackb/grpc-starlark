package protodescriptorset

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

//go:embed routeguide_proto_descriptor.pb
var routeguideProtoDescriptor []byte

func TestMergeFilesIgnoreConflicts(t *testing.T) {
	fileDescriptorSet, err := Parse(routeguideProtoDescriptor)
	if err != nil {
		t.Fatal(err)
	}
	files, err := protodesc.NewFiles(fileDescriptorSet)
	if err != nil {
		t.Fatal(err)
	}
	var wantNames []string
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		wantNames = append(wantNames, string(fd.FullName()))
		return true
	})

	for name, tc := range map[string]struct {
		all []*protoregistry.Files
		// protoreflect.FileDescriptor is an interface, so in order to test with
		// this we'd need a custom comparer.  So, just use the names for comparison
		want []string
	}{
		"degenerate case": {
			want: nil,
		},
		"simple case": {
			all:  []*protoregistry.Files{files},
			want: wantNames,
		},
		"merge case": {
			all:  []*protoregistry.Files{files, files},
			want: wantNames,
		},
	} {
		t.Run(name, func(t *testing.T) {
			merged := MergeFilesIgnoreConflicts(tc.all...)

			var got []string
			merged.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
				got = append(got, string(fd.FullName()))
				return true
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}

}

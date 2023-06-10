package protodescriptorset

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func LoadFiles(filename string) (*protoregistry.Files, error) {
	dpb, err := LoadFileDescriptorSet(filename)
	if err != nil {
		return nil, err
	}
	files, err := protodesc.NewFiles(dpb)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func LoadFileDescriptorSet(filename string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading protoset file: %w", err)
	}
	return Unmarshal(data)
}

func Unmarshal(data []byte) (*descriptorpb.FileDescriptorSet, error) {
	var dpb descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &dpb); err != nil {
		return nil, fmt.Errorf("unmarshaling protoset file: %v", err)
	}
	return &dpb, nil
}

func Parse(data []byte) (*descriptorpb.FileDescriptorSet, error) {
	var dpb descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &dpb); err != nil {
		return nil, fmt.Errorf("parsing protoset file: %v", err)
	}

	return &dpb, nil
}

func ParseFiles(data []byte) (*protoregistry.Files, error) {
	descriptor, err := Parse(data)
	if err != nil {
		return nil, err
	}
	files, err := protodesc.NewFiles(descriptor)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func FileTypes(files *protoregistry.Files) *protoregistry.Types {
	var types protoregistry.Types
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		messages := fd.Messages()
		for i := 0; i < messages.Len(); i++ {
			md := messages.Get(i)
			msg := dynamicpb.NewMessage(md)
			msgType := msg.Type()
			types.RegisterMessage(msgType)
		}
		enums := fd.Enums()
		for i := 0; i < enums.Len(); i++ {
			ed := enums.Get(i)
			enumType := dynamicpb.NewEnumType(ed)
			types.RegisterEnum(enumType)
		}
		return true
	})
	return &types
}

func MergeFilesIgnoreConflicts(all ...*protoregistry.Files) *protoregistry.Files {
	merged := &protoregistry.Files{}
	for _, files := range all {
		files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
			// RegisterFile only return err due to file or name conflicts.  This
			// function is about ignoring conflicts, so we can ignore the error.
			merged.RegisterFile(fd)
			return true
		})
	}
	return merged
}

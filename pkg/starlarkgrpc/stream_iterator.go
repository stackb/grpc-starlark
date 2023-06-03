package starlarkgrpc

import (
	"io"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type streamIterator struct {
	recvMsg func(interface{}) error
	md      protoreflect.MessageDescriptor
}

func (it *streamIterator) Next(p *starlark.Value) bool {
	msg := dynamicpb.NewMessage(it.md)
	msg.Reset()
	if err := it.recvMsg(msg); err != nil {
		if err != io.EOF {
			log.Println("stream recv error:", err)
		}
		return false
	}
	next, err := protomodule.NewMessage(msg)
	if err != nil {
		log.Println("stream recv conversion error:", err)
		return false
	}
	*p = next
	return true
}

func (*streamIterator) Done() {}

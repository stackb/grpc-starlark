package starlarkgrpc

import (
	"fmt"
	"io"
	"log"

	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type clientStreamRecv struct {
	recvMsg func(interface{}) error
	md      protoreflect.MessageDescriptor
}

func (csr *clientStreamRecv) String() string { return fmt.Sprintf("<%s %s>", csr.Type(), csr.Name()) }

func (csr *clientStreamRecv) Type() string { return fmt.Sprintf("%T", csr) }

func (*clientStreamRecv) Truth() starlark.Bool { return starlark.True }

func (*clientStreamRecv) Freeze() {} // immutable

func (csr *clientStreamRecv) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable: %s", csr.Type())
}

func (csr *clientStreamRecv) Iterate() starlark.Iterator { return csr }

func (csr *clientStreamRecv) Next(p *starlark.Value) bool {
	msg := dynamicpb.NewMessage(csr.md)
	msg.Reset()
	if err := csr.recvMsg(msg); err != nil {
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

func (*clientStreamRecv) Done() {}

func (csr *clientStreamRecv) Name() string {
	return "ClientStreamRecv"
}

func (csr *clientStreamRecv) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	msg := dynamicpb.NewMessage(csr.md)
	msg.Reset()
	if err := csr.recvMsg(msg); err != nil {
		if err != io.EOF {
			return nil, err
		}
		return starlark.None, nil
	}
	next, err := protomodule.NewMessage(msg)
	if err != nil {
		return nil, err
	}
	return next, nil
}

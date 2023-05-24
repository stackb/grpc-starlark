package starlarkgrpc

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func getStructAttrUint64(strct *starlarkstruct.Struct, attrName string) (uint64, error) {
	value, err := strct.Attr(attrName)
	if err != nil {
		return 0, fmt.Errorf("expected attribute %q: got %v", attrName, value)
	}
	intValue, ok := value.(starlark.Int)
	if !ok {
		return 0, fmt.Errorf("expected int attribute %q: got %v (%T)", attrName, value, value)
	}
	val, ok := intValue.Uint64()
	if !ok {
		return 0, fmt.Errorf("expected uint64 attribute %q: got %v", attrName, value)
	}
	return val, nil
}

func getStructAttrString(strct *starlarkstruct.Struct, attrName string) (string, error) {
	value, err := strct.Attr(attrName)
	if err != nil {
		return "", fmt.Errorf("expected attribute %q: got %v", attrName, value)
	}
	str, ok := value.(starlark.String)
	if !ok {
		return "", fmt.Errorf("expected string attribute %q: got %v (%T)", attrName, value, value)
	}
	val := str.GoString()
	return val, nil
}

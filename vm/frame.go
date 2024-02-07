package vm

import (
	"boludolang/code"
	"boludolang/object"
)

type Frame struct {
	obj         object.Object
	ip          int
	basePointer int
}

func NewFrame(obj object.Object, basePointer int) *Frame {
	return &Frame{
		obj:         obj,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	switch obj := f.obj.(type) {
	case *object.Closure:
		return obj.Fn.Instructions
	case *object.CompiledFor:
		return obj.Instructions
	default:
		return nil
	}
}

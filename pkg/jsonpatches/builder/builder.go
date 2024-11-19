package builder

import (
	"bytes"
	"reflect"
	"strconv"
)

type GsJsonPatchBuilder struct {
	bb bytes.Buffer
}

func NewGsJsonPatchBuilder() *GsJsonPatchBuilder {
	b := &GsJsonPatchBuilder{
		bb: bytes.Buffer{},
	}

	b.bb.WriteByte('[')

	return b
}

func replace[T any](bb *bytes.Buffer, path string, value T) {
	op(bb, "replace", path, value)
}

func add[T any](bb *bytes.Buffer, path string, value T) {
	op(bb, "add", path, value)
}

func op[T any](bb *bytes.Buffer, op, path string, value T) {
	if bb.Len() != 1 {
		bb.WriteByte(',')
	}
	bb.WriteString("{\"op\":\"")
	bb.WriteString(op)
	bb.WriteString("\",\"path\":\"")
	bb.WriteString(path)

	bb.WriteString("\",\"value\":")
	if reflect.TypeOf(value).Kind() == reflect.String {
		bb.WriteByte('"')
	}

	bb.WriteString(toString(value))
	if reflect.TypeOf(value).Kind() == reflect.String {
		bb.WriteByte('"')
	}

	bb.WriteByte('}')
}

func toString[T any](value T) string {
	switch v := any(value).(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	}
	panic("unsupported toString type: " + reflect.TypeOf(value).String())
}

func (g *GsJsonPatchBuilder) ReplaceOpsState(opsState string) *GsJsonPatchBuilder {
	replace(&g.bb, "/spec/opsState", opsState)
	return g
}

func (g *GsJsonPatchBuilder) ReplaceUpdatePriority(updatePriority int) *GsJsonPatchBuilder {
	replace(&g.bb, "/spec/updatePriority", updatePriority)
	return g
}

func (g *GsJsonPatchBuilder) Build() []byte {
	g.bb.WriteByte(']')
	return g.bb.Bytes()
}

func (g *GsJsonPatchBuilder) Reset() {
	g.bb.Reset()
}

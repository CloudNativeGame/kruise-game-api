package builder

import (
	"bytes"
	"reflect"
	"strconv"
)

type GsJsonPatchBuilder struct {
	bb    bytes.Buffer
	stack []byte
}

func NewGsJsonPatchBuilder() *GsJsonPatchBuilder {
	b := &GsJsonPatchBuilder{
		bb: bytes.Buffer{},
	}

	b.bb.WriteByte('[')
	b.stack = append(b.stack, ']')

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
	for len(g.stack) != 0 {
		g.bb.WriteByte(g.popByte())
	}
	return g.bb.Bytes()
}

func (g *GsJsonPatchBuilder) popByte() byte {
	c := g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]
	return c
}

func (g *GsJsonPatchBuilder) topByte() byte {
	if len(g.stack) == 0 {
		return 0
	}
	return g.stack[len(g.stack)-1]
}

func (g *GsJsonPatchBuilder) Reset() {
	g.bb.Reset()
}

package builder

import (
	"bytes"
	"reflect"
	"strconv"
)

func replace[T any](bb *bytes.Buffer, path string, value T) {
	opWithValue(bb, "replace", path, value)
}

func add[T any](bb *bytes.Buffer, path string, value T) {
	opWithValue(bb, "add", path, value)
}

func remove(bb *bytes.Buffer, path string) {
	op(bb, "remove", path)
}

func op(bb *bytes.Buffer, op, path string) {
	if bb.Len() != 1 {
		bb.WriteByte(',')
	}

	bb.WriteString("{\"op\":\"")
	bb.WriteString(op)
	bb.WriteString("\",\"path\":\"")
	bb.WriteString(path)
	bb.WriteString("\"}")
}

func opWithValue[T any](bb *bytes.Buffer, op, path string, value T) {
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

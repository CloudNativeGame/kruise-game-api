package builder

import (
	"bytes"
)

type GsJsonPatchBuilder struct {
	bb bytes.Buffer
	bs BytesStack
}

func NewGsJsonPatchBuilder() *GsJsonPatchBuilder {
	b := &GsJsonPatchBuilder{
		bb: bytes.Buffer{},
		bs: BytesStack{},
	}

	b.bb.WriteByte('[')
	b.bs.AppendByte(']')

	return b
}

func (g *GsJsonPatchBuilder) ReplaceOpsState(opsState string) *GsJsonPatchBuilder {
	replace(&g.bb, "/spec/opsState", opsState)
	return g
}

func (g *GsJsonPatchBuilder) ReplaceUpdatePriority(updatePriority int) *GsJsonPatchBuilder {
	replace(&g.bb, "/spec/updatePriority", updatePriority)
	return g
}

func (g *GsJsonPatchBuilder) ReplaceDeletionPriority(deletionPriority int) *GsJsonPatchBuilder {
	replace(&g.bb, "/spec/deletionPriority", deletionPriority)
	return g
}

func (g *GsJsonPatchBuilder) Build() []byte {
	for !g.bs.IsEmpty() {
		g.bb.WriteByte(g.bs.PopByte())
	}
	return g.bb.Bytes()
}

func (g *GsJsonPatchBuilder) Reset() {
	g.bb.Reset()
}

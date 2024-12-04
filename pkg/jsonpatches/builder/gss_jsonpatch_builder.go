package builder

import (
	"bytes"
	"strconv"
)

type GssJsonPatchBuilder struct {
	bb bytes.Buffer
	bs BytesStack
}

func NewGssJsonPatchBuilder() *GssJsonPatchBuilder {
	b := &GssJsonPatchBuilder{
		bb: bytes.Buffer{},
		bs: BytesStack{},
	}

	b.bb.WriteByte('[')
	b.bs.AppendByte(']')

	return b
}

func (g *GssJsonPatchBuilder) ReplaceRollingUpdatePartition(partition int) *GssJsonPatchBuilder {
	replace(&g.bb, "/spec/updateStrategy/rollingUpdate/partition", partition)
	return g
}

func (g *GssJsonPatchBuilder) AppendReserveGameServerIds(ids []int) *GssJsonPatchBuilder {
	for _, id := range ids {
		add(&g.bb, "/spec/reserveGameServerIds/-", id)
	}

	return g
}

func (g *GssJsonPatchBuilder) RemoveGameServerIds(ids []int, toRemoveId int) *GssJsonPatchBuilder {
	for i, id := range ids {
		if id == toRemoveId {
			remove(&g.bb, "/spec/reserveGameServerIds/"+strconv.Itoa(i))
		}
	}
	return g
}

func (g *GssJsonPatchBuilder) ReplaceReplicas(replicas int) *GssJsonPatchBuilder {
	replace(&g.bb, "/spec/replicas", replicas)
	return g
}

func (g *GssJsonPatchBuilder) Build() []byte {
	for !g.bs.IsEmpty() {
		g.bb.WriteByte(g.bs.PopByte())
	}
	return g.bb.Bytes()
}

func (g *GssJsonPatchBuilder) Reset() {
	g.bb.Reset()
}

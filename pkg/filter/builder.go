package filter

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/builder"
)

type GsFilterBuilder struct {
	filterBuilder builder.IFilterBuilder
}

func NewGsFilterBuilder() *GsFilterBuilder {
	return &GsFilterBuilder{
		filterBuilder: builder.NewFilterBuilder(),
	}
}

func (g *GsFilterBuilder) Or() *GsFilterBuilder {
	g.filterBuilder.Or()
	return g
}

func (g *GsFilterBuilder) And() *GsFilterBuilder {
	g.filterBuilder.And()
	return g
}

func (g *GsFilterBuilder) KStringV(key string, value string) *GsFilterBuilder {
	g.filterBuilder.KStringV(key, value)
	return g
}

func (g *GsFilterBuilder) KBoolV(key string, value bool) *GsFilterBuilder {
	g.filterBuilder.KBoolV(key, value)
	return g
}

func (g *GsFilterBuilder) KNumberV(key string, value float64) *GsFilterBuilder {
	g.filterBuilder.KNumberV(key, value)
	return g
}

func (g *GsFilterBuilder) KObjectV(key string, value builder.FilterBuilderObject) *GsFilterBuilder {
	g.filterBuilder.KObjectV(key, value)
	return g
}

func (g *GsFilterBuilder) Build() string {
	return g.filterBuilder.Build()
}

func (g *GsFilterBuilder) Reset() {
	g.filterBuilder.Reset()
}

func (g *GsFilterBuilder) UpdatePriority(updatePriority int) *GsFilterBuilder {
	g.filterBuilder.KNumberV("updatePriority", float64(updatePriority))
	return g
}

func (g *GsFilterBuilder) UpdatePriorityObject(obj builder.FilterBuilderObject) *GsFilterBuilder {
	g.filterBuilder.KObjectV("updatePriority", obj)
	return g
}

func (g *GsFilterBuilder) OpsState(opsState string) *GsFilterBuilder {
	g.filterBuilder.KStringV("opsState", opsState)
	return g
}

func (g *GsFilterBuilder) OpsStateObject(obj builder.FilterBuilderObject) *GsFilterBuilder {
	g.filterBuilder.KObjectV("opsState", obj)
	return g
}

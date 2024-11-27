package builder

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/builder"
)

type GssFilterBuilder struct {
	filterBuilder builder.IFilterBuilder
}

func NewGssFilterBuilder() *GssFilterBuilder {
	return &GssFilterBuilder{
		filterBuilder: builder.NewFilterBuilder(),
	}
}

func (g *GssFilterBuilder) Or() *GssFilterBuilder {
	g.filterBuilder.Or()
	return g
}

func (g *GssFilterBuilder) And() *GssFilterBuilder {
	g.filterBuilder.And()
	return g
}

func (g *GssFilterBuilder) KStringV(key string, value string) *GssFilterBuilder {
	g.filterBuilder.KStringV(key, value)
	return g
}

func (g *GssFilterBuilder) KBoolV(key string, value bool) *GssFilterBuilder {
	g.filterBuilder.KBoolV(key, value)
	return g
}

func (g *GssFilterBuilder) KNumberV(key string, value float64) *GssFilterBuilder {
	g.filterBuilder.KNumberV(key, value)
	return g
}

func (g *GssFilterBuilder) KObjectV(key string, value builder.FilterBuilderObject) *GssFilterBuilder {
	g.filterBuilder.KObjectV(key, value)
	return g
}

func (g *GssFilterBuilder) Build() string {
	return g.filterBuilder.Build()
}

func (g *GssFilterBuilder) Reset() {
	g.filterBuilder.Reset()
}

func (g *GssFilterBuilder) Namespace(namespace string) *GssFilterBuilder {
	g.filterBuilder.KStringV("namespace", namespace)
	return g
}

func (g *GssFilterBuilder) NamespaceObject(obj builder.FilterBuilderObject) *GssFilterBuilder {
	g.filterBuilder.KObjectV("namespace", obj)
	return g
}

package gss

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scenes"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

func NewNamespaceFilter(filterFactory *factory.FilterFactory[*v1alpha1.GameServerSet]) scene_filter.ISceneFilter[*v1alpha1.GameServerSet] {
	return scenes.NewStringSceneFilter[*v1alpha1.GameServerSet]("namespace", func(gss *v1alpha1.GameServerSet) string {
		return gss.Namespace
	}, filterFactory)
}

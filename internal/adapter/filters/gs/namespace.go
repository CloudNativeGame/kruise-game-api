package gs

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scenes"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

func NewNamespaceFilter(filterFactory *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
	return scenes.NewStringSceneFilter[*v1alpha1.GameServer]("namespace", func(gs *v1alpha1.GameServer) string {
		return gs.Namespace
	}, filterFactory)
}

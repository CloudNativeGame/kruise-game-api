package filters

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scenes"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

func NewUpdatePriorityFilter(filterFactory *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
	return scenes.NewNumberSceneFilter[*v1alpha1.GameServer]("updatePriority", func(gs *v1alpha1.GameServer) float64 {
		return float64(gs.Spec.UpdatePriority.IntValue())
	}, filterFactory)
}

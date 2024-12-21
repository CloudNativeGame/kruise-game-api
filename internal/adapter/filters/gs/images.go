package gs

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scenes"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func NewImagesFilter(filterFactory *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
	return scenes.NewStringArraySceneFilter[*v1alpha1.GameServer]("images", func(gs *v1alpha1.GameServer) []string {
		return toStringArray(gs.Status.PodStatus.ContainerStatuses)
	}, filterFactory)
}

func toStringArray(containers []v1.ContainerStatus) []string {
	results := make([]string, 0, len(containers))
	for _, container := range containers {
		results = append(results, container.Name+","+container.Image)
	}

	return results
}

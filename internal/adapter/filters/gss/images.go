package gss

import (
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scenes"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func NewImagesFilter(filterFactory *factory.FilterFactory[*v1alpha1.GameServerSet]) scene_filter.ISceneFilter[*v1alpha1.GameServerSet] {
	return scenes.NewStringArraySceneFilter[*v1alpha1.GameServerSet]("images", func(gss *v1alpha1.GameServerSet) []string {
		return toStringArray(gss.Spec.GameServerTemplate.Spec.Containers)
	}, filterFactory)
}

func toStringArray(containers []v1.Container) []string {
	results := make([]string, 0, len(containers))
	for _, container := range containers {
		results = append(results, container.Name+","+container.Image)
	}

	return results
}

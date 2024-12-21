package filter

import (
	gssfilters "github.com/CloudNativeGame/kruise-game-api/internal/adapter/filters/gss"
	"github.com/CloudNativeGame/kruise-game-api/internal/queryer"
	filter "github.com/CloudNativeGame/structured-filter-go/pkg"
	filtererrors "github.com/CloudNativeGame/structured-filter-go/pkg/errors"
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

type GssFilter struct {
	filterService *filter.FilterService[*v1alpha1.GameServerSet]
	queryer       *queryer.Queryer
}

func NewGssFilter(option *FilterOption) *GssFilter {
	filterService := filter.NewFilterService[*v1alpha1.GameServerSet]().WithSceneFilters([]filter.SceneFilterCreator[*v1alpha1.GameServerSet]{
		func(f *factory.FilterFactory[*v1alpha1.GameServerSet]) scene_filter.ISceneFilter[*v1alpha1.GameServerSet] {
			return gssfilters.NewNamespaceFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServerSet]) scene_filter.ISceneFilter[*v1alpha1.GameServerSet] {
			return gssfilters.NewImagesFilter(f)
		},
	})

	return &GssFilter{
		filterService: filterService,
		queryer:       queryer.NewQueryer(&option.KubeOption),
	}
}

func (f *GssFilter) GetFilteredGameServerSets(rawFilter string) ([]*v1alpha1.GameServerSet, error) {
	gameServerSets, err := f.queryer.GetGameServerSets()
	if err != nil {
		return nil, err
	}

	filteredGameServers := make([]*v1alpha1.GameServerSet, 0)
	for _, gss := range gameServerSets {
		err := f.filterService.Match(rawFilter, gss)
		if err != nil && err.Type() == filtererrors.InvalidFilter {
			return nil, err
		}
		if err == nil {
			filteredGameServers = append(filteredGameServers, gss)
		}
	}

	return filteredGameServers, nil
}

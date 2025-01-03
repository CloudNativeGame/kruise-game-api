package filter

import (
	gsfilters "github.com/CloudNativeGame/kruise-game-api/internal/adapter/filters/gs"
	"github.com/CloudNativeGame/kruise-game-api/internal/queryer"
	filter "github.com/CloudNativeGame/structured-filter-go/pkg"
	filtererrors "github.com/CloudNativeGame/structured-filter-go/pkg/errors"
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

type GsFilter struct {
	filterService *filter.FilterService[*v1alpha1.GameServer]
	queryer       *queryer.Queryer
}

func NewGsFilter(option *FilterOption) *GsFilter {
	filterService := filter.NewFilterService[*v1alpha1.GameServer]().WithSceneFilters([]filter.SceneFilterCreator[*v1alpha1.GameServer]{
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewOpsStateFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewUpdatePriorityFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewNamespaceFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewNameFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewDeletionPriorityFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewCurrentStateFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewCurrentNetworkStateFilter(f)
		},
		func(f *factory.FilterFactory[*v1alpha1.GameServer]) scene_filter.ISceneFilter[*v1alpha1.GameServer] {
			return gsfilters.NewImagesFilter(f)
		},
	})

	return &GsFilter{
		filterService: filterService,
		queryer:       queryer.NewQueryer(&option.KubeOption),
	}
}

func (f *GsFilter) GetFilteredGameServers(rawFilter string) ([]*v1alpha1.GameServer, error) {
	gameServers, err := f.queryer.GetGameServers()
	if err != nil {
		return nil, err
	}

	filteredGameServers := make([]*v1alpha1.GameServer, 0)
	for _, gs := range gameServers {
		err := f.filterService.Match(rawFilter, gs)
		if err != nil && err.Type() == filtererrors.InvalidFilter {
			return nil, err
		}
		if err == nil {
			filteredGameServers = append(filteredGameServers, gs)
		}
	}

	return filteredGameServers, nil
}

package filter

import (
	internalfilters "github.com/CloudNativeGame/kruise-game-api/internal/adapter/filters"
	"github.com/CloudNativeGame/kruise-game-api/internal/queryer"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	filter "github.com/CloudNativeGame/structured-filter-go/pkg"
	filtererrors "github.com/CloudNativeGame/structured-filter-go/pkg/errors"
	"github.com/CloudNativeGame/structured-filter-go/pkg/factory"
	"github.com/CloudNativeGame/structured-filter-go/pkg/filters/scene_filter"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

type Filter struct {
	filterService *filter.FilterService[*v1alpha1.GameServer]
	queryer       *queryer.Queryer
}

type FilterOption struct {
	options.KubeOption
}

func NewFilter(option *FilterOption) *Filter {
	filterFactory := factory.NewFilterFactory[*v1alpha1.GameServer]()
	filterService := filter.NewFilterService(filterFactory.WithSceneFilters([]scene_filter.ISceneFilter[*v1alpha1.GameServer]{
		internalfilters.NewOpsStateFilter(filterFactory),
		internalfilters.NewUpdatePriorityFilter(filterFactory),
		internalfilters.NewNamespaceFilter(filterFactory),
	}))
	return &Filter{
		filterService: filterService,
		queryer:       queryer.NewQueryer(&option.KubeOption),
	}
}

func (f *Filter) GetFilteredGameServers(rawFilter string) ([]*v1alpha1.GameServer, error) {
	gameServers, err := f.queryer.GetGameServers()
	if err != nil {
		return nil, err
	}

	filteredGameServers := make([]*v1alpha1.GameServer, 0)
	for _, gs := range gameServers {
		err := f.filterService.MatchFilter(rawFilter, gs)
		if err != nil && err.Type() == filtererrors.InvalidFilter {
			return nil, err
		}
		if err == nil {
			filteredGameServers = append(filteredGameServers, gs)
		}
	}

	return filteredGameServers, nil
}

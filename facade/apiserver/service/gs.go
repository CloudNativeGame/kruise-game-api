package service

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/rest/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"os"
	"sync"
	"time"
)

type GsService struct {
	filter  *filter.GsFilter
	updater *updater.Updater
	deleter *deleter.Deleter
}

var gsServiceSingleton *GsService
var gsServiceOnce sync.Once

func GetGsService() *GsService {
	gsServiceOnce.Do(func() {
		gsService := newGsService()
		gsServiceSingleton = gsService
	})

	return gsServiceSingleton
}

func newGsService() *GsService {
	kubeOption := options.KubeOption{
		KubeConfigPath:          os.Getenv("KUBECONFIG"),
		InformersReSyncInterval: time.Second * 30,
	}
	return &GsService{
		filter: filter.NewGsFilter(&filter.FilterOption{
			KubeOption: kubeOption,
		}),
		updater: updater.NewUpdater(&updater.UpdaterOption{
			KubeOption: kubeOption,
		}),
		deleter: deleter.NewDeleter(&deleter.DeleterOption{
			KubeOption: kubeOption,
		}),
	}
}

func (g *GsService) GetGameServers(rawFilter string) ([]*v1alpha1.GameServer, error) {
	gameServers, err := g.filter.GetFilteredGameServers(rawFilter)
	if err != nil {
		return nil, err
	}

	return gameServers, nil
}

func (g *GsService) UpdateGameServers(request *apimodels.UpdateGameServersRequest) ([]updater.UpdateGsResult, error) {
	gameServers, err := g.filter.GetFilteredGameServers(request.Filter)
	if err != nil {
		return nil, err
	}

	results := g.updater.UpdateGameServers(gameServers, []byte(request.JsonPatch))
	return results, nil
}

func (g *GsService) DeleteGameServers(rawFilter string) ([]deleter.DeleteGsResult, error) {
	gameServers, err := g.filter.GetFilteredGameServers(rawFilter)
	if err != nil {
		return nil, err
	}

	results := g.deleter.DeleteGameServers(gameServers)
	return results, nil
}

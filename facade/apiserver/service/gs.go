package service

import (
	"os"
	"sync"
	"time"

	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/internal/queryer"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
)

type GsService struct {
	queryer *queryer.Queryer
	filter  *filter.GsFilter
	updater *updater.Updater
	deleter *deleter.Deleter
}

var (
	gsServiceSingleton *GsService
	gsServiceOnce      sync.Once
)

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
		queryer: queryer.NewQueryer(&kubeOption),
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

func (g *GsService) GetGameServer(namespace, name string) (*v1alpha1.GameServer, error) {
	return g.queryer.GetGameServer(namespace, name)
}

func (g *GsService) GetGameServers(rawFilter string) ([]*v1alpha1.GameServer, error) {
	gameServers, err := g.filter.GetFilteredGameServers(rawFilter)
	if err != nil {
		return nil, err
	}
	return gameServers, nil
}

func (g *GsService) UpdateGameServer(namespace, name string, jsonPatch []byte) (*updater.UpdateGsResult, error) {
	gameserver, err := g.queryer.GetGameServer(namespace, name)
	if err != nil {
		return nil, err
	}
	results := g.updater.UpdateGameServers([]*v1alpha1.GameServer{gameserver}, jsonPatch)
	return &results[0], nil // length is always equal 1, no need to check
}

func (g *GsService) UpdateGameServers(request *apimodels.UpdateGameServersRequest) ([]updater.UpdateGsResult, error) {
	gameServers, err := g.filter.GetFilteredGameServers(request.Filter)
	if err != nil {
		return nil, err
	}
	results := g.updater.UpdateGameServers(gameServers, []byte(request.JsonPatch))
	return results, nil
}

func (g *GsService) DeleteGameServer(namespace, name string) (*deleter.DeleteGsResult, error) {
	gameserver, err := g.queryer.GetGameServer(namespace, name)
	if err != nil {
		return nil, err
	}
	results := g.deleter.DeleteGameServers([]*v1alpha1.GameServer{gameserver})
	return &results[0], nil // length is always equal 1, no need to check
}

func (g *GsService) DeleteGameServers(rawFilter string) ([]deleter.DeleteGsResult, error) {
	gameServers, err := g.filter.GetFilteredGameServers(rawFilter)
	if err != nil {
		return nil, err
	}
	results := g.deleter.DeleteGameServers(gameServers)
	return results, nil
}

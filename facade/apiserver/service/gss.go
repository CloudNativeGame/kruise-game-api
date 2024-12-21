package service

import (
	"github.com/CloudNativeGame/kruise-game-api/facade/apiserver/apimodels"
	"github.com/CloudNativeGame/kruise-game-api/pkg/deleter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/CloudNativeGame/kruise-game-api/pkg/updater"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	"os"
	"sync"
	"time"
)

type GssService struct {
	filter  *filter.GssFilter
	updater *updater.Updater
	deleter *deleter.Deleter
}

var gssServiceSingleton *GssService
var gssServiceOnce sync.Once

func GetGssService() *GssService {
	gssServiceOnce.Do(func() {
		gssService := newGssService()
		gssServiceSingleton = gssService
	})

	return gssServiceSingleton
}

func newGssService() *GssService {
	kubeOption := options.KubeOption{
		KubeConfigPath:          os.Getenv("KUBECONFIG"),
		InformersReSyncInterval: time.Second * 30,
	}
	return &GssService{
		filter: filter.NewGssFilter(&filter.FilterOption{
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

func (g *GssService) GetGameServerSets(rawFilter string) ([]*v1alpha1.GameServerSet, error) {
	gameServerSets, err := g.filter.GetFilteredGameServerSets(rawFilter)
	if err != nil {
		return nil, err
	}

	return gameServerSets, nil
}

func (g *GssService) UpdateGameServerSets(request *apimodels.UpdateGameServerSetsRequest) ([]updater.UpdateGssResult, error) {
	gameServerSets, err := g.filter.GetFilteredGameServerSets(request.Filter)
	if err != nil {
		return nil, err
	}

	results := g.updater.UpdateGameServerSets(gameServerSets, []byte(request.JsonPatch))
	return results, nil
}

func (g *GssService) DeleteGameServerSets(rawFilter string) ([]deleter.DeleteGssResult, error) {
	gameServerSets, err := g.filter.GetFilteredGameServerSets(rawFilter)
	if err != nil {
		return nil, err
	}

	results := g.deleter.DeleteGameServerSets(gameServerSets)
	return results, nil
}

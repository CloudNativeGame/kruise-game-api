package queryer

import (
	"log/slog"

	"github.com/CloudNativeGame/kruise-game-api/internal/utils"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1alpha1client "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	"github.com/openkruise/kruise-game/pkg/client/informers/externalversions"
	v1alpha1Lister "github.com/openkruise/kruise-game/pkg/client/listers/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

type IQueryer interface {
	GetGameServers() ([]*v1alpha1.GameServer, error)
}

type Queryer struct {
	selector            labels.Selector
	gameServerLister    v1alpha1Lister.GameServerLister
	gameServerSetLister v1alpha1Lister.GameServerSetLister
	stopCh              chan struct{}
}

func NewQueryer(option *options.KubeOption) *Queryer {
	config, err := utils.BuildKubeConfig(option.KubeConfigPath)
	if err != nil {
		slog.Error("failed to build kubeConfig", "error", err)
		panic(err)
	}

	kruiseGameClient, err := v1alpha1client.NewForConfig(config)
	if err != nil {
		slog.Error("failed to create kruise game client", "error", err)
		panic(err)
	}

	gameServerInformerFactory := externalversions.NewSharedInformerFactory(kruiseGameClient, option.InformersReSyncInterval)

	queryer := &Queryer{
		selector:            labels.Everything(),
		gameServerLister:    gameServerInformerFactory.Game().V1alpha1().GameServers().Lister(),
		gameServerSetLister: gameServerInformerFactory.Game().V1alpha1().GameServerSets().Lister(),
		stopCh:              make(chan struct{}),
	}

	queryer.start(gameServerInformerFactory)

	return queryer
}

func (q *Queryer) start(gameServerInformerFactory externalversions.SharedInformerFactory) {
	go gameServerInformerFactory.Start(q.stopCh)
	informerHasSynced := gameServerInformerFactory.WaitForCacheSync(q.stopCh)
	for informer, synced := range informerHasSynced {
		if synced == false {
			slog.Error("failed to sync informer", "informer", informer)
		}
		slog.Info("informer has synced", "informer", informer)
	}
	slog.Info("all informers have synced")
}

func (q *Queryer) GetGameServer(namespace, name string) (*v1alpha1.GameServer, error) {
	return q.gameServerLister.GameServers(namespace).Get(name)
}

func (q *Queryer) GetGameServers() ([]*v1alpha1.GameServer, error) {
	gameServers, err := q.gameServerLister.List(q.selector)
	if err != nil {
		return nil, err
	}
	return gameServers, nil
}

func (q *Queryer) GetGameServerSet(namespace, name string) (*v1alpha1.GameServerSet, error) {
	return q.gameServerSetLister.GameServerSets(namespace).Get(name)
}

func (q *Queryer) GetGameServerSets() ([]*v1alpha1.GameServerSet, error) {
	gameServerSets, err := q.gameServerSetLister.List(q.selector)
	if err != nil {
		return nil, err
	}
	return gameServerSets, nil
}

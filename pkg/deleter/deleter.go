package deleter

import (
	"context"
	"github.com/CloudNativeGame/kruise-game-api/internal/utils"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1alpha1client "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
)

type Deleter struct {
	kruiseGameClient v1alpha1client.Interface
}

func NewDeleter(option *DeleterOption) *Deleter {
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

	return &Deleter{
		kruiseGameClient: kruiseGameClient,
	}
}

type DeleteGsResult struct {
	Gs  *v1alpha1.GameServer `json:"gs"`
	Err error                `json:"err"`
}

func (d *Deleter) DeleteGameServers(gameServers []*v1alpha1.GameServer) []DeleteGsResult {
	results := make([]DeleteGsResult, 0, len(gameServers))
	ctx := context.Background()
	for _, gs := range gameServers {
		err := d.deleteGameServer(ctx, gs.Name, gs.Namespace)
		results = append(results, DeleteGsResult{
			Gs:  gs,
			Err: err,
		})
	}

	return results
}

func (d *Deleter) deleteGameServer(ctx context.Context, gsName, namespace string) error {
	return d.kruiseGameClient.GameV1alpha1().GameServers(namespace).Delete(ctx, gsName, metav1.DeleteOptions{})
}

type DeleteGssResult struct {
	Gss *v1alpha1.GameServerSet `json:"gss"`
	Err error                   `json:"err"`
}

func (d *Deleter) DeleteGameServerSets(gameServerSets []*v1alpha1.GameServerSet) []DeleteGssResult {
	results := make([]DeleteGssResult, 0, len(gameServerSets))
	ctx := context.Background()
	for _, gss := range gameServerSets {
		err := d.deleteGameServerSet(ctx, gss.Name, gss.Namespace)
		results = append(results, DeleteGssResult{
			Gss: gss,
			Err: err,
		})
	}

	return results
}

func (d *Deleter) deleteGameServerSet(ctx context.Context, gssName, namespace string) error {
	return d.kruiseGameClient.GameV1alpha1().GameServerSets(namespace).Delete(ctx, gssName, metav1.DeleteOptions{})
}

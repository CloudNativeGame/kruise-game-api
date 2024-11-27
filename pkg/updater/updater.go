package updater

import (
	"context"
	"github.com/CloudNativeGame/kruise-game-api/internal/utils"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1alpha1client "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log/slog"
)

type Updater struct {
	kruiseGameClient v1alpha1client.Interface
}

func NewUpdater(option *UpdaterOption) *Updater {
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

	return &Updater{
		kruiseGameClient: kruiseGameClient,
	}
}

type UpdateGsResult struct {
	Gs        *v1alpha1.GameServer `json:"gs"`
	UpdatedGs *v1alpha1.GameServer `json:"updatedGs"`
	Err       error                `json:"err"`
}

func (u *Updater) UpdateGameServers(gameServers []*v1alpha1.GameServer, jsonPatch []byte) []UpdateGsResult {
	results := make([]UpdateGsResult, 0, len(gameServers))
	ctx := context.Background()
	for _, gs := range gameServers {
		updatedGs, err := u.updateGameServer(ctx, gs.Name, gs.Namespace, jsonPatch)
		results = append(results, UpdateGsResult{
			Gs:        gs,
			UpdatedGs: updatedGs,
			Err:       err,
		})
	}

	return results
}

func (u *Updater) updateGameServer(ctx context.Context, gsName, namespace string, jsonPatch []byte) (*v1alpha1.GameServer, error) {
	gs, err := u.kruiseGameClient.GameV1alpha1().GameServers(namespace).Patch(ctx,
		gsName, types.JSONPatchType, jsonPatch, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	return gs, nil
}

type UpdateGssResult struct {
	Gss        *v1alpha1.GameServerSet `json:"gss"`
	UpdatedGss *v1alpha1.GameServerSet `json:"updatedGss"`
	Err        error                   `json:"err"`
}

func (u *Updater) UpdateGameServerSets(gameServerSets []*v1alpha1.GameServerSet, jsonPatch []byte) []UpdateGssResult {
	results := make([]UpdateGssResult, 0, len(gameServerSets))
	ctx := context.Background()
	for _, gss := range gameServerSets {
		updatedGss, err := u.updateGameServerSet(ctx, gss.Name, gss.Namespace, jsonPatch)
		results = append(results, UpdateGssResult{
			Gss:        gss,
			UpdatedGss: updatedGss,
			Err:        err,
		})
	}

	return results
}

func (u *Updater) updateGameServerSet(ctx context.Context, gssName, namespace string, jsonPatch []byte) (*v1alpha1.GameServerSet, error) {
	gs, err := u.kruiseGameClient.GameV1alpha1().GameServerSets(namespace).Patch(ctx,
		gssName, types.JSONPatchType, jsonPatch, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	return gs, nil
}

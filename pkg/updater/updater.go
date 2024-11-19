package updater

import (
	"context"
	"github.com/CloudNativeGame/kruise-game-api/internal/utils"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"github.com/openkruise/kruise-game/apis/v1alpha1"
	v1alpha1client "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log/slog"
)

type IUpdater interface {
	Update(gameServers []*v1alpha1.GameServer, jsonPatch []byte) ([]*v1alpha1.GameServer, error)
}

type Updater struct {
	kruiseGameClient v1alpha1client.Interface
}

type UpdaterOption struct {
	options.KubeOption
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

type UpdateResult struct {
	Gs        *v1alpha1.GameServer `json:"gs"`
	UpdatedGs *v1alpha1.GameServer `json:"updatedGs"`
	Err       error                `json:"err"`
}

func (u *Updater) Update(gameServers []*v1alpha1.GameServer, jsonPatch []byte) []UpdateResult {
	results := make([]UpdateResult, 0, len(gameServers))
	ctx := context.Background()
	for _, gs := range gameServers {
		updatedGs, err := u.update(ctx, gs.Name, gs.Namespace, jsonPatch)
		results = append(results, UpdateResult{
			Gs:        gs,
			UpdatedGs: updatedGs,
			Err:       err,
		})
	}

	return results
}

func (u *Updater) update(ctx context.Context, gsName, namespace string, jsonPatch []byte) (*v1alpha1.GameServer, error) {
	gs, err := u.kruiseGameClient.GameV1alpha1().GameServers(namespace).Patch(ctx,
		gsName, types.JSONPatchType, jsonPatch, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}
	return gs, nil
}

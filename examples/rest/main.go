package main

import (
	"fmt"
	kruisegameapiclient "github.com/CloudNativeGame/kruise-game-api/facade/rest/client"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	jsonpatchbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/jsonpatches/builder"
	"log/slog"
	"os"
)

func main() {
	_ = os.Setenv("SERVER_URL", "http://192.168.2.2")

	client := kruisegameapiclient.NewKruiseGameApiHttpClient()
	f := filter.NewGsFilterBuilder().OpsState("None")
	gameServers, err := client.GetGameServers(f)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		panic(err)
	}

	slog.Info(fmt.Sprintf("%d GameServers matched filter %s", len(gameServers), f.Build()))
	for _, gs := range gameServers {
		slog.Info("filtered GS", "gs", gs)
	}

	results, err := client.UpdateGameServers(f, jsonpatchbuilder.NewGsJsonPatchBuilder().ReplaceUpdatePriority(1))
	if err != nil {
		return
	}

	for _, result := range results {
		if result.Err != nil {
			slog.Error(fmt.Sprintf("update GameServer %s/%s failed", result.Gs.Namespace, result.Gs.Name),
				"error", result.Err.Error())
		} else {
			slog.Info(fmt.Sprintf("update GameServer %s/%s success", result.Gs.Namespace, result.Gs.Name),
				"gs", result.Gs, "updatedGs", result.UpdatedGs)
		}
	}
}

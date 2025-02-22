package main

import (
	"fmt"
	kruisegameapiclient "github.com/CloudNativeGame/kruise-game-api/facade/apiserver/client"
	filterbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/filter/builder"
	jsonpatchbuilder "github.com/CloudNativeGame/kruise-game-api/pkg/jsonpatches/builder"
	"log/slog"
	"os"
)

func gsDemo() {
	client := kruisegameapiclient.NewKruiseGameApiHttpClient()
	f := filterbuilder.NewGsFilterBuilder().OpsState("None")
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

func gssDemo() {
	client := kruisegameapiclient.NewKruiseGameApiHttpClient()
	f := filterbuilder.NewGssFilterBuilder().Namespace("minecraft")
	gameServerSets, err := client.GetGameServerSets(f)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
		panic(err)
	}

	slog.Info(fmt.Sprintf("%d GameServerSets matched filter %s", len(gameServerSets), f.Build()))
	for _, gs := range gameServerSets {
		slog.Info("filtered GSS", "gss", gs)
	}
}

func main() {
	_ = os.Setenv("KRUISE_GAME_API_SERVER_URL", "http://192.168.2.2")

	gsDemo()
	gssDemo()
}

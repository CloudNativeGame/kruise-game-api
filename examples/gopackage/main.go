package main

import (
	"fmt"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	"log/slog"
	"time"
)

func main() {
	f := filter.NewFilter(&filter.FilterOption{
		KubeOption: options.KubeOption{
			KubeConfigPath:          "~/.kube/config",
			InformersReSyncInterval: time.Second * 30,
		},
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	filterBuilder := filter.NewGsFilterBuilder()
	rawFilter := filterBuilder.And().OpsState("None").UpdatePriority(0).Build()
	filterBuilder.Reset()
	//rawFilter := "{\"$and\":[{\"opsState\": \"None\"}, {\"updatePriority\": 0}]}"
	gameServers, err := f.GetFilteredGameServers(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		panic(err)
	}

	slog.Info(fmt.Sprintf("%d GameServers matched filter %s", len(gameServers), rawFilter))
	for _, gs := range gameServers {
		slog.Info("filtered GS", "gs", gs)
	}
}

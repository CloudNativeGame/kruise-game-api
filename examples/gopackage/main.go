package main

import (
	"fmt"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter"
	"github.com/CloudNativeGame/kruise-game-api/pkg/filter/builder"
	"github.com/CloudNativeGame/kruise-game-api/pkg/options"
	objectbuilder "github.com/CloudNativeGame/structured-filter-go/pkg/builder"
	"log/slog"
	"time"
)

func gsDemo() {
	gsFilter := filter.NewGsFilter(&filter.FilterOption{
		KubeOption: options.KubeOption{
			KubeConfigPath:          "~/.kube/config",
			InformersReSyncInterval: time.Second * 30,
		},
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	filterBuilder := builder.NewGsFilterBuilder()
	rawFilter := filterBuilder.And().OpsState("None").UpdatePriority(0).Build()
	filterBuilder.Reset()
	//rawFilter := "{\"$and\":[{\"opsState\": \"None\"}, {\"updatePriority\": 0}]}"
	gameServers, err := gsFilter.GetFilteredGameServers(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServers failed", "error", err)
		panic(err)
	}

	slog.Info(fmt.Sprintf("%d GameServers matched filter %s", len(gameServers), rawFilter))
	for _, gs := range gameServers {
		slog.Info("filtered GS", "gs", gs)
	}
}

func gssDemo() {
	gssFilter := filter.NewGssFilter(&filter.FilterOption{
		KubeOption: options.KubeOption{
			KubeConfigPath:          "~/.kube/config",
			InformersReSyncInterval: time.Second * 30,
		},
	})

	// TODO: fix need sleep
	time.Sleep(time.Second)

	filterBuilder := builder.NewGssFilterBuilder()
	rawFilter := filterBuilder.ImagesObject(objectbuilder.StringArrayAll(builder.ContainerImagesToStringArray([]builder.ContainerImage{
		{
			ContainerName: "busybox",
			Image:         "busybox:latest",
		},
	}))).Build()
	filterBuilder.Reset()
	//rawFilter := "{\"namespace\": \"minecraft\"}"
	gameServerSets, err := gssFilter.GetFilteredGameServerSets(rawFilter)
	if err != nil {
		slog.Error("get filtered GameServerSets failed", "error", err)
		panic(err)
	}

	slog.Info(fmt.Sprintf("%d GameServerSets matched filter %s", len(gameServerSets), rawFilter))
	for _, gs := range gameServerSets {
		slog.Info("filtered GSS", "gss", gs)
	}
}

func main() {
	gsDemo()
	gssDemo()
}
